package gogol

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"log"
)

func glInit() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	log.Printf("OpenGL version: %v", gl.GoStr(gl.GetString(gl.VERSION)))

	vertexShader, err := glCompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := glCompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	pointer := gl.CreateProgram()
	gl.AttachShader(pointer, vertexShader)
	gl.AttachShader(pointer, fragmentShader)
	gl.LinkProgram(pointer)

	return pointer
}

