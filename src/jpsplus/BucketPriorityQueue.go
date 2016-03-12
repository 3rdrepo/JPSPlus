package jpsplus

import (
// "fmt"
)

type BucketPriorityQueue struct {
	m_numBuckets        int
	m_lowestNonEmptyBin int
	m_numNodesTracked   int
	m_division          int
	m_baseCost          int64
	m_bin               []*UnsortedPriorityQueue
	m_maxFreeBuckets    int
	m_nextFreeBucket    int
	m_freeBuckets       []*UnsortedPriorityQueue
}

func newBucketPriorityQueue(buckets int, arraySize int, division int) *BucketPriorityQueue {
	b := new(BucketPriorityQueue)
	b.m_numBuckets = buckets
	b.m_division = division

	b.Reset()

	// Allocate a bunch of free buckets
	b.m_maxFreeBuckets = 200
	b.m_nextFreeBucket = 0
	b.m_freeBuckets = make([]*UnsortedPriorityQueue, b.m_maxFreeBuckets)
	for m := 0; m < b.m_maxFreeBuckets; m++ {
		b.m_freeBuckets[m] = newUnsortedPriorityQueue(arraySize)
	}

	// Allocate bucket slots
	b.m_bin = make([]*UnsortedPriorityQueue, b.m_numBuckets)
	// for m := 0; m < b.m_numBuckets; m++ {
	// 	b.m_bin[m] = 0
	// }
	return b
}

func (b BucketPriorityQueue) GetBinIndex(cost int64) int {
	return int((cost - b.m_baseCost) / int64(b.m_division))
}

func (b *BucketPriorityQueue) Pop() *DijkstraPathfindingNode {
	node := b.m_bin[b.m_lowestNonEmptyBin].Pop()
	b.m_numNodesTracked -= 1
	// fmt.Printf("node %#v\n", node)
	// fmt.Printf("b.m_lowestNonEmptyBin %v\n", b.m_lowestNonEmptyBin)
	// fmt.Printf("b.m_bin[b.m_lowestNonEmptyBin] %#v\n", b.m_bin[b.m_lowestNonEmptyBin])

	if b.m_bin[b.m_lowestNonEmptyBin].Empty(node.m_iteration) {
		b.m_nextFreeBucket -= 1
		b.m_freeBuckets[b.m_nextFreeBucket] = b.m_bin[b.m_lowestNonEmptyBin]
		b.m_bin[b.m_lowestNonEmptyBin] = nil
	}

	if b.m_numNodesTracked > 0 {
		// Find the next non-empty bin
		for b.m_lowestNonEmptyBin < b.m_numBuckets {
			if b.m_bin[b.m_lowestNonEmptyBin] != nil &&
				!b.m_bin[b.m_lowestNonEmptyBin].Empty(node.m_iteration) {
				break
			}
			b.m_lowestNonEmptyBin += 1
		}
	} else {
		b.m_lowestNonEmptyBin = b.m_numBuckets
	}
	return node
}

func (b *BucketPriorityQueue) Push(node *DijkstraPathfindingNode) {
	b.m_numNodesTracked += 1
	index := b.GetBinIndex(node.m_givenCost)
	if nil == b.m_bin[index] {
		b.m_bin[index] = b.m_freeBuckets[b.m_nextFreeBucket]
		b.m_nextFreeBucket += 1
	}
	b.m_bin[index].Push(node)
	if index < b.m_lowestNonEmptyBin {
		b.m_lowestNonEmptyBin = index
	}
}

func (b *BucketPriorityQueue) DecreaseKey(node *DijkstraPathfindingNode, lastCost int64) {
	// Remove node
	index := b.GetBinIndex(lastCost)
	b.m_bin[index].Remove(node)

	if b.m_bin[index].Empty(node.m_iteration) {
		b.m_nextFreeBucket -= 1
		b.m_freeBuckets[b.m_nextFreeBucket] = b.m_bin[index]
		b.m_bin[index] = nil
	}

	// Push node
	index = b.GetBinIndex(node.m_givenCost)

	if b.m_bin[index] == nil {
		b.m_bin[index] = b.m_freeBuckets[b.m_nextFreeBucket]
		b.m_nextFreeBucket += 1
	}

	b.m_bin[index].Push(node)

	if index < b.m_lowestNonEmptyBin {
		b.m_lowestNonEmptyBin = index
	}
}

func (b BucketPriorityQueue) Empty() bool {
	return 0 == b.m_numNodesTracked
}

func (b *BucketPriorityQueue) Reset() {
	b.m_lowestNonEmptyBin = b.m_numBuckets
	b.m_numNodesTracked = 0
	b.m_baseCost = 0
}

func (b *BucketPriorityQueue) SetBaseCost(baseCost int64) {
	b.m_baseCost = baseCost
}
