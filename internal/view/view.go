package view

import (
	"AOI/internal/g"
	p "AOI/internal/player"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"math"
	"sync"
)

var DrawCon *DrawContainer

type DrawContainer struct {
	sync.RWMutex
	scopes    []*rect
	players   []*p.Player
	window    *sdlcanvas.Window
	canvas    *canvas.Canvas
	frameFunc func()
}

func init() {
	window, cvs, err := sdlcanvas.CreateWindow(g.MapX, g.MapY, g.ViewTitle)
	if err != nil {
		panic(err)
	}

	DrawCon = &DrawContainer{
		scopes:  make([]*rect, 0),
		players: make([]*p.Player, 0),
		window:  window,
		canvas:  cvs,
	}

	window.Window.SetTitle(g.ViewTitle)
	window.Window.SetBordered(true)

	DrawCon.frameFunc = func() {
		DrawCon.RLock()
		cvs.SetFillStyle(g.ViewBackgroundColor)
		cvs.FillRect(0, 0, float64(cvs.Width()), float64(cvs.Height()))

		cvs.SetStrokeStyle(g.ViewScopeLineColor)
		for _, scope := range DrawCon.scopes {
			cvs.StrokeRect(scope.x, scope.y, scope.w, scope.h)
		}

		cvs.SetFillStyle(g.ViewScopePlayerColor)
		for _, p := range DrawCon.players {
			p.RLock()
			cvs.BeginPath()
			cvs.Arc(p.X, p.Y, g.PlayerRadius, 0, math.Pi*2, false)
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

func (c *DrawContainer) AddScope(startX, startY, endX, endY float64) {
	c.Lock()
	defer c.Unlock()
	rect := newRect(startX, startY, endX-startX, endY-startY)
	c.scopes = append(c.scopes, rect)
}

func (c *DrawContainer) RegisterPlayer(p *p.Player) {
	c.Lock()
	defer c.Unlock()
	c.players = append(c.players, p)
}
