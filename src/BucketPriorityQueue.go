package jpsplus

type BucketPriorityQueue struct {
	m_numBuckets        int
	m_lowestNonEmptyBin int
	m_numNodesTracked   int
	m_division          uint
	m_baseCost          uint
	m_bin               []*UnsortedPriorityQueue
	m_maxFreeBuckets    int
	m_nextFreeBucket    int
	m_freeBuckets       []*UnsortedPriorityQueue
}

func newBucketPriorityQueue(buckets int, arraySize int, division uint) *BucketPriorityQueue {
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

func (b BucketPriorityQueue) GetBinIndex(cost uint) int {
	return ((cost - b.m_baseCost) / b.m_division)
}

func (b *BucketPriorityQueue) Pop() *DijkstraPathfindingNode {
	node := b.m_bin[b.m_lowestNonEmptyBin].Pop()
	b.m_numNodesTracked -= 1

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
	index := GetBinIndex(node.m_givenCost)
	if 0 == b.m_bin[index] {
		b.m_nextFreeBucket += 1
		b.m_bin[index] = b.m_freeBuckets[b.m_nextFreeBucket]
	}
	b.m_bin[index].Push(node)
	if index < b.m_lowestNonEmptyBin {
		b.m_lowestNonEmptyBin = index
	}
}

func (b *BucketPriorityQueue) DecreaseKey(node *DijkstraPathfindingNode, lastCost uint) {
	// Remove node
	index := GetBinIndex(lastCost)
	b.m_bin[index].Remove(node)

	if b.m_bin[index].Empty(node.m_iteration) {
		b.m_nextFreeBucket -= 1
		b.m_freeBuckets[b.m_nextFreeBucket] = b.m_bin[index]
		b.m_bin[index] = 0
	}

	// Push node
	index = GetBinIndex(node.m_givenCost)

	if b.m_bin[index] == 0 {
		b.m_nextFreeBucket += 1
		b.m_bin[index] = b.m_freeBuckets[b.m_nextFreeBucket]
	}

	b.m_bin[index].Push(node)

	if index < b.m_lowestNonEmptyBin {
		b.m_lowestNonEmptyBin = index
	}
}

func (b BucketPriorityQueue) Empty() {
	return 0 == b.m_numNodesTracked
}

func (b *BucketPriorityQueue) Reset() {
	b.m_lowestNonEmptyBin = m_numBuckets
	b.m_numNodesTracked = 0
	b.m_baseCost = 0
}

func (b *BucketPriorityQueue) SetBaseCost(baseCost uint) {
	b.m_baseCost = baseCost
}
