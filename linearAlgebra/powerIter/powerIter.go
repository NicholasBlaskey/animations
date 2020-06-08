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
	"github.com/nicholasblaskey/animations/graphing"
)

const pointsPerVertex = 5

func init() {
	runtime.LockOSThread()
}

func makeVecBuffs(params graphing.Params2D, vec mgl.Vec2,
	col mgl.Vec3, arrowSize float32) (uint32, uint32, int32) {

	xCord := (2.0 - params.XBoarder*2) *
		(vec[0] / (params.XRange[1] - params.XRange[0]))
	yCord := (2.0 - params.YBoarder*2) *
		(vec[1] / (params.YRange[1] - params.YRange[0]))

	// https://stackoverflow.com/questions/1622762/draw-an-arrow-in-opengl-es
	v := mgl.Vec2{xCord, yCord}.Normalize()
	v1 := v.Add(mgl.Vec2{-v[1], v[0]}).Normalize().Mul(-arrowSize)
	v2 := v.Add(mgl.Vec2{v[1], -v[0]}).Normalize().Mul(-arrowSize)

	vertices := []float32{
		// Draw origin point of vec
		0, 0, col[0], col[1], col[2],
		xCord, yCord, col[0], col[1], col[2],

		xCord, yCord, col[0], col[1], col[2],
		xCord + v1[0], yCord + v1[1], col[0], col[1], col[2],

		xCord, yCord, col[0], col[1], col[2],
		xCord + v2[0], yCord + v2[1], col[0], col[1], col[2],
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
	title := "Power iteration"
	fmt.Println("Starting")

	params := graphing.Params2D{
		XBoarder:   0.1,
		YBoarder:   0.1,
		XRange:     mgl.Vec2{-10, 10},
		YRange:     mgl.Vec2{-10, 10},
		Dx:         0.01,
		XAxisColor: mgl.Vec3{1.0, 1.0, 1.0},
		YAxisColor: mgl.Vec3{0.3, 0.5, 0.3},
	}

	window := glfwBoilerplate.InitGLFW(title,
		500, 500, false)
	defer glfw.Terminate()

	ourShader := shader.MakeShaders("powerIter.vs", "powerIter.fs")

	axisVAO, axisVBO, axisVertexCount := graphing.MakeAxisBuffs(params)

	vecVAOs := []uint32{}
	vertexCounts := []int32{}

	vecVAO, _, vecVertexCount := makeVecBuffs(params,
		mgl.Vec2{5, 1}, mgl.Vec3{0.5, 0.9, 0.3}, 0.06)
	vecVAOs = append(vecVAOs, vecVAO)
	vertexCounts = append(vertexCounts, vecVertexCount)

	vecVAO, _, vecVertexCount = makeVecBuffs(params,
		mgl.Vec2{3, 3}, mgl.Vec3{0.9, 0.5, 0.3}, 0.06)
	vecVAOs = append(vecVAOs, vecVAO)
	vertexCounts = append(vertexCounts, vecVertexCount)

	vecVAO, _, vecVertexCount = makeVecBuffs(params,
		mgl.Vec2{-1, -7}, mgl.Vec3{0.5, 0.3, 0.9}, 0.06)
	vecVAOs = append(vecVAOs, vecVAO)
	vertexCounts = append(vertexCounts, vecVertexCount)

	defer gl.DeleteVertexArrays(1, &axisVAO)
	defer gl.DeleteVertexArrays(1, &axisVBO)
	//defer gl.DeleteVertexArrays(1, &vecVAO)
	//defer gl.DeleteVertexArrays(1, &vecVBO)

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

		//gl.BindVertexArray(vecVAO)
		//gl.DrawArrays(gl.LINES, 0, vecVertexCount)
		//gl.BindVertexArray(0)

		// Draw vecs
		for i := 0; i < len(vecVAOs); i++ {
			gl.BindVertexArray(vecVAOs[i])
			gl.DrawArrays(gl.LINES, 0, vertexCounts[i])
			gl.BindVertexArray(0)
		}

		window.SwapBuffers()
		glfw.PollEvents()

		time.Sleep(0 * time.Millisecond)
	}
}
