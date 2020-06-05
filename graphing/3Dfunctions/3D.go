package main

import (
	"fmt"
	//	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/nicholasblaskey/go-learn-opengl/includes/camera"
	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"

	"github.com/nicholasblaskey/animations/glfwBoilerplate"
)

const windowWidth = 1280
const windowHeight = 720

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

type twoVarFunc func(x, y float32) float32

func init() {
	runtime.LockOSThread()
}

func getPositions(numX, numY int) []mgl.Vec2 {
	translations := []mgl.Vec2{}
	xOffset := 1.0 / float32(numX)
	yOffset := 1.0 / float32(numY)
	for y := -numY; y < numY; y += 2 {
		for x := -numX; x < numX; x += 2 {
			translations = append(translations,
				mgl.Vec2{float32(x)/float32(numX) + xOffset,
					float32(y)/float32(numY) + yOffset})
		}
	}
	return translations
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

func makeFunctionBuffs(params graphParams, fx twoVarFunc,
	col mgl.Vec3) (uint32, uint32, int32) {

	numX := int(params.xRange[1] - params.xRange[0]/params.gridSpacing)
	numY := int(params.yRange[1] - params.yRange[0]/params.gridSpacing)
	pos := getPositions(numX, numY)
	xOff := 1.0 / float32(numX)
	yOff := 1.0 / float32(numY)

	getZ := func(xPos, yPos float32) float32 {
		xVal := (xPos / (2.0 - params.xBoarder*2) *
			(params.xRange[1] - params.xRange[0]))
		yVal := (yPos / (2.0 - params.yBoarder*2) *
			(params.yRange[1] - params.yRange[0]))
		return (2.0 - params.zBoarder*2) * fx(xVal, yVal) /
			(params.zRange[1] - params.zRange[0])
	}

	vertices := []float32{}
	for i := 0; i < len(pos); i++ {
		vertices = append(vertices,
			-xOff+pos[i][0], getZ(-xOff+pos[i][0], yOff+pos[i][1]), yOff+pos[i][1],
			col[0], col[1], col[2],
			xOff+pos[i][0], getZ(xOff+pos[i][0], -yOff+pos[i][1]), -yOff+pos[i][1],
			col[0], col[1], col[2],
			-xOff+pos[i][0], getZ(-xOff+pos[i][0], -yOff+pos[i][1]), -yOff+pos[i][1],
			col[0], col[1], col[2],

			-xOff+pos[i][0], getZ(-xOff+pos[i][0], yOff+pos[i][1]), yOff+pos[i][1],
			col[0], col[1], col[2],
			xOff+pos[i][0], getZ(xOff+pos[i][0], -yOff+pos[i][1]), -yOff+pos[i][1],
			col[0], col[1], col[2],
			xOff+pos[i][0], getZ(xOff+pos[i][0], yOff+pos[i][1]), yOff+pos[i][1],
			col[0], col[1], col[2],
		)
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

func makeAsymptoteBuffs(params graphParams) (
	[]uint32, []int32, uint32, bool) {

	funcVAOs := []uint32{}
	funcVertexCounts := []int32{}

	params.xRange = mgl.Vec2{-100, 100}
	params.yRange = mgl.Vec2{-100, 100}
	params.zRange = mgl.Vec2{-1, 1}
	params.gridSpacing = 1.0

	funcVAO, _, funcVertexCount := makeFunctionBuffs(params,
		func(x, y float32) float32 {
			return 1/x + 0.00001 + 1/y + 0.0001
		},
		mgl.Vec3{0.20, 0.7, 0.55})
	funcVAOs = append(funcVAOs, funcVAO)
	funcVertexCounts = append(funcVertexCounts, funcVertexCount)

	return funcVAOs, funcVertexCounts, gl.TRIANGLES, true
}

func makeRandPlane(params graphParams) (
	[]uint32, []int32, uint32, bool) {
	funcVAOs := []uint32{}
	funcVertexCounts := []int32{}

	rand.Seed(100010)

	count := 100
	for i := -count; i <= count; i++ {
		funcVAO, _, funcVertexCount := makeFunctionBuffs(params,
			func(x, y float32) float32 {
				return x + y + rand.Float32()
			},
			mgl.Vec3{0.5 + float32(i)/float32(count),
				0.5 + float32(i)/float32(count),
				0.5 + float32(-i)/float32(count)})
		funcVAOs = append(funcVAOs, funcVAO)
		funcVertexCounts = append(funcVertexCounts, funcVertexCount)
	}
	return funcVAOs, funcVertexCounts, gl.TRIANGLES, false
}

func makeRandColorPlane(params graphParams) (
	[]uint32, []int32, uint32, bool) {
	funcVAOs := []uint32{}
	funcVertexCounts := []int32{}

	rand.Seed(100010)

	count := 10
	for i := -count; i <= count; i++ {
		funcVAO, _, funcVertexCount := makeFunctionBuffs(params,
			func(x, y float32) float32 {
				return 33 * rand.Float32() * rand.Float32()
			},
			mgl.Vec3{rand.Float32(), rand.Float32(), rand.Float32()})
		funcVAOs = append(funcVAOs, funcVAO)
		funcVertexCounts = append(funcVertexCounts, funcVertexCount)
	}
	return funcVAOs, funcVertexCounts, gl.TRIANGLES, false
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
		windowWidth, windowHeight, true)
	defer glfw.Terminate()

	// Add in camera callbacks
	window.SetCursorPosCallback(glfw.CursorPosCallback(mouse_callback))
	window.SetScrollCallback(glfw.ScrollCallback(scroll_callback))
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetKeyCallback(keyCallback)

	//gl.Enable(gl.MULTISAMPLE) // Enable anti aliasing

	ourShader := shader.MakeShaders("3D.vs", "3D.fs")
	axisVAO, axisVBO, axisVertexCount := makeAxisBuffs(params)
	funcVAOs, funcVertexCounts, drawingType, isWireframe := makeAsymptoteBuffs(params)
	//funcVAOs, funcVertexCounts, drawingType, isWireframe := makeRandPlane(params)
	//funcVAOs, funcVertexCounts, drawingType, isWireframe := makeRandColorPlane(params)

	/*funcVAO, funcVBO, funcVertexCount := makeFunctionBuffs(params,
	func(x, y float32) float32 {
		return float32(math.Sin(float64(x)))*y +
			float32(math.Cos(float64(y)))*x
	},
	mgl.Vec3{0.3, 0.6, 0.3})
	*/

	defer gl.DeleteVertexArrays(1, &axisVAO)
	defer gl.DeleteVertexArrays(1, &axisVBO)
	//defer gl.DeleteVertexArrays(1, &funcVAO)
	//defer gl.DeleteVertexArrays(1, &funcVBO)

	lastTime := 0.0
	numFrames := 0.0

	if isWireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}
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
			float32(windowHeight)/windowWidth, 0.1, 1000.0)
		view := ourCamera.GetViewMatrix()
		ourShader.SetMat4("projection", projection)
		ourShader.SetMat4("view", view)
		ourShader.SetMat4("model", mgl.Ident4())

		// Draw Axis
		gl.BindVertexArray(axisVAO)
		gl.DrawArrays(gl.LINES, 0, axisVertexCount)
		gl.BindVertexArray(0)

		// Draw functions
		for i := 0; i < len(funcVAOs); i++ {
			gl.BindVertexArray(funcVAOs[i])
			gl.DrawArrays(drawingType, 0, funcVertexCounts[i])
			gl.BindVertexArray(0)
		}
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
