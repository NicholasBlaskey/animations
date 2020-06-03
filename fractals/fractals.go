package fractals

import (
	"math"
	"sync"

	mgl "github.com/go-gl/mathgl/mgl32"
)

func UpdateKoch(vertices []float32) []float32 {
	var wg sync.WaitGroup

	flattenedVertices := []float32{}
	verticesOrdered := make([](*[]float32), len(vertices)/5)

	// For each line segment
	for i := 0; i < len(vertices); i += 5 {
		wg.Add(1)
		curVertex := &[]float32{}
		go workerKoch(&wg, vertices, curVertex, i)
		verticesOrdered[i/5] = curVertex
	}
	wg.Wait()

	for i := 0; i < len(verticesOrdered); i++ {
		flattenedVertices = append(flattenedVertices, *verticesOrdered[i]...)
	}
	return flattenedVertices
}

func workerKoch(wg *sync.WaitGroup, vertices []float32,
	updatedVertices *[]float32, segID int) {

	defer wg.Done()

	pointsPerVertex := 5
	from := vertices[segID : segID+pointsPerVertex]
	startingIndex := (segID + pointsPerVertex) % len(vertices)
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

	// Add in first segment
	*updatedVertices = append(*updatedVertices, froms[0]...)
	*updatedVertices = append(*updatedVertices, froms[1]...)
	// Add in the triangle segments
	*updatedVertices = append(*updatedVertices, froms[1]...)
	*updatedVertices = append(*updatedVertices, fullPoint...)
	*updatedVertices = append(*updatedVertices, fullPoint...)
	*updatedVertices = append(*updatedVertices, froms[2]...)
	// Add in the final segment
	*updatedVertices = append(*updatedVertices, froms[2]...)
	*updatedVertices = append(*updatedVertices, to...)
}
