package kit

import (
	"AOI/player"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"image/color"
	"math"
	"sync"
)

var DrawCon *DrawContainer

type DrawContainer struct {
	sync.RWMutex
	scopes []*rect
	players []*player.Player
	window *sdlcanvas.Window
	canvas *canvas.Canvas
	frameFunc func()
}

const canvasX = 1024
const canvasY = 512

func init() {
	window, cvs, err := sdlcanvas.CreateWindow(canvasX, canvasY, "AOI")
	if err != nil {
		panic(err)
	}

	DrawCon = &DrawContainer{
		scopes: make([]*rect, 0),
		players: make([]*player.Player, 0),
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
			cvs.StrokeRect(scope.x, scope.y, scope.w, scope.h)
		}

		cvs.SetFillStyle("#DA49D3")
		for _, p := range DrawCon.players {
			p.RLock()
			cvs.BeginPath()
			cvs.Arc(p.X, p.Y, 10, 0, math.Pi*2, false)
			cvs.Fill()
			p.RUnlock()
		}

		DrawCon.RUnlock()
	}
}

func CvsStart() {
	defer DrawCon.window.Close()
	DrawCon.window.MainLoop(DrawCon.frameFunc)
}

func newRect(x, y, w, h float64) *rect {
	return &rect{x, y, w, h}
}

type rect struct {
	x float64
	y float64
	w float64
	h float64
}

func (c *DrawContainer)AddScope(startX, startY, endX, endY float64) {
	c.Lock()
	defer c.Unlock()
	rect := newRect(startX, startY, endX-startX, endY-startY)
	c.scopes = append(c.scopes, rect)
}


func (c *DrawContainer)RegisterPlayer(p *player.Player) {
	c.Lock()
	defer c.Unlock()
	c.players = append(c.players, p)
}

