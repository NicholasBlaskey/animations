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

func init() {
	runtime.LockOSThread()
}

func updateVertices(vertices []float32) []float32 {
	updatedVertices := []float32{}

	// https://goinbigdata.com/golang-wait-for-all-goroutines-to-finish/
	// TODO rewrite with concurency in mind
	// allocate
	//fmt.Println("STARTING")

	// For each line segment
	n := len(vertices)
	for i := 0; i < n; i += 5 {
		from := vertices[i : i+pointsPerVertex]
		startingIndex := (i + pointsPerVertex) % n
		to := vertices[startingIndex : startingIndex+pointsPerVertex]

		// Divide each segment into 3 equal parts keeping track of the froms
		// of each segment (tos is just froms[i + 1] and the to value
		ratio := float32(1.0 / 3.0)
		froms := make([][]float32, 3)
		for i := range froms {
			froms[i] = make([]float32, 5)
		}
		froms[0] = from
		for j := range from {
			froms[1][j] = from[j]*ratio*2 + to[j]*ratio*1
			froms[2][j] = from[j]*ratio*1 + to[j]*ratio*2
		}

		// Add in first segment
		updatedVertices = append(updatedVertices, froms[0]...)
		updatedVertices = append(updatedVertices, froms[1]...)

		// Get third triangle point using this
		// This method
		//https://stackoverflow.com/questions/50547068/creating-an-equilateral-triangle-for-given-two-points-in-the-plane-python
		mid := mgl.Vec2{(froms[1][0] + froms[2][0]) / 2.0,
			(froms[1][1] + froms[2][1]) / 2.0}
		orig := mgl.Vec2{(froms[1][0] - mid[0]), (froms[1][1] - mid[1])}
		orig.Mul(3 * float32(math.Sqrt(3)))
		transform := mgl.Rotate2D(mgl.DegToRad(90))
		point := mid.Add(transform.Mul2x1(orig))
		fullPoint := []float32{point[0], point[1],
			froms[1][2] + froms[2][2], froms[1][3] + froms[2][3],
			froms[1][4] + froms[2][4]}

		// Add in the triangle segments
		updatedVertices = append(updatedVertices, froms[1]...)
		updatedVertices = append(updatedVertices, fullPoint...)
		updatedVertices = append(updatedVertices, fullPoint...)
		updatedVertices = append(updatedVertices, froms[2]...)

		// Add in the final segment
		updatedVertices = append(updatedVertices, froms[2]...)
		updatedVertices = append(updatedVertices, to...)
	}

	return updatedVertices
}

func makeBuffers(offset float32) ([]float32, uint32, uint32) {
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

	return vertices, VAO, VBO
}

func main() {
	title := "Koch snowflake"
	fmt.Println("Starting")

	window := glfwBoilerplate.InitGLFW(title,
		500, 500, false)
	defer glfw.Terminate()
	//gl.LineWidth(100.0)
	gl.Enable(gl.MULTISAMPLE) // Enable anti aliasing

	ourShader := shader.MakeShaders("koch.vs", "koch.fs")

	triangleSize := float32(.75)
	//triangleSize := float32(1.5)
	vertices, VAO, VBO := makeBuffers(triangleSize)
	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteVertexArrays(1, &VBO)

	lastTime := 0.0
	numFrames := 0.0

	maxIters := 5
	for !window.ShouldClose() {
		// Preframe
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		//fmt.Println(len(vertices))

		// Actually render the fractal
		ourShader.Use()
		ourShader.SetMat4("transform", mgl.Ident4())

		//scaleF := float32(
		//	(math.Sin(float64(glfw.GetTime()*2.0)) * 5) + 10)
		//ourShader.SetMat4("transform", mgl.Scale3D(
		//	scaleF, scaleF, 0).Mul4(
		//	mgl.Translate3D(0, -1.5, 0)))

		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.LINE_LOOP, 0, int32(len(vertices)/pointsPerVertex))
		//gl.DrawArrays(gl.POINTS, 0, int32(len(vertices)/pointsPerVertex))
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()

		if maxIters > 0 {
			// Update triangle and VBO
			vertices = updateVertices(vertices)
			gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
			gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4,
				gl.Ptr(vertices), gl.STATIC_DRAW)
			gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		}

		maxIters -= 1
		time.Sleep(0 * time.Second)
	}
}
