package aoi

import (
	"AOI/internal/g"
	"fmt"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"math"
	"sync"
)

var WinView *WindowView

type WindowView struct {
	sync.RWMutex
	scopes    []*rect
	players   []*Player
	window    *sdlcanvas.Window
	canvas    *canvas.Canvas
	mouseX    int
	mouseY    int
	frameFunc func()
}

func init() {
	window, cvs, err := sdlcanvas.CreateWindow(g.MapX, g.MapY, g.ViewTitle)
	if err != nil {
		panic(err)
	}

	WinView = &WindowView{
		scopes:  make([]*rect, 0),
		players: make([]*Player, 0),
		window:  window,
		canvas:  cvs,
		mouseX:  g.MapX / 2,
		mouseY:  g.MapY / 2,
	}

	window.Window.SetTitle(g.ViewTitle)
	window.Window.SetBordered(true)

	window.MouseMove = func(x, y int) {
		WinView.Lock()
		WinView.mouseX, WinView.mouseY = x, y
		WinView.Unlock()
	}
	window.MouseWheel = func(x, y int) {
		WinView.Lock()
		updateSelectScopeRange(y)
		WinView.Unlock()
	}

	WinView.frameFunc = func() {
		WinView.RLock()
		cvs.SetFillStyle(g.ViewBackgroundColor)
		cvs.FillRect(0, 0, float64(cvs.Width()), float64(cvs.Height()))

		cvs.SetStrokeStyle(g.ViewScopeLineColor)
		for _, scope := range WinView.scopes {
			cvs.StrokeRect(scope.x, scope.y, scope.w, scope.h)
		}

		cvs.SetStrokeStyle(g.ViewSelectScopeLineColor)
		cvs.StrokeRect(float64(WinView.mouseX)-g.ViewSelectScopeLineWidth/2, float64(WinView.mouseY)-g.ViewSelectScopeLineHeight/2, g.ViewSelectScopeLineWidth, g.ViewSelectScopeLineHeight)

		cvs.SetFillStyle(g.ViewScopePlayerColor)
		for _, p := range WinView.players {
			p.RLock()
			cvs.BeginPath()
			cvs.Arc(p.X, p.Y, g.PlayerRadius, 0, math.Pi*2, false)
			cvs.Fill()
			p.RUnlock()
		}

		var fontLineIdx float64
		getFontLineHeight := func() float64 {
			fontLineIdx++
			return g.ViewFontLineHeight * fontLineIdx
		}

		cvs.SetFont("font/Roboto-Light.ttf", g.ViewFontSize)
		cvs.SetFillStyle(g.ViewFontColor1)
		cvs.FillText(fmt.Sprintf("FPS:%.2f", window.FPS()), 0, getFontLineHeight())
		cvs.SetFillStyle(g.ViewFontColor2)
		cvs.FillText(fmt.Sprintf("MapSize [%d,%d] , PlayerNum:%d", g.MapX, g.MapY, len(WinView.players)), 1, getFontLineHeight())

		num := getScopePlayerNum(
			float64(WinView.mouseX)-g.ViewSelectScopeLineWidth/2,
			float64(WinView.mouseY)-g.ViewSelectScopeLineHeight/2,
			float64(WinView.mouseX)+g.ViewSelectScopeLineWidth/2,
			float64(WinView.mouseY)+g.ViewSelectScopeLineHeight/2,
		)

		cvs.FillText(fmt.Sprintf("SelectView [%.0f,%.0f  %.0f,%.0f] , PlayerNum:%d",
			float64(WinView.mouseX)-g.ViewSelectScopeLineWidth/2,
			float64(WinView.mouseY)-g.ViewSelectScopeLineHeight/2,
			float64(WinView.mouseX)+g.ViewSelectScopeLineWidth/2,
			float64(WinView.mouseY)+g.ViewSelectScopeLineHeight/2,
			num), 0, getFontLineHeight(),
		)

		WinView.RUnlock()
	}
}

func CvsStart() {
	defer WinView.window.Close()
	WinView.window.MainLoop(WinView.frameFunc)
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

func (c *WindowView) RegisterScope(startX, startY, endX, endY float64) {
	c.Lock()
	defer c.Unlock()
	rect := newRect(startX, startY, endX-startX, endY-startY)
	c.scopes = append(c.scopes, rect)
}

func (c *WindowView) RegisterPlayer(p *Player) {
	c.Lock()
	defer c.Unlock()
	c.players = append(c.players, p)
}

func updateSelectScopeRange(y int) {
	y *= g.ViewMouseWheelSensitivity
	g.ViewSelectScopeLineHeight += float64(y)
	if g.ViewSelectScopeLineHeight > g.ViewSelectScopeLineHeightMax {
		g.ViewSelectScopeLineHeight = g.ViewSelectScopeLineHeightMax
	}
	if g.ViewSelectScopeLineHeight < g.ViewSelectScopeLineHeightMin {
		g.ViewSelectScopeLineHeight = g.ViewSelectScopeLineHeightMin
	}

	g.ViewSelectScopeLineWidth += float64(y)
	if g.ViewSelectScopeLineWidth > g.ViewSelectScopeLineWidthMax {
		g.ViewSelectScopeLineWidth = g.ViewSelectScopeLineWidthMax
	}
	if g.ViewSelectScopeLineWidth < g.ViewSelectScopeLineWidthMin {
		g.ViewSelectScopeLineWidth = g.ViewSelectScopeLineWidthMin
	}
}

func getScopePlayerNum(startX, startY, endX, endY float64) int {
	players := QTree.Range(startX, startY, endX, endY)
	return len(players)
}
