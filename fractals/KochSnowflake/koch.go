package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/animations/glfwBoilerplate"
)

func init() {
	runtime.LockOSThread()
}

func makeBuffers(offset float32) (uint32, uint32) {
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
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(2*4))

	return VAO, VBO
}

func main() {
	title := "Koch snowflake"
	fmt.Println("Starting")

	window := glfwBoilerplate.InitGLFW(title,
		800, 600, false)
	defer glfw.Terminate()

	ourShader := shader.MakeShaders("koch.vs", "koch.fs")

	VAO, VBO := makeBuffers(0.75)
	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteVertexArrays(1, &VBO)

	lastTime := 0.0
	numFrames := 0.0
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	for !window.ShouldClose() {
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		// Render
		ourShader.Use()
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.LINE_LOOP, 0, 3)
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
