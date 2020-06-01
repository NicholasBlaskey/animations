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

		//fmt.Println(from)
		//fmt.Println(to)
		//fmt.Println("from above then to then endl\n")
		//fmt.Printf("from=(%f, %f),to=(%f, %f)\n",
		//	vertices[(i)%lenVertices], vertices[(i+1)%lenVertices],
		//	vertices[(i+5)%lenVertices], vertices[(i+5+1)%lenVertices],
		//)

		// TODO clean this up into a single for loop
		// Step 1 divide the segment into 3 equal parts
		ratio := float32(1.0 / 3.0)
		froms := [][]float32{}
		tos := [][]float32{}

		// First from is just the first from
		froms = append(froms, from)

		// Second from is 66% of from and 33% of to
		curVertex := []float32{}
		for j := 0; j < len(from); j++ {
			curVertex = append(curVertex, from[j]*ratio*2+to[j]*ratio*1)
		}
		froms = append(froms, curVertex)

		// Third from is 33% of from and 66% of to
		curVertex = []float32{}
		for j := 0; j < len(from); j++ {
			curVertex = append(curVertex, from[j]*ratio*1+to[j]*ratio*2)
		}
		froms = append(froms, curVertex)

		//fmt.Println(tos)
		tos = append(tos, froms[1], froms[2], to)

		/*
			//fmt.Println("not!hERE?")
			fmt.Println(from)
			fmt.Println(to)
			fmt.Println(froms)
			fmt.Println(tos)
			fmt.Println("from then to then froms then tos above\n")
		*/

		// Add in first segment
		for i := 0; i < len(froms[0]); i++ {
			updatedVertices = append(updatedVertices, froms[0][i])
		}
		for i := 0; i < len(froms[0]); i++ {
			updatedVertices = append(updatedVertices, tos[0][i])
		}

		// Get third triangle point
		mid := mgl.Vec2{(froms[1][0] + froms[2][0]) / 2.0,
			(froms[1][1] + froms[2][1]) / 2.0}
		orig := mgl.Vec2{(froms[1][0] - mid[0]), (froms[1][1] - mid[1])}
		orig.Mul(3 * float32(math.Sqrt(3)))
		transform := mgl.Rotate2D(mgl.DegToRad(90))
		topTri := mid.Add(transform.Mul2x1(orig))

		// This is bugged (also check with the thing)
		// TODO rewrite with matrices
		// Add in the middle two segments
		// https://stackoverflow.com/questions/50547068/creating-an-equilateral-triangle-for-given-two-points-in-the-plane-python
		// We need some way to decide which it is. One is to check if the vecs are othogonal which will def work within error?
		// The other which may or may not work is to see which one is further away from the origin.
		/*
			mid := []float32{(froms[1][0] + froms[2][0]) / 2.0,
				(froms[1][1] + froms[2][1]) / 2.0}
			origin := []float32{(froms[1][0] - mid[0]) * 3 * float32(math.Sqrt(3)),
				(froms[1][1] - mid[1]) * 3 * float32(math.Sqrt(3))}
			//topTri := []float32{mid[0] + origin[1], mid[1] - origin[0], 0, 0, 0}
			topTri := []float32{mid[0] - origin[1], mid[1] + origin[0], 0, 0, 0}
		*/

		for i := 0; i < len(froms[0]); i++ {
			updatedVertices = append(updatedVertices, froms[1][i])
		}
		updatedVertices = append(updatedVertices, topTri[0], topTri[1],
			froms[1][2]+froms[2][2], froms[1][3]+froms[2][3],
			froms[1][4]+froms[2][4])
		updatedVertices = append(updatedVertices, topTri[0], topTri[1],
			froms[1][2]+froms[2][2], froms[1][3]+froms[2][3],
			froms[1][4]+froms[2][4])
		for i := 0; i < len(froms[0]); i++ {
			updatedVertices = append(updatedVertices, froms[2][i])
		}

		// Add in the final segment
		for i := 0; i < len(froms[0]); i++ {
			updatedVertices = append(updatedVertices, froms[2][i])
		}
		for i := 0; i < len(froms[0]); i++ {
			updatedVertices = append(updatedVertices, tos[2][i])
		}

		//for i := 0; i < 3; i++ {
		// Copy the from over
		//	froms = append(froms, from)
		//	froms[i][0] *= 1 / float32(i) *
		//}

		// Step 2 draw an equilateral triangle that has the middle segment
		// from step 1 as its base and points outward.

		// Step 3 remove the line segment that is the base of the triangle
		// from step 2
		//updatedVertices = append(updatedVertices,
		//	froms[0][0], froms[0][1], froms[0][2], froms[0][3], froms[0][4],
		//	tos[0][0], tos[0][1], tos[0][2], tos[0][3], tos[0][4],
		//	froms[1][0], froms[1][1], froms[1][2], froms[1][3], froms[1][4],
		//	tos[1][0], tos[1][1], tos[1][2], tos[1][3], tos[1][4],
		//	froms[2][0], froms[2][1], froms[2][2], froms[2][3], froms[2][4],
		//	tos[2][0], tos[2][1], tos[2][2], tos[2][3], tos[2][4])

	}

	//fmt.Println(updatedVertices)
	//fmt.Println()
	//panic("LOL")
	//return vertices
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
		800, 600, false)
	defer glfw.Terminate()
	gl.LineWidth(2.0)

	ourShader := shader.MakeShaders("koch.vs", "koch.fs")

	triangleSize := float32(0.75)
	vertices, VAO, VBO := makeBuffers(triangleSize)
	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteVertexArrays(1, &VBO)

	lastTime := 0.0
	numFrames := 0.0

	maxIters := 7
	for !window.ShouldClose() {
		// Preframe
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		//fmt.Println(len(vertices))

		if maxIters > 0 {
			// Update triangle and VBO
			vertices = updateVertices(vertices)
			gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
			gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4,
				gl.Ptr(vertices), gl.STATIC_DRAW)
			gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		}

		// Actually render the fractal
		ourShader.Use()
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.LINE_LOOP, 0, int32(len(vertices)/pointsPerVertex))
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()

		maxIters -= 1
		time.Sleep(0 * time.Second)
	}
}
