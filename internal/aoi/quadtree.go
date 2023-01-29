package aoi

import (
	"AOI/internal/g"
	"math/rand"
	"sync"
	"time"
)

var QTree *quadTree

func Start() {
	quadTreeOpts := NewQuadTreeOpts(g.MaxSize, g.MaxSplit, float64(g.MapX), float64(g.MapY))
	QTree = NewQuadTree(quadTreeOpts)
	QTree.drawScope()
}

func StartEnterTest() {
	go func() {
		var uid = 0
		for {
			time.Sleep(time.Second / 100)
			x := rand.Int31n(int32(g.MapX))
			y := rand.Int31n(int32(g.MapY))
			uid++
			if uid > g.TestMaxEnterPlayer {
				break
			}
			p := NewPlayer(uid, float64(x), float64(y))
			QTree.Add(p)
			WinView.RegisterPlayer(p)
		}
	}()
}

func NewQuadTreeOpts(maxSize, maxSplit int, maxX, maxY float64) *quadTreeOpts {
	opts := &quadTreeOpts{
		MaxSize:  maxSize,
		MaxSplit: maxSplit,
		MaxX:     maxX,
		MaxY:     maxY,
	}
	return opts
}

type quadTreeOpts struct {
	MaxSize  int
	MaxSplit int
	MaxX     float64
	MaxY     float64
}

func NewQuadTree(opts *quadTreeOpts) *quadTree {
	qt := &quadTree{
		items:        make(map[int]*Player),
		startX:       0,
		endX:         opts.MaxX,
		startY:       0,
		endY:         opts.MaxY,
		quadTreeOpts: opts,
	}
	return qt
}

type quadTree struct {
	sync.RWMutex
	*quadTreeOpts
	depth     int
	size      int
	items     map[int]*Player
	startX    float64
	endX      float64
	startY    float64
	endY      float64
	leftUp    *quadTree
	leftDown  *quadTree
	rightUp   *quadTree
	rightDown *quadTree
}

func (q *quadTree) PlayerMove(ox,oy float64, p *Player) {
	q.Lock()
	defer q.Unlock()
	//q.Del(ox,oy,p)
	//q.Add(p)
}


func (q *quadTree) Get(startX, startY, endX, endY float64) {
	q.Range(startX, startY, endX, endY)
}

func (q *quadTree) Add(item *Player) {
	q.size++
	if q.isSplit() {
		// 往子节点加
		q.addToChild(item)
	} else {
		if q.canSplit() {
			// 分裂
			q.toSplit()
			q.addToChild(item)
		} else {
			// 自己加
			q.items[item.Id] = item
		}
	}
}

func (q *quadTree) Del(ox,oy float64, p *Player) {
	q.size--
	if q.isSplit() {
		q.getOpNode(ox, oy).Del(ox,oy,p)
	} else {
		delete(q.items, p.Id)
	}
}

func (q *quadTree) Range(startX, startY, endX, endY float64) []*Player {
	q.Lock()
	defer q.Unlock()
	res := make([]*Player, 0)
	q.RangeRecursive(startX, startY, endX, endY, &res)
	return res
}

func (q *quadTree) RangeRecursive(startX, startY, endX, endY float64, found *[]*Player) {
	if !q.collision(startX, startY, endX, endY) {
		return
	}
	if q.isSplit() {
		q.leftUp.RangeRecursive(startX, startY, endX, endY, found)
		q.leftDown.RangeRecursive(startX, startY, endX, endY, found)
		q.rightUp.RangeRecursive(startX, startY, endX, endY, found)
		q.rightDown.RangeRecursive(startX, startY, endX, endY, found)
	} else {
		q.RangeItems(startX, startY, endX, endY, found)
	}
}

func (q *quadTree) RangeItems(startX, startY, endX, endY float64, found *[]*Player) {
	for _, item := range q.items {
		item.RLock()
		if item.X >= startX && item.X < endX && item.Y >= startY && item.Y < endY {
			*found = append(*found, item)
		}
		item.RUnlock()
	}
}

func (q *quadTree) collision(startX, startY, endX, endY float64) bool {
	return !(q.endX <= startX || q.startX > endX || q.endY <= startY || q.startY > endY)
}

func (q *quadTree) addToChild(p *Player) {
	q.getOpNode(p.X, p.Y).Add(p)
}

func (q *quadTree) getOpNode(x, y float64) *quadTree {
	if x >= q.leftUp.endX {
		if y >= q.leftUp.endY {
			return q.rightDown
		}
		return q.rightUp
	} else {
		if y >= q.leftUp.endY {
			return q.leftDown
		}
		return q.leftUp
	}
}

func (q *quadTree) isSplit() bool {
	return q.leftUp != nil
}

func (q *quadTree) canSplit() bool {
	return q.size > q.MaxSize && q.depth < q.MaxSplit
}

func (q *quadTree) newChild() *quadTree {
	return &quadTree{
		items:        make(map[int]*Player),
		depth:        q.depth + 1,
		quadTreeOpts: q.quadTreeOpts,
	}
}

func (q *quadTree) toSplit() {

	q.leftUp = q.newChild()
	q.leftDown = q.newChild()
	q.rightUp = q.newChild()
	q.rightDown = q.newChild()

	q.leftUp.startX = q.startX
	q.leftUp.endX = q.startX + (q.endX-q.startX)/2
	q.leftUp.startY = q.startY
	q.leftUp.endY = q.startY + (q.endY-q.startY)/2
	q.leftUp.drawScope()

	q.leftDown.startX = q.startX
	q.leftDown.endX = q.leftUp.endX
	q.leftDown.startY = q.leftUp.endY
	q.leftDown.endY = q.endY
	q.leftDown.drawScope()

	q.rightUp.startX = q.leftUp.endX
	q.rightUp.endX = q.endX
	q.rightUp.startY = q.startY
	q.rightUp.endY = q.startY + (q.endY-q.startY)/2
	q.rightUp.drawScope()

	q.rightDown.startX = q.leftUp.endX
	q.rightDown.endX = q.endX
	q.rightDown.startY = q.rightUp.endY
	q.rightDown.endY = q.endY
	q.rightDown.drawScope()

	for _, item := range q.items {
		q.addToChild(item)
	}
	q.items = make(map[int]*Player)

	println("分裂")
}

func (q *quadTree) drawScope() {
	WinView.RegisterScope(q.startX, q.startY, q.endX, q.endY)
}
