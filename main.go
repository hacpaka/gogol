package gogol

import (
	"errors"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	MinWidth  = 200
	MinHeight = 200
)

type Engine struct {
	handler uint32

	Units [][]*TUnit
}

func (e *Engine) init (duration int64, width int, height int) error {
	if duration < 1 {
		errors.New("invalid uiDuration")
	}

	if width < MinWidth {
		errors.New("invalid uiWidth")
	}

	if height < MinHeight {
		errors.New("invalid uiHeight")
	}

	if err := glfw.Init(); err != nil {
		panic(err)
	}

	defer glfw.Terminate()
	e.handler = glInit()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Decorated, glfw.False)
	glfw.WindowHint(glfw.Focused, glfw.True)

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Test window", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.Maximize()

	timeMark := time.Now()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(e.handler)

		for x := range e.Units {
			for y := range e.Units[x] {
				if e.Units[x][y].life > 0 {
					e.draw(e.Units[x][y].data, e.Units[x][y].color)
				}
			}
		}

		glfw.PollEvents()
		window.SwapBuffers()

		time.Sleep(time.Second/time.Duration(duration) - time.Since(timeMark))
		timeMark = time.Now()
	}

	return nil
}

func (e *Engine) draw(data []float32, color TColor)  {
	loc := gl.GetUniformLocation(e.handler, gl.Str("un_color" + "\x00"))
	gl.Uniform3f(loc, float32(color.R), float32(color.G), float32(color.B))

	gl.BindVertexArray(glMakeVao(data))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data) / 3))
}