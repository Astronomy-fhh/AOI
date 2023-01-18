package quadtree

import (
	"AOI/kit"
	"AOI/player"
	"math/rand"
	"sync"
	"time"
)

const MapX = 1024
const MapY = 512
const MaxSize = 16
const MaxSplit = 8
const MaxItem = 100

var AOI *quadTree

func Start()  {
	quadTreeOpts := NewQuadTreeOpts(MaxSize, MaxSplit, MapX - 1, MapY - 1)
	AOI = NewQuadTree(quadTreeOpts)
	AOI.drawScope()
}

func StartEnterTest()  {
	go func() {
		var uid = 0
		for  {
			time.Sleep(time.Second/100)
			x := rand.Int31n(MapX)
			y := rand.Int31n(MapY)
			uid ++
			if uid > MaxItem {
				break
			}
			p := player.NewPlayer(uid, float64(x), float64(y))
			AOI.Enter(p)
			kit.DrawCon.RegisterPlayer(p)
		}
	}()
}

func NewQuadTreeOpts(maxSize, maxSplit int , maxX, maxY float64)*quadTreeOpts {
	opts := &quadTreeOpts{
		MaxSize: maxSize,
		MaxSplit: maxSplit,
		MaxX: maxX,
		MaxY: maxY,
	}
	return opts
}

type quadTreeOpts struct {
	MaxSize int
	MaxSplit int
	MaxX float64
	MaxY float64
}

func NewQuadTree(opts *quadTreeOpts) *quadTree {
	qt := &quadTree{
		items:  make(map[int]*player.Player),
		startX: 0,
		endX:   opts.MaxX,
		startY: 0,
		endY:   opts.MaxY,
		quadTreeOpts: opts,
	}
	return qt
}

type quadTree struct {
	sync.RWMutex
	*quadTreeOpts
	depth     int
	size      int
	items     map[int]*player.Player
	startX    float64
	endX      float64
	startY    float64
	endY      float64
	leftUp    *quadTree
	leftDown  *quadTree
	rightUp   *quadTree
	rightDown *quadTree
}

func (q *quadTree) Enter(item *player.Player) {
	q.Add(item)
	println("进入+1")
}

func (q *quadTree) Leave(item *player.Player) {
	q.Del(item)
}

func (q *quadTree) Get(startX, startY, endX, endY float64) {
	q.Range(startX, startY, endX, endY)
}


func (q *quadTree) Add(item *player.Player) {
	q.Lock()
	defer q.Unlock()

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

func (q *quadTree) Del(p *player.Player) {
	q.Lock()
	defer q.Unlock()

	q.size--
	if q.isSplit() {
		q.getOpNode(p.X, p.Y).Del(p)
	} else {
		delete(q.items, p.Id)
	}
}

func (q *quadTree) Range(startX, startY, endX, endY float64)[]*player.Player {
	q.Lock()
	defer q.Unlock()
	res := make([]*player.Player, 0)
	q.RangeRecursive(startX, startY, endX, endY, res)
	return res
}

func (q *quadTree) RangeRecursive(startX, startY, endX, endY float64, found []*player.Player) {
	if !q.collision(startX, startY, endX, endY) {
		return
	}
	if q.isSplit() {
		q.leftUp.RangeRecursive(startX, startY, endX, endY, found)
		q.leftDown.RangeRecursive(startX, startY, endX, endY, found)
		q.rightUp.RangeRecursive(startX, startY, endX, endY, found)
		q.rightDown.RangeRecursive(startX, startY, endX, endY, found)
	}else{
		q.RangeItems(startX, startY, endX, endY, found)
	}
}

func (q *quadTree) RangeItems(startX, startY, endX, endY float64, found []*player.Player){
	q.Lock()
	defer q.Unlock()

	for _, item := range q.items {
		if item.X >= startX && item.X <= endX && item.Y >= startY && item.Y <= endY {
			found = append(found, item)
		}
	}
}

func (q *quadTree) collision(startX, startY, endX, endY float64)bool  {
	return !(q.endX < startX || q.startX > endX || q.endY < startY || q.startY > endY)
}

func (q *quadTree) addToChild(p *player.Player) {
	q.getOpNode(p.X, p.Y).Add(p)
}

func (q *quadTree) getOpNode(x,y float64)*quadTree {
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
		items: make(map[int]*player.Player),
		depth: q.depth + 1,
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
	q.leftDown.startY = q.leftUp.endY + 1
	q.leftDown.endY = q.endY
	q.leftDown.drawScope()


	q.rightUp.startX = q.leftUp.endX + 1
	q.rightUp.endX = q.endX
	q.rightUp.startY = q.startY
	q.rightUp.endY = q.startY + (q.endY-q.startY)/2
	q.rightUp.drawScope()


	q.rightDown.startX = q.leftUp.endX + 1
	q.rightDown.endX = q.endX
	q.rightDown.startY = q.rightUp.endY + 1
	q.rightDown.endY = q.endY
	q.rightDown.drawScope()


	for _, item := range q.items {
		q.addToChild(item)
	}
	q.items = make(map[int]*player.Player)

	println("分裂")
}

func (q *quadTree) drawScope() {
	println(q.startX, q.startY, q.endX, q.endY)
	kit.DrawCon.AddScope(q.startX, q.startY, q.endX, q.endY)
}
