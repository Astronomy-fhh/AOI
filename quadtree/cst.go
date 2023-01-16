package quadtree

import "sync"

const MapX = 1024
const MapY = 1024
const MaxSize = 128
const MaxSplit = 8


func NewAOI() *quadTree {
	qt := &quadTree{
		items:  make(map[int]*item),
		startX: 0,
		endX:   MapX,
		startY: 0,
		endY:   MapY,
	}
	return qt
}

type quadTree struct {
	sync.RWMutex
	depth     int
	size      int
	items     map[int]*item
	startX    int
	endX      int
	startY    int
	endY      int
	leftUp    *quadTree
	leftDown  *quadTree
	rightUp   *quadTree
	rightDown *quadTree
}

type item struct {
	id int
	x  int
	y  int
}

func (q *quadTree) Enter(item *item) {
	q.Add(item)
}

func (q *quadTree) Leave(item *item) {
	q.Del(item)
}

func (q *quadTree) Get(startX, startY, endX, endY int) {
	q.Range(startX, startY, endX, endY)
}


func (q *quadTree) Add(item *item) {
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

func (q *quadTree) Del(item *item) {
	q.Lock()
	defer q.Unlock()

	q.size--
	if q.isSplit() {
		q.getOpNode(item.x, item.y).Del(item)
	} else {
		delete(q.items, item.id)
	}
}

func (q *quadTree) Range(startX, startY, endX, endY int)[]*item {
	q.Lock()
	defer q.Unlock()
	res := make([]*item, 0)
	q.RangeRecursive(startX, startY, endX, endY, res)
	return res
}

func (q *quadTree) RangeRecursive(startX, startY, endX, endY int, found []*item) {
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

func (q *quadTree) RangeItems(startX, startY, endX, endY int, found []*item){
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

func (q *quadTree) addToChild(item *item) {
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
	return q.size > MaxSize && q.depth < MaxSplit
}

func (q *quadTree) toSplit() {

	q.leftUp = &quadTree{
		items: make(map[int]*item),
		depth: q.depth + 1,
	}
	q.leftDown = &quadTree{
		items: make(map[int]*item),
		depth: q.depth + 1,
	}
	q.rightUp = &quadTree{
		items: make(map[int]*item),
		depth: q.depth + 1,
	}
	q.rightDown = &quadTree{
		items: make(map[int]*item),
		depth: q.depth + 1,
	}

	q.leftUp.startX = q.startX
	q.leftUp.endX = q.startX + (q.endX-q.startX)/2
	q.leftUp.startY = q.startY
	q.leftUp.endY = q.startY + (q.endY-q.startY)/2

	q.leftDown.startX = q.startX
	q.leftDown.endX = q.leftUp.endX
	q.leftDown.startY = q.leftUp.endY
	q.leftDown.endY = q.endY

	q.rightUp.startX = q.leftUp.endX + 1
	q.rightUp.endX = q.endX
	q.rightUp.startY = q.startY
	q.rightUp.endY = q.startY + (q.endY-q.startY)/2

	q.rightDown.startX = q.leftUp.endX + 1
	q.rightDown.endX = q.endX
	q.rightDown.startY = q.rightUp.endY + 1
	q.rightDown.endY = q.endY

	for _, item := range q.items {
		q.addToChild(item)
	}
	q.items = make(map[int]*item)
}
