// https://en.wikipedia.org/wiki/Koch_snowflake

package main

import (
	"fmt"
	"math"
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
	dx         float32
	xAxisColor mgl.Vec3
	yAxisColor mgl.Vec3
}

type singleVarFunc func(x float32) float32

func init() {
	runtime.LockOSThread()
}

// todo change axis location depending on xRange and yRange
// FE if yRange is (0, 10) then we need not have lower quads
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

// TODO same idea change the range to fit in (I think its an issue with us starting at 0 instead of negative???)
func makeFunctionBuffs(params graphParams, fx singleVarFunc,
	col mgl.Vec3) (uint32, uint32, int32) {

	vertices := []float32{}
	for i := params.xRange[0]; i <= params.xRange[1]; i += params.dx {
		//y := fx(i)
		//if y >= params.yRange[0] && y <= params.yRange[1] {
		vertices = append(vertices,
			(2.0-params.xBoarder*2)*
				(i/(params.xRange[1]-params.xRange[0])),
			(2.0-params.yBoarder*2)*
				(fx(i)/(params.yRange[1]-params.yRange[0])),
			col[0], col[1], col[2])
		//}
	}

	fmt.Println(len(vertices) / 5)

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

func makeRibbonBuffs(params graphParams) ([]uint32, []int32) {
	funcVAOs := []uint32{}
	funcVertexCounts := []int32{}
	for i := -100; i < 100; i++ {
		funcVAO, _, funcVertexCount := makeFunctionBuffs(params,
			func(x float32) float32 {
				return float32(math.Sin(float64(x) + float64(i)*0.01))
			},
			mgl.Vec3{(float32(i))*0.006 + 0.45,
				float32(i)*0.005 + 0.35, float32(i)*0.007 + 0.55})
		funcVAOs = append(funcVAOs, funcVAO)
		funcVertexCounts = append(funcVertexCounts, funcVertexCount)
	}

	return funcVAOs, funcVertexCounts
}

func makeWaveBuffs(params graphParams) ([]uint32, []int32) {
	funcVAOs := []uint32{}
	funcVertexCounts := []int32{}
	for i := -100; i < 100; i++ {
		funcVAO, _, funcVertexCount := makeFunctionBuffs(params,
			func(x float32) float32 {
				return float32(math.Sin(float64(x))) + float32(i)*0.01
			},
			mgl.Vec3{(float32(i))*0.01 + 0.30,
				float32(i)*0.03 + 0.25, float32(i)*0.02 + 0.15})
		funcVAOs = append(funcVAOs, funcVAO)
		funcVertexCounts = append(funcVertexCounts, funcVertexCount)
	}

	return funcVAOs, funcVertexCounts
}

func makeWeirdBuffs(params graphParams) ([]uint32, []int32) {
	funcVAOs := []uint32{}
	funcVertexCounts := []int32{}
	for i := -100; i < 100; i++ {
		funcVAO, _, funcVertexCount := makeFunctionBuffs(params,
			func(x float32) float32 {
				return float32(math.Sin(float64(x*3.0))) * float32(i) * 0.03
			},
			mgl.Vec3{(float32(-i))*0.015 + 0.68,
				float32(i)*0.03 + 0.53, float32(i)*0.002 + 0.83})
		funcVAOs = append(funcVAOs, funcVAO)
		funcVertexCounts = append(funcVertexCounts, funcVertexCount)
	}

	return funcVAOs, funcVertexCounts
}

func main() {
	title := "2D graping"
	fmt.Println("Starting")

	params := graphParams{
		xBoarder:   0.1,
		yBoarder:   0.1,
		xRange:     mgl.Vec2{-10, 10},
		yRange:     mgl.Vec2{-2, 2},
		dx:         0.01,
		xAxisColor: mgl.Vec3{1.0, 1.0, 1.0},
		yAxisColor: mgl.Vec3{0.3, 0.5, 0.3},
	}

	window := glfwBoilerplate.InitGLFW(title,
		500, 500, false)
	defer glfw.Terminate()
	//gl.Enable(gl.MULTISAMPLE) // Enable anti aliasing

	ourShader := shader.MakeShaders("2D.vs", "2D.fs")
	axisVAO, axisVBO, axisVertexCount := makeAxisBuffs(params)

	//funcVAOs, funcVertexCounts := makeRibbonBuffs(params)
	//funcVAOs, funcVertexCounts := makeWaveBuffs(params)
	funcVAOs, funcVertexCounts := makeWeirdBuffs(params)

	defer gl.DeleteVertexArrays(1, &axisVAO)
	defer gl.DeleteVertexArrays(1, &axisVBO)
	//defer gl.DeleteVertexArrays(1, &funcVAO)
	//defer gl.DeleteVertexArrays(1, &funcVBO)

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
		// Draw Axis
		gl.BindVertexArray(axisVAO)
		gl.DrawArrays(gl.LINES, 0, axisVertexCount)
		gl.BindVertexArray(0)

		// Draw functions
		for i := 0; i < len(funcVAOs); i++ {
			gl.BindVertexArray(funcVAOs[i])
			//gl.DrawArray(gl.POINTS, 0, funcVertexCount)
			gl.DrawArrays(gl.LINE_STRIP, 0, funcVertexCounts[i])
			gl.BindVertexArray(0)
		}

		window.SwapBuffers()
		glfw.PollEvents()

		time.Sleep(0 * time.Millisecond)
	}
}
