package player

import (
	"sync"
	"time"
)

func NewPlayer(id int, x, y float64) *Player {
	p := &Player{
		Id: id,
		X:  x,
		Y:  y,
	}
	go p.move()
	return p
}

type Player struct {
	sync.RWMutex
	Id int
	X  float64
	Y  float64
	moveAttr *moveAttr
}

type moveAttr struct {
	speed float64
	colTime    time.Duration
}

func (p *Player) move()  {
	for  {
		time.Sleep(time.Millisecond * 10)
		p.Lock()
		p.X += 2
		p.Y += 1
		p.Unlock()
	}
}