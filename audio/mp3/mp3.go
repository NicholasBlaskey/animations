// https://github.com/hajimehoshi/go-mp3/blob/master/example/main.go

package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

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

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	//read, _ :=

	fmt.Println("here?")
	if _, err := io.Copy(p, d); err != nil {
		panic(err)
	}
	fmt.Println("End method")
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
