// https://github.com/hajimehoshi/go-mp3/blob/master/example/main.go

package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	/*
		"math"
		"time"

		"github.com/go-gl/gl/v4.1-core/gl"
		"github.com/go-gl/glfw/v3.1/glfw"
		mgl "github.com/go-gl/mathgl/mgl32"

		"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

		"github.com/nicholasblaskey/animations/glfwBoilerplate"
	*/)

const pointsPerVertex = 5

func init() {
	runtime.LockOSThread()
}

/*
func makeBuffers(offset float32) (uint32, uint32, int32) {
	vertices := []float32{
		// Positions // Color coords
		-offset, -offset, 1.0, 0.0, 0.0, // Bot left
		offset, -offset, 0.0, 1.0, 0.0, // Bot right
		0.0, offset, 0.0, 0.0, 1.0, // Top
	}

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
*/

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
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		panic(err)
	}

	//read, _ :=

	// https://larsimmisch.github.io/pyalsaaudio/terminology.html
	// https://github.com/hajimehoshi/oto/blob/master/internal/mux/mux.go
	// https://www.codeproject.com/Articles/8295/MPEG-Audio-Frame-Header
	// https://github.com/hajimehoshi/go-mp3/blob/master/decode.go

	//https://github.com/hajimehoshi/oto/blob/master/player.go
	// Write writes PCM samples to the Player.
	//
	// The format is as follows:
	//   [data]      = [sample 1] [sample 2] [sample 3] ...
	//   [sample *]  = [channel 1] ...
	//   [channel *] = [byte 1] [byte 2] ...
	// Byte ordering is little endian.
	//
	// Idea is we need to take divide the sample rate into 60 buckets for fps
	// We then create frequency buckets like 20-40hz,... all the way to the max
	// Then we put a bargraph per frequency?

	for i := 0; i < 1; i++ {
		buff := make([]byte, 10000)
		n, err := d.Read(buff)

		if err != nil {
			fmt.Println(i)
			panic(err)
		}
		//fmt.Println(buff)

		sum := 0
		for j := 0; j < n; j++ {
			sum += int(buff[j])
		}
		fmt.Printf("n=%d,sum=%d\n", n, sum)

		if i == 1000-1 {
			fmt.Println(buff)
		}
	}

	fmt.Println("here?")
	go playMp3("../kubernetesMixtape.mp3")
	fmt.Println("Ending")

	fmt.Println(d.SampleRate())

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
