// https://github.com/hajimehoshi/go-mp3/blob/master/example/main.go

package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"

	"github.com/nicholasblaskey/animations/glfwBoilerplate"
	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"
	/*
		"math"
		"time"

		mgl "github.com/go-gl/mathgl/mgl32"

	*/)

const pointsPerVertex = 5

type freqRangeInfo struct {
	maxPos float32
	minNeg float32
	rms    float32
}

func init() {
	runtime.LockOSThread()
}

func makeBuffers(vertices []float32) (uint32, uint32, int32) {
	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4,
		gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, pointsPerVertex*4,
		gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, pointsPerVertex*4,
		gl.PtrOffset(2*4))

	return VAO, VBO, int32(len(vertices) / pointsPerVertex)
}

func playMp3(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		panic(err)
	}

	// Sample rate, channelNum, bitDepthInBytes, bufferSizeInBytes?
	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("WORKING")

	f, err := os.Open("../kubernetesMixtape.mp3")
	//f, err := os.Open("C:/Users/nblas/Desktop/ltir01.mp3")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		panic(err)
	}

	// https://larsimmisch.github.io/pyalsaaudio/terminology.html
	// https://github.com/hajimehoshi/oto/blob/master/internal/mux/mux.go
	// https://www.codeproject.com/Articles/8295/MPEG-Audio-Frame-Header
	// https://github.com/hajimehoshi/go-mp3/blob/master/decode.go
	// https://github.com/hajimehoshi/oto/blob/master/player.go
	// https://wiki.multimedia.cx/index.php/PCM
	//
	// Write writes PCM samples to the Player.
	//
	// The format is as follows:
	//   [data]      = [sample 1] [sample 2] [sample 3] ...
	//   [sample *]  = [channel 1] ...
	//   [channel *] = [byte 1] [byte 2] ...
	//
	// For example:
	//   s1c1b1 s1c1b2 s1c2b1 s1c2b2 s2c1b1 s2c1b2 s2c2b1 s2c2b2
	//
	// We want to take every two bytes every 4 bytes.
	//
	// Byte ordering is little endian.
	//
	// Idea is we need to take divide the sample rate into 60 buckets for fps
	// We then create frequency buckets like 20-40hz,... all the way to the max
	// Then we put a bargraph per frequency?
	//
	// Each readCall does a frame
	// each frame lasts for
	// 26ms (26/1000 of a second). This works out to around 38fps.
	// http://www.mp3-converter.com/mp3codec/frames.htm
	//
	// https://stackoverflow.com/questions/5890499/pcm-audio-amplitude-values
	//
	// http://www.geosci.usyd.edu.au/users/jboyden/vad/
	//
	// https://stackoverflow.com/questions/26663494/algorithm-to-draw-waveform-from-audio
	samplesPerPoint := 100
	numChannels := 2
	bytesPerChannel := 2
	offset := numChannels * bytesPerChannel

	// First get all freq values

	//	maxVal := -1
	freqValues := []freqRangeInfo{}
	largestMag := float32(0)
	for true {
		buff := make([]byte, 10000)
		n, err := d.Read(buff)
		if n == 0 {
			break
		}
		if err != nil {
			panic(err)
		}

		rms := float32(0)
		maxPos := float32(0)
		minNeg := float32(0)
		for j := 0; j < n; j += offset {
			if (j/offset)%samplesPerPoint == 0 {
				rms = float32(math.Sqrt(
					float64(rms / float32(samplesPerPoint))))
				freqValues = append(freqValues,
					freqRangeInfo{maxPos, minNeg, rms})

				rms = 0.0
				maxPos = 0.0
				minNeg = 0.0
			}

			val := float32(int16(binary.LittleEndian.Uint16(buff[j:])))
			rms += val * val
			if val > maxPos {
				maxPos = val
			} else if val < minNeg {
				minNeg = val
			}

			if float32(math.Abs(float64(val))) > largestMag {
				largestMag = val
			}
		}
	}

	// Now turn frequency values into vertices
	vertices := []float32{}
	xOffset := 1.0 / float32(len(freqValues))
	j := 0

	//c1 := []float32{0.5, 0.5, 0.6}
	//c2 := []float32{0.2, 0.2, 0.4}
	c1 := []float32{0.3, 0.7, 0.3}
	c2 := []float32{0.1, 0.3, 0.3}
	for i := -len(freqValues); i < len(freqValues)-3; i += 2 {
		j += 1
		xVal := float32(i)/float32(len(freqValues)) + xOffset

		minY := float32(freqValues[j].minNeg) / (largestMag * 5)
		maxY := float32(freqValues[j].maxPos) / (largestMag * 5)
		rms := float32(freqValues[j].rms) / (largestMag * 5)
		vertices = append(vertices,
			// Draw dark line
			xVal-xOffset, maxY, c1[0], c1[1], c1[2],
			xVal+xOffset, minY, c1[0], c1[1], c1[2],
			xVal-xOffset, minY, c1[0], c1[1], c1[2],

			// Draw light line
			xVal-xOffset, maxY, c1[0], c1[1], c1[2],
			xVal+xOffset, minY, c1[0], c1[1], c1[2],
			xVal+xOffset, maxY, c1[0], c1[1], c1[2],

			// Draw upperline
			xVal-xOffset, rms, c2[0], c2[1], c2[2],
			xVal+xOffset, -rms, c2[0], c2[1], c2[2],
			xVal-xOffset, -rms, c2[0], c2[1], c2[2],

			xVal-xOffset, rms, c2[0], c2[1], c2[2],
			xVal+xOffset, -rms, c2[0], c2[1], c2[2],
			xVal+xOffset, rms, c2[0], c2[1], c2[2],
		)
	}

	//fmt.Println(vertices)
	fmt.Println("PANIC?")

	title := "Mp3"
	fmt.Println("Starting")

	window := glfwBoilerplate.InitGLFW(title,
		500, 500, false)
	defer glfw.Terminate()
	//gl.Enable(gl.MULTISAMPLE) // Enable anti aliasing

	ourShader := shader.MakeShaders("mp3.vs", "mp3.fs")

	VAO, VBO, vertexCount := makeBuffers(vertices)
	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteVertexArrays(1, &VBO)

	lastTime := 0.0
	numFrames := 0.0
	for !window.ShouldClose() {
		// Preframe
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		ourShader.Use()
		gl.BindVertexArray(VAO)
		//gl.DrawArrays(gl.POINTS, 0, vertexCount)
		gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	//fmt.Println("here?")
	//go playMp3("../kubernetesMixtape.mp3")
	//fmt.Println("Ending")

	time.Sleep(10 * time.Second)
}

/*
func main() {
	title := "Mp3"
	fmt.Println("Starting")

	window := glfwBoilerplate.InitGLFW(title,
		500, 500, false)
	defer glfw.Terminate()
	//gl.Enable(gl.MULTISAMPLE) // Enable anti aliasing

	ourShader := shader.MakeShaders("mp3.vs", "mp3.fs")

	VAO, VBO, vertexCount := makeBuffers()
	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteVertexArrays(1, &VBO)

	lastTime := 0.0
	numFrames := 0.0
	for !window.ShouldClose() {
		// Preframe
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		ourShader.Use()
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.LINES, 0, vertexCount)
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()

		time.Sleep(0 * time.Millisecond)
	}
}
*/
