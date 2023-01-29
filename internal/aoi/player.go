package aoi

import (
	"AOI/internal/g"
	"math/rand"
	"sync"
	"time"
)

func NewPlayer(id int, x, y float64) *Player {
	p := &Player{
		Id:       id,
		X:        x,
		Y:        y,
		moveAttr: newMoveAttr(),
	}

	go p.move()
	return p
}

type Player struct {
	sync.RWMutex
	Id int
	X  float64
	Y  float64
	*moveAttr
}

func newMoveAttr() *moveAttr {
	// 置随机速度，使运动更具随机性
	return &moveAttr{
		randSpeed(),
		randSpeed(),
		newMoment(),
		0,
		false,
	}
}

type moveAttr struct {
	xSpeed   float64
	ySpeed   float64
	moveTime int64
	colTime  int64
	col      bool
}

func (p *Player) move() {
	if !g.PlayerMove {
		return
	}
	for {
		time.Sleep(g.PlayerMoveInternal)
		p.Lock()
		//ox, oy := p.X, p.Y
		t := newMoment()
		p.collision()
		p.X += p.xSpeed * float64(t-p.moveTime)
		p.Y += p.ySpeed * float64(t-p.moveTime)
		p.moveTime = t
		//QTree.PlayerMove(ox, oy, p)
		p.Unlock()
	}
}

func (p *Player) changeColXSpeed() {
	// 碰撞速度直接重置，这样球的运动更有随机性
	if p.xSpeed > 0 {
		p.xSpeed = -1 * randSpeed()
	} else {
		p.xSpeed = randSpeed()
	}
}

func (p *Player) changeColYSpeed() {
	if p.ySpeed > 0 {
		p.ySpeed = -1 * randSpeed()
	} else {
		p.ySpeed = randSpeed()
	}
}

func (p *Player) collision() {
	// 碰撞后需要归位，否则连续检测帧内会重复碰撞
	if p.X < g.PlayerRadius {
		p.changeColXSpeed()
		p.X = g.PlayerRadius
	} else if p.X > float64(g.MapX)-g.PlayerRadius {
		p.changeColXSpeed()
		p.X = float64(g.MapX) - g.PlayerRadius
	}

	if p.Y < g.PlayerRadius {
		p.changeColYSpeed()
		p.Y = g.PlayerRadius
	} else if p.Y > float64(g.MapY)-g.PlayerRadius {
		p.changeColYSpeed()
		p.Y = float64(g.MapY) - g.PlayerRadius
	}
}

func newMoment() int64 {
	return time.Now().UnixNano()
}

func randSpeed() float64 {
	return float64(rand.Intn(g.PlayerSpeedRandFactor)+1) * g.PlayerBaseSpeed
}
