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
	"github.com/nicholasblaskey/animations/graphing"
)

func init() {
	runtime.LockOSThread()
}

func makeRibbonBuffs(params graphing.Params2D) ([]uint32, []int32) {
	funcVAOs := []uint32{}
	funcVertexCounts := []int32{}
	for i := -100; i < 100; i++ {
		funcVAO, _, funcVertexCount := graphing.MakeFunctionBuffs(params,
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

func makeWaveBuffs(params graphing.Params2D) ([]uint32, []int32) {
	funcVAOs := []uint32{}
	funcVertexCounts := []int32{}
	for i := -100; i < 100; i++ {
		funcVAO, _, funcVertexCount := graphing.MakeFunctionBuffs(params,
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

func makeWeirdBuffs(params graphing.Params2D) ([]uint32, []int32) {
	funcVAOs := []uint32{}
	funcVertexCounts := []int32{}
	for i := -100; i < 100; i++ {
		funcVAO, _, funcVertexCount := graphing.MakeFunctionBuffs(params,
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

func makeCubicBuffs(params graphing.Params2D) ([]uint32, []int32) {
	params.XRange = mgl.Vec2{-2, 2}
	params.YRange = mgl.Vec2{-3, 3}

	funcVAOs := []uint32{}
	funcVertexCounts := []int32{}
	for i := -400; i < 400; i++ {
		col := float32((i + 1000) % 20)
		if col > 10 {
			col = 1.0
		} else {
			col = 0.0
		}

		funcVAO, _, funcVertexCount := graphing.MakeFunctionBuffs(params,
			func(x float32) float32 {
				return x*x*x + float32(i)*0.005
			},
			mgl.Vec3{col, col, col})
		funcVAOs = append(funcVAOs, funcVAO)
		funcVertexCounts = append(funcVertexCounts, funcVertexCount)
	}

	return funcVAOs, funcVertexCounts
}

func makeCubic2Buffs(params graphing.Params2D) ([]uint32, []int32) {
	params.XRange = mgl.Vec2{-2, 2}
	params.YRange = mgl.Vec2{-3, 3}

	colVec1 := mgl.Vec3{0.0, 0.8, 0.8}
	colVec2 := mgl.Vec3{0.8, 0.8, 0.1}
	colVec3 := mgl.Vec3{0.8, 0.1, 0.8}

	funcVAOs := []uint32{}
	funcVertexCounts := []int32{}
	sign := float32(-1.0)
	for j := 0; j < 2; j++ {
		if j == 1 {
			sign = 1.0
		}
		for i := -400; i < 400; i++ {
			if i == 0 {
				sign = 1.0
			}

			col := float32((i + 3000) % 30)
			var colVec mgl.Vec3
			if col > 20 {
				colVec = colVec1
			} else if col > 10 {
				colVec = colVec2
			} else {
				colVec = colVec3
			}

			funcVAO, _, funcVertexCount := graphing.MakeFunctionBuffs(params,
				func(x float32) float32 {
					return sign*x*x*x + float32(i)*0.09
				},
				colVec)
			funcVAOs = append(funcVAOs, funcVAO)
			funcVertexCounts = append(funcVertexCounts, funcVertexCount)
		}
	}

	return funcVAOs, funcVertexCounts
}

func main() {
	title := "2D graping"
	fmt.Println("Starting")

	/*
		params := graphParams{
			xBoarder:   0.1,
			yBoarder:   0.1,
			xRange:     mgl.Vec2{-10, 10},
			yRange:     mgl.Vec2{-2, 2},
			dx:         0.01,
			xAxisColor: mgl.Vec3{1.0, 1.0, 1.0},
			yAxisColor: mgl.Vec3{0.3, 0.5, 0.3},
		}
	*/

	params := graphing.Params2D{
		XBoarder:   0.1,
		YBoarder:   0.1,
		XRange:     mgl.Vec2{-2, 2},
		YRange:     mgl.Vec2{0, 2},
		Dx:         0.01,
		XAxisColor: mgl.Vec3{1.0, 1.0, 1.0},
		YAxisColor: mgl.Vec3{0.3, 0.5, 0.3},
	}

	window := glfwBoilerplate.InitGLFW(title,
		500, 500, false)
	defer glfw.Terminate()
	//gl.Enable(gl.MULTISAMPLE) // Enable anti aliasing

	ourShader := shader.MakeShaders("2D.vs", "2D.fs")
	axisVAO, axisVBO, axisVertexCount := graphing.MakeAxisBuffs(params)

	// TODO non centered axis
	//funcVAOs, funcVertexCounts := makeNonCenteredAxisBuffs(params)
	//funcVAOs, funcVertexCounts := makeRibbonBuffs(params)
	//funcVAOs, funcVertexCounts := makeWaveBuffs(params)
	//funcVAOs, funcVertexCounts := makeCubicBuffs(params)
	funcVAOs, funcVertexCounts := makeCubic2Buffs(params)
	//funcVAOs, funcVertexCounts := makeWeirdBuffs(params)

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
