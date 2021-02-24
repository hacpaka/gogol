package gogol

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

func glMakeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func glPrepareTriangles(columns, rows uint, x, y int) []float32 {
	var (
		triangles = []float32{
			-1 , 1 * (2 / float32(rows) - 1), 0,
			1 * (2 / float32(columns) - 1), -1, 0,
			-1, -1, 0,

			-1 , 1 * (2 / float32(rows) - 1), 0,
			1 * (2 / float32(columns) - 1), 1 * (2 / float32(rows) - 1), 0,
			1 * (2 / float32(columns) - 1), -1, 0,
		}
	)

	offset := []float32 {
		0.0,

		2 / float32(columns) * float32(x),
		2 / float32(rows) * float32(y),
	}

	for i := 0; i < len(triangles); i++ {
		triangles[i] += offset[(i + 1) % 3]
	}

	return triangles
}
