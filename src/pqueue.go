package jpsplus

import (
	"container/heap"
	//"fmt"
)

type PriorityQueue struct {
	pos  int
	node map[int]*Node
}

func newPriorityQueue() *PriorityQueue {
	p := new(PriorityQueue)
	p.node = make(map[int]*Node)
	return p
}

func (p *PriorityQueue) Len() int {
	return len(p.node)
}

func (p *PriorityQueue) Less(i, j int) bool {
	return p.node[i].finalCost < p.node[j].finalCost
}

func (p *PriorityQueue) Swap(i, j int) {
	p.node[i], p.node[j] = p.node[j], p.node[i]
	p.node[i].heapIndex = i
	p.node[j].heapIndex = j
}

func (p *PriorityQueue) Push(x interface{}) {
	item, ok := x.(*Node)
	if ok {
		item.heapIndex = p.pos
		p.node[p.pos] = item
		p.pos++
	}
}

func (p *PriorityQueue) Pop() interface{} {
	p.pos--
	item := p.node[p.pos]
	delete(p.node, p.pos)
	return item
}

func (p *PriorityQueue) PushNode(n *Node) {
	heap.Push(p, n)
}

func (p *PriorityQueue) PopNode() *Node {
	return heap.Pop(p).(*Node)
}

func (p *PriorityQueue) RemoveNode(n *Node) {
	heap.Remove(p, n.heapIndex)
}
