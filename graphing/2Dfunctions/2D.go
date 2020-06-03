// https://en.wikipedia.org/wiki/Koch_snowflake

package main

import (
	"fmt"
	//	"math"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/animations/glfwBoilerplate"
)

const pointsPerVertex = 5

type graphParams struct {
	xBoarder   float32
	yBoarder   float32
	xRange     mgl.Vec2
	yRange     mgl.Vec2
	xAxisColor mgl.Vec3
	yAxisColor mgl.Vec3
}

func init() {
	runtime.LockOSThread()
}

func makeAxisBuffs(params graphParams) (uint32, uint32, int32) {
	xCol := params.xAxisColor
	yCol := params.yAxisColor
	vertices := []float32{
		// Positions         // Color coords
		1 - params.xBoarder, 0.0, xCol[0], xCol[1], xCol[2],
		-1 + params.xBoarder, 0.0, xCol[0], xCol[1], xCol[2],
		0.0, 1 - params.yBoarder, yCol[0], yCol[1], yCol[2],
		0.0, -1 + params.yBoarder, yCol[0], yCol[1], yCol[2],
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

func main() {
	title := "2D graping"
	fmt.Println("Starting")

	params := graphParams{
		xBoarder:   0.1,
		yBoarder:   0.1,
		xRange:     mgl.Vec2{-10, 10},
		yRange:     mgl.Vec2{-10, 10},
		xAxisColor: mgl.Vec3{1.0, 1.0, 1.0},
		yAxisColor: mgl.Vec3{0.3, 0.5, 0.3},
	}

	window := glfwBoilerplate.InitGLFW(title,
		500, 500, false)
	defer glfw.Terminate()
	//gl.Enable(gl.MULTISAMPLE) // Enable anti aliasing

	ourShader := shader.MakeShaders("2D.vs", "2D.fs")
	axisVAO, axisVBO, axisVertexCount := makeAxisBuffs(params)
	//funcVAO, funcVBO := makeFunctionBuffs(params, func1, color)
	defer gl.DeleteVertexArrays(1, &axisVAO)
	defer gl.DeleteVertexArrays(1, &axisVBO)

	lastTime := 0.0
	numFrames := 0.0
	for !window.ShouldClose() {
		// Preframe
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		// Draw axis ?
		//fmt.Println(axisVertexCount)
		ourShader.Use()
		gl.BindVertexArray(axisVAO)
		gl.DrawArrays(gl.LINES, 0, axisVertexCount)
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()

		/*
			// Update triangle and VBO
			vertices = fractals.UpdateKoch(vertices)
			gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
			gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4,
				gl.Ptr(vertices), gl.STATIC_DRAW)
			gl.BindBuffer(gl.ARRAY_BUFFER, 0)
			maxIters -= 1
		*/
		time.Sleep(0 * time.Millisecond)
	}
}
