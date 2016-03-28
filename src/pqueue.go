package jpsplus

import (
	"container/heap"
	//"fmt"
)

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].f < pq[j].f
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].heap_index = i
	pq[j].heap_index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	a := *pq
	n := len(a)
	a = a[0 : n+1]
	item := x.(*Node)
	item.heap_index = n
	a[n] = item
	*pq = a
}

func (pq *PriorityQueue) Pop() interface{} {
	a := *pq
	n := len(a)
	item := a[n-1]
	item.heap_index = -1
	*pq = a[0 : n-1]
	return item
}

func (pq *PriorityQueue) PushNode(n *Node) {
	heap.Push(pq, n)
}

func (pq *PriorityQueue) PopNode() *Node {
	return heap.Pop(pq).(*Node)
}

func (pq *PriorityQueue) RemoveNode(n *Node) {
	heap.Remove(pq, n.heap_index)
}
