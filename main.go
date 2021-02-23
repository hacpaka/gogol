package gogol

import (
	"errors"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	Size uint = 20
	Duration uint = 10
	MinWidth uint = 200
	MinHeight uint = 200
)

type Engine struct {
	handler uint32
	units [][]*TUnit
}

func (e *Engine) Init (action func([][]*TUnit) error, width uint, height uint) error {
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

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Decorated, glfw.False)
	glfw.WindowHint(glfw.Focused, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(int(width), int(height), "Test window", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.Maximize()

	e.handler = glInit()
	e.units = Units(width / Size, height / Size, 1500)

	timeMark := time.Now()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(e.handler)

		if err := action(e.units); err != nil {
			glfw.Terminate()
		} else {
			for x := range e.units {
				for y := range e.units[x] {
					if e.units[x][y].Life > 0 {
						e.draw(e.units[x][y].data, e.units[x][y].color)
					}
				}
			}
		}

		glfw.PollEvents()
		window.SwapBuffers()

		time.Sleep(time.Second/time.Duration(Duration) - time.Since(timeMark))
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