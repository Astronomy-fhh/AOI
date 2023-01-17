package kit

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"image/color"
	"math/rand"
	"sync"
	"time"
)

var DrawCon *DrawContainer

type DrawContainer struct {
	sync.RWMutex
	scopes []*rect
	window *sdlcanvas.Window
	canvas *canvas.Canvas
	frameFunc func()
}

const canvasX = 1024
const canvasY = 1024

func init() {
	window, cvs, err := sdlcanvas.CreateWindow(canvasX, canvasY, "AOI")
	if err != nil {
		panic(err)
	}

	DrawCon = &DrawContainer{
		scopes: make([]*rect, 0),
		window: window,
		canvas: cvs,
	}

	window.Window.SetTitle("AOI")
	window.Window.SetBordered(true)


	DrawCon.frameFunc = func() {
		DrawCon.RLock()
		cvs.SetFillStyle(color.RGBA{R: 255, G: 255, B: 255, A: 255})
		cvs.FillRect(0, 0, float64(cvs.Width()), float64(cvs.Height()))

		cvs.SetStrokeStyle(color.RGBA{R: 0, G: 0, B: 0, A: 255})
		for _, scope := range DrawCon.scopes {
			cvs.StrokeRect(float64(scope.x), float64(scope.y), float64(scope.w), float64(scope.h))
		}
		DrawCon.RUnlock()
	}
}

func CvsStart() {
	defer DrawCon.window.Close()
	DrawCon.window.MainLoop(DrawCon.frameFunc)
}

func newRect(x, y, w, h int) *rect {
	return &rect{x, y, w, h}
}

type rect struct {
	x int
	y int
	w int
	h int
}

func test() {
	for {
		time.Sleep(time.Second)
		sx := rand.Int31n(500)
		sy := rand.Int31n(500)
		ex := sx + rand.Int31n(500)
		ey := sy + rand.Int31n(500)
		CvsAddScope(int(sx), int(sy), int(ex), int(ey))
	}
}

func (c *DrawContainer)AddScope(startX, startY, endX, endY int) {
	DrawCon.Lock()
	defer DrawCon.Unlock()
	rect := newRect(startX, startY, endX-startX, endY-startY)
	c.scopes = append(c.scopes, rect)
}


func CvsAddScope(startX, startY, endX, endY int) {
	DrawCon.Lock()
	defer DrawCon.Unlock()
	rect := newRect(startX, startY, endX-startX, endY-startY)
	DrawCon.scopes = append(DrawCon.scopes, rect)
}
