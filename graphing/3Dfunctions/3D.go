package main

import (
	"fmt"
	//	"math"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/nicholasblaskey/go-learn-opengl/includes/camera"
	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/animations/glfwBoilerplate"
)

const windowWidth = 500
const windowHeight = 500

// Camera
var ourCamera camera.Camera = camera.NewCamera(
	0.0, 0.0, 3.0, // pos xyz
	0.0, 1.0, 0.0, // up xyz
	-90.0, 0.0, // Yaw and pitch
	80.0, 45.0, 0.1) // Speed, zoom, and mouse sensitivity
var firstMouse bool = true
var lastX float32 = windowWidth / 2
var lastY float32 = windowHeight / 2

// Timing
var deltaTime float32 = 0.0
var lastFrame float32 = 0.0

// Controls
var heldW bool = false
var heldA bool = false
var heldS bool = false
var heldD bool = false

const pointsPerVertex = 6

type graphParams struct {
	xBoarder    float32
	yBoarder    float32
	zBoarder    float32
	xRange      mgl.Vec2
	yRange      mgl.Vec2
	zRange      mgl.Vec2
	gridSpacing float32
	//	xAxisColor mgl.Vec3
	//	yAxisColor mgl.Vec3
}

func init() {
	runtime.LockOSThread()
}

func makeAxisBuffs(params graphParams) (uint32, uint32, int32) {
	vertices := []float32{
		// Positions         // Color coords
		1 - params.xBoarder, 0.0, 0.0, 1, 0.0, 0.0,
		-1 + params.xBoarder, 0.0, 0.0, 1, 0.0, 0.0,
		0.0, 1 - params.yBoarder, 0.0, 0.0, 1.0, 0.0,
		0.0, -1 + params.yBoarder, 0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 1 - params.zBoarder, 0.0, 0.0, 1.0,
		0.0, 0.0, -1 + params.zBoarder, 0.0, 0.0, 1.0,
	}

	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4,
		gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, pointsPerVertex*4,
		gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, pointsPerVertex*4,
		gl.PtrOffset(3*4))

	return VAO, VBO, int32(len(vertices) / pointsPerVertex)
}

func main() {
	title := "3D graping"
	fmt.Println("Starting")

	params := graphParams{
		xBoarder:    0.1,
		yBoarder:    0.1,
		zBoarder:    0.1,
		xRange:      mgl.Vec2{-5, 5},
		yRange:      mgl.Vec2{-5, 5},
		zRange:      mgl.Vec2{-5, 5},
		gridSpacing: 0.1,
	}

	window := glfwBoilerplate.InitGLFW(title,
		windowWidth, windowHeight, false)
	defer glfw.Terminate()

	// Add in camera callbacks
	window.SetCursorPosCallback(glfw.CursorPosCallback(mouse_callback))
	window.SetScrollCallback(glfw.ScrollCallback(scroll_callback))
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetKeyCallback(keyCallback)

	//gl.Enable(gl.MULTISAMPLE) // Enable anti aliasing

	ourShader := shader.MakeShaders("3D.vs", "3D.fs")
	axisVAO, axisVBO, axisVertexCount := makeAxisBuffs(params)

	defer gl.DeleteVertexArrays(1, &axisVAO)
	defer gl.DeleteVertexArrays(1, &axisVBO)
	//defer gl.DeleteVertexArrays(1, &funcVAO)
	//defer gl.DeleteVertexArrays(1, &funcVBO)

	lastTime := 0.0
	numFrames := 0.0
	for !window.ShouldClose() {
		// Pre frame logic
		currentFrame := float32(glfw.GetTime())
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame
		glfw.PollEvents()

		// Preframe
		lastTime, numFrames = glfwBoilerplate.DisplayFrameRate(
			window, title, numFrames, lastTime)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		// View / projection transformations
		ourShader.Use()
		projection := mgl.Perspective(mgl.DegToRad(ourCamera.Zoom),
			float32(windowHeight)/windowWidth, 0.1, 100.0)
		view := ourCamera.GetViewMatrix()
		ourShader.SetMat4("projection", projection)
		ourShader.SetMat4("view", view)
		ourShader.SetMat4("model", mgl.Ident4())

		// Draw Axis
		gl.BindVertexArray(axisVAO)
		gl.DrawArrays(gl.LINES, 0, axisVertexCount)
		gl.BindVertexArray(0)

		// Draw functions
		//gl.BindVertexArray(funcVAOs[i])
		//gl.DrawArray(gl.POINTS, 0, funcVertexCount)
		//gl.BindVertexArray(funcVAO)
		//gl.DrawArrays(gl.LINE_STRIP, 0, funcVertexCount)
		//gl.BindVertexArray(0)
		//}

		window.SwapBuffers()

		time.Sleep(0 * time.Millisecond)
	}
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int,
	action glfw.Action, mods glfw.ModifierKey) {

	// Escape closes window
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}

	if key == glfw.KeyW && action == glfw.Press || heldW {
		ourCamera.ProcessKeyboard(camera.FORWARD, deltaTime)
		heldW = true
	}
	if key == glfw.KeyS && action == glfw.Press || heldS {
		ourCamera.ProcessKeyboard(camera.BACKWARD, deltaTime)
		heldS = true
	}
	if key == glfw.KeyA && action == glfw.Press || heldA {
		ourCamera.ProcessKeyboard(camera.LEFT, deltaTime)
		heldA = true
	}
	if key == glfw.KeyD && action == glfw.Press || heldD {
		ourCamera.ProcessKeyboard(camera.RIGHT, deltaTime)
		heldD = true
	}

	if key == glfw.KeyW && action == glfw.Release {
		heldW = false
	}
	if key == glfw.KeyS && action == glfw.Release {
		heldS = false
	}
	if key == glfw.KeyA && action == glfw.Release {
		heldA = false
	}
	if key == glfw.KeyD && action == glfw.Release {
		heldD = false
	}
}

func mouse_callback(w *glfw.Window, xPos float64, yPos float64) {
	if firstMouse {
		lastX = float32(xPos)
		lastY = float32(yPos)
		firstMouse = false
	}

	xOffset := float32(xPos) - lastX
	// Reversed due to y coords go from bot up
	yOffset := lastY - float32(yPos)

	lastX = float32(xPos)
	lastY = float32(yPos)

	ourCamera.ProcessMouseMovement(xOffset, yOffset, true)
}

func scroll_callback(w *glfw.Window, xOffset float64, yOffset float64) {
	ourCamera.ProcessMouseScroll(float32(yOffset))
}

func framebuffer_size_callback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}
