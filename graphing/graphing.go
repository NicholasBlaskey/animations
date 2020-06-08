package graphing

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type SingleVarFunc func(x float32) float32

type Params2D struct {
	XBoarder   float32
	YBoarder   float32
	XRange     mgl.Vec2
	YRange     mgl.Vec2
	Dx         float32
	XAxisColor mgl.Vec3
	YAxisColor mgl.Vec3
}

const pointsPerVertex = 5

func MakeAxisBuffs(params Params2D) (uint32, uint32, int32) {
	xCol := params.XAxisColor
	yCol := params.YAxisColor

	vertices := []float32{
		// Positions         // Color coords
		1 - params.XBoarder, 0.0, xCol[0], xCol[1], xCol[2],
		-1 + params.XBoarder, 0.0, xCol[0], xCol[1], xCol[2],
		0.0, 1 - params.YBoarder, yCol[0], yCol[1], yCol[2],
		0.0, -1 + params.YBoarder, yCol[0], yCol[1], yCol[2],
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

func MakeFunctionBuffs(params Params2D, fx SingleVarFunc,
	col mgl.Vec3) (uint32, uint32, int32) {

	vertices := []float32{}
	for i := params.XRange[0]; i <= params.XRange[1]; i += params.Dx {
		vertices = append(vertices,
			(2.0-params.XBoarder*2)*
				(i/(params.XRange[1]-params.XRange[0])),
			(2.0-params.YBoarder*2)*
				(fx(i)/(params.YRange[1]-params.YRange[0])),
			col[0], col[1], col[2])
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
