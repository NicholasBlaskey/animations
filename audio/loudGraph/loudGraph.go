// https://github.com/hajimehoshi/go-mp3/blob/master/example/main.go

package main

import (
	"encoding/binary"
	"fmt"

	"math"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/hajimehoshi/go-mp3"
	//	"github.com/hajimehoshi/oto"

	"github.com/nicholasblaskey/animations/glfwBoilerplate"
	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"
)

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

func getFreqValues(d *mp3.Decoder, samplesPerPoint,
	numChannels, bytesPerChannel int) ([]freqRangeInfo, float32) {

	offset := numChannels * bytesPerChannel
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

	return freqValues, largestMag
}

func freqIntoVertices(freqValues []freqRangeInfo, xShift, largestMag,
	scaleFactor float32, c1 []float32, c2 []float32) []float32 {

	vertices := []float32{}
	xOffset := 1.0 / float32(len(freqValues))
	j := 0
	for i := -len(freqValues); i < len(freqValues)-3; i += 2 {
		j += 1
		xVal := float32(i)/float32(len(freqValues)) + xOffset + xShift

		minY := float32(freqValues[j].minNeg) / (largestMag * scaleFactor)
		maxY := float32(freqValues[j].maxPos) / (largestMag * scaleFactor)
		rms := float32(freqValues[j].rms) / (largestMag * scaleFactor)
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
	return vertices
}

func singleGraph(fileName string) ([]uint32, []uint32, []int32) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		panic(err)
	}

	freqValues, largestMag := getFreqValues(d, 100, 2, 2)
	xShift := float32(0.0)
	scaleFactor := float32(5.0)
	vertices := freqIntoVertices(freqValues, xShift, largestMag,
		scaleFactor, []float32{0.3, 0.7, 0.3}, []float32{0.1, 0.3, 0.3})

	VAO, VBO, vertexCount := makeBuffers(vertices)

	return []uint32{VAO}, []uint32{VBO}, []int32{vertexCount}
}

func stackedGraphs(fileName string) ([]uint32, []uint32, []int32) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		panic(err)
	}

	freqValues, largestMag := getFreqValues(d, 100, 2, 2)
	VAOs := []uint32{}
	VBOs := []uint32{}
	vertexCounts := []int32{}
	for i := float32(30.0); i > 1.0; i -= float32(0.2) {
		c1 := []float32{-i/30.0 + 0.5, -i/30.0 + 0.3, -i/5.0 + 0.2}
		vertices := freqIntoVertices(freqValues, i/20.0-1.0, largestMag,
			i/1.5, c1, c1)
		VAO, VBO, vertexCount := makeBuffers(vertices)

		VAOs = append(VAOs, VAO)
		VBOs = append(VBOs, VBO)
		vertexCounts = append(vertexCounts, vertexCount)
	}

	return VAOs, VBOs, vertexCounts
}

func readFramesGraph(fileName string) ([]uint32, []uint32, []int32) {

	buff := make([]byte, 10000)

	VAOs := []uint32{}
	VBOs := []uint32{}
	vertexCounts := []int32{}
	for i := 0; i < 25; i++ {
		f, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		d, err := mp3.NewDecoder(f)

		if err != nil {
			panic(err)
		}
		// Read i amount of frames to advance the buffer
		for j := 0; j < i*50; j++ {
			d.Read(buff)
		}
		freqValues, largestMag := getFreqValues(d, 100, 2, 2)
		vertices := freqIntoVertices(freqValues, 0.0, largestMag,
			float32(i), []float32{0.3, 0.7, 0.3}, []float32{0.1, 0.3, 0.3})

		VAO, VBO, vertexCount := makeBuffers(vertices)
		VAOs = append(VAOs, VAO)
		VBOs = append(VBOs, VBO)
		vertexCounts = append(vertexCounts, vertexCount)
	}

	return VAOs, VBOs, vertexCounts
}

func main() {
	fmt.Println("Starting")

	title := "Mp3"
	window := glfwBoilerplate.InitGLFW(title,
		500, 500, false)
	defer glfw.Terminate()

	ourShader := shader.MakeShaders("loudGraph.vs", "loudGraph.fs")

	fileName := "../kubernetesMixtape.mp3"

	//VAOs, _, vertexCounts := singleGraph(fileName)
	//VAOs, _, vertexCounts := stackedGraphs(fileName)
	VAOs, _, vertexCounts := readFramesGraph(fileName)

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
		for i, VAO := range VAOs {
			gl.BindVertexArray(VAO)
			gl.DrawArrays(gl.TRIANGLES, 0, vertexCounts[i])
			gl.BindVertexArray(0)
		}
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
