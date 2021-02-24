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

type Color struct {
	R int
	G int
	B int
}

func (c *Color) Default () {
	c.R, c.G, c.B = 0, 0, 0
}

type Point struct {
	Color Color

	x uint
	y uint

	data  []float32
}

func (p *Point) Default () {
	p.Color, p.x, p.y = Color{0,0,0}, 0, 0
}

func (p *Point) Draw (location int32)  {
	gl.Uniform3f(location, float32(p.Color.R), float32(p.Color.G), float32(p.Color.B))

	gl.BindVertexArray(glMakeVao(p.data))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(p.data) / 3))
}

type Engine struct {
	handler uint32
	points  [][]*Point

	Columns uint
	Rows uint
}

func (e *Engine) Init (action func([][]*Point) error, width uint, height uint) error {
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
	e.Rows, e.Columns = width / Size, height / Size

	e.points = make([][]*Point, e.Columns)
	for x := range e.points {
		e.points[x] = make([]*Point, e.Rows)

		for y :=  range e.points[x] {
			e.points[x][y] = &Point{
				Color{255, 255, 255 },
				uint(x),
				uint(y),
				glPrepareTriangles(e.Columns, e.Rows, x, y),
			}
		}
	}

	timeMark := time.Now()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(e.handler)

		if err := action(e.points); err != nil {
			glfw.Terminate()
		} else {
			for x := range e.points {
				for y := range e.points[x] {
					e.draw(e.points[x][y])
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

func (e *Engine) draw(point *Point)  {
	loc := gl.GetUniformLocation(e.handler, gl.Str("un_color" + "\x00"))
	point.Draw(loc)
}