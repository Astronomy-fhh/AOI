package quadtree

import (
	"AOI/kit"
	"math/rand"
	"sync"
	"time"
)

const MapX = 1024
const MapY = 1024
const MaxSize = 1
const MaxSplit = 4


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
			time.Sleep(time.Second/10)
			x := rand.Int31n(MapX)
			y := rand.Int31n(MapY)
			uid ++
			AOI.Enter(&Item{uid, int(x), int(y)})
		}
	}()
}

func NewQuadTreeOpts(maxSize, maxSplit, maxX, maxY int)*quadTreeOpts {
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
	MaxX int
	MaxY int
}

func NewQuadTree(opts *quadTreeOpts) *quadTree {
	qt := &quadTree{
		items:  make(map[int]*Item),
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
	items     map[int]*Item
	startX    int
	endX      int
	startY    int
	endY      int
	leftUp    *quadTree
	leftDown  *quadTree
	rightUp   *quadTree
	rightDown *quadTree
}

type Item struct {
	id int
	x  int
	y  int
}

func (q *quadTree) Enter(item *Item) {
	q.Add(item)
	println("进入+1")
}

func (q *quadTree) Leave(item *Item) {
	q.Del(item)
}

func (q *quadTree) Get(startX, startY, endX, endY int) {
	q.Range(startX, startY, endX, endY)
}


func (q *quadTree) Add(item *Item) {
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
			q.items[item.id] = item
		}
	}
}

func (q *quadTree) Del(item *Item) {
	q.Lock()
	defer q.Unlock()

	q.size--
	if q.isSplit() {
		q.getOpNode(item.x, item.y).Del(item)
	} else {
		delete(q.items, item.id)
	}
}

func (q *quadTree) Range(startX, startY, endX, endY int)[]*Item {
	q.Lock()
	defer q.Unlock()
	res := make([]*Item, 0)
	q.RangeRecursive(startX, startY, endX, endY, res)
	return res
}

func (q *quadTree) RangeRecursive(startX, startY, endX, endY int, found []*Item) {
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

func (q *quadTree) RangeItems(startX, startY, endX, endY int, found []*Item){
	q.Lock()
	defer q.Unlock()

	for _, item := range q.items {
		if item.x >= startX && item.x <= endX && item.y >= startY && item.y <= endY {
			found = append(found, item)
		}
	}
}

func (q *quadTree) collision(startX, startY, endX, endY int)bool  {
	return !(q.endX < startX || q.startX > endX || q.endY < startY || q.startY > endY)
}

func (q *quadTree) addToChild(item *Item) {
	q.getOpNode(item.x, item.y).Add(item)
}

func (q *quadTree) getOpNode(x,y int)*quadTree {
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
		items: make(map[int]*Item),
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
	q.items = make(map[int]*Item)

	println("分裂")
}

func (q *quadTree) drawScope() {
	println(q.startX, q.startY, q.endX, q.endY)
	kit.DrawCon.AddScope(q.startX, q.startY, q.endX, q.endY)
}
