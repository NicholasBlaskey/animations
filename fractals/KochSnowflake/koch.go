// https://en.wikipedia.org/wiki/Koch_snowflake

package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/animations/fractals"
	"github.com/nicholasblaskey/animations/glfwBoilerplate"
)

const pointsPerVertex = 5

func init() {
	runtime.LockOSThread()
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
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, pointsPerVertex*4,
		gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, pointsPerVertex*4,
		gl.PtrOffset(2*4))

	return vertices, VAO, VBO
}

func drawSlowly(window *glfw.Window, title string,
	vertices []float32, ourShader shader.Shader, VAO, VBO uint32) {

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

		// Actually render the fractal
		ourShader.Use()
		ourShader.SetMat4("transform", mgl.Scale3D(0.5, 0.5, 0))
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.LINE_LOOP, 0, int32(len(vertices)/pointsPerVertex))
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()

		if maxIters > 0 {
			// Update triangle and VBO
			vertices = fractals.UpdateKoch(vertices)
			gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
			gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4,
				gl.Ptr(vertices), gl.STATIC_DRAW)
			gl.BindBuffer(gl.ARRAY_BUFFER, 0)
			maxIters -= 1
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func slowZoom(window *glfw.Window, title string,
	vertices []float32, ourShader shader.Shader, VAO, VBO uint32) {

	lastTime := 0.0
	numFrames := 0.0
	maxIters := 7
	iters := 0
	for !window.ShouldClose() {
		// Preframe
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		// Actually render the fractal
		ourShader.Use()
		ourShader.SetMat4("transform", mgl.Ident4())

		scaleF := float32(iters) / 100.0
		iters += 1
		//scaleF := float32(
		//	(math.Sin(float64(glfw.GetTime()*2.0)) * 5) + 10)
		ourShader.SetMat4("transform", mgl.Scale3D(
			scaleF, scaleF, 0).Mul4(
			mgl.Translate3D(0, -1.5, 0)))

		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.LINE_LOOP, 0, int32(len(vertices)/pointsPerVertex))
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()

		if maxIters > 0 && iters%100 == 0 {
			// Update triangle and VBO
			vertices = fractals.UpdateKoch(vertices)
			gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
			gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4,
				gl.Ptr(vertices), gl.STATIC_DRAW)
			gl.BindBuffer(gl.ARRAY_BUFFER, 0)
			maxIters -= 1
		}

		time.Sleep(0 * time.Second)
	}
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

	//triangleSize := float32(.75)
	triangleSize := float32(1.5)
	vertices, VAO, VBO := makeBuffers(triangleSize)
	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteVertexArrays(1, &VBO)

	drawSlowly(window, title, vertices, ourShader, VAO, VBO)
	//slowZoom(window, title, vertices, ourShader, VAO, VBO)
}
