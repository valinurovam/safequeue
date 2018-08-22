package safequeue

import (
	"sync"
)

// SafeQueue implements thread safe FIFO queue
type SafeQueue struct {
	sync.Mutex
	shards    [][]interface{}
	shardSize int
	tailIdx   int
	tail      []interface{}
	tailPos   int
	headIdx   int
	head      []interface{}
	headPos   int
	length    uint64
}

// NewSafeQueue returns new instance of SafeQueue
func NewSafeQueue(shardSize int) *SafeQueue {
	queue := &SafeQueue{
		shardSize: shardSize,
		shards:    [][]interface{}{make([]interface{}, shardSize)},
	}

	queue.tailIdx = 0
	queue.tail = queue.shards[queue.tailIdx]
	queue.headIdx = 0
	queue.head = queue.shards[queue.headIdx]
	return queue
}

// Push append item into queue tail
func (queue *SafeQueue) Push(item interface{}) {
	queue.Lock()
	defer queue.Unlock()

	queue.tail[queue.tailPos] = item
	queue.tailPos++
	queue.length++

	if queue.tailPos == queue.shardSize {
		queue.tailPos = 0
		queue.tailIdx = len(queue.shards)

		buffer := make([][]interface{}, len(queue.shards)+1)
		buffer[queue.tailIdx] = make([]interface{}, queue.shardSize)
		copy(buffer, queue.shards)

		queue.shards = buffer
		queue.tail = queue.shards[queue.tailIdx]
	}
}

// PushHead append item into queue head
func (queue *SafeQueue) PushHead(item interface{}) {
	queue.Lock()
	defer queue.Unlock()

	if queue.headPos == 0 {
		buffer := make([][]interface{}, len(queue.shards)+1)
		copy(buffer[1:], queue.shards)
		buffer[queue.headIdx] = make([]interface{}, queue.shardSize)

		queue.shards = buffer
		queue.tailIdx++
		queue.headPos = queue.shardSize
		queue.tail = queue.shards[queue.tailIdx]
		queue.head = queue.shards[queue.headIdx]
	}
	queue.length++
	queue.headPos--
	queue.head[queue.headPos] = item
}

// Pop returns item from queue head
func (queue *SafeQueue) Pop() (item interface{}) {
	queue.Lock()
	item = queue.DirtyPop()
	queue.Unlock()
	return
}

// DirtyPop returns item from queue head
// DirtyPop is not thread-safe
func (queue *SafeQueue) DirtyPop() (item interface{}) {
	item, queue.head[queue.headPos] = queue.head[queue.headPos], nil
	if item == nil {
		return item
	}
	queue.headPos++
	queue.length--
	if queue.headPos == queue.shardSize {
		buffer := make([][]interface{}, len(queue.shards)-1)
		copy(buffer, queue.shards[queue.headIdx+1:])

		queue.shards = buffer

		queue.headPos = 0
		queue.tailIdx--
		queue.head = queue.shards[queue.headIdx]
	}
	return
}

// Length returns queue length
func (queue *SafeQueue) Length() uint64 {
	queue.Lock()
	defer queue.Unlock()
	return queue.length
}

// DirtyLength returns queue length
// DirtyLength is not thread-safe
func (queue *SafeQueue) DirtyLength() uint64 {
	return queue.length
}

// HeadItem returns queue head item
// HeadItem is not thread-safe
func (queue *SafeQueue) HeadItem() (res interface{}) {
	return queue.head[queue.headPos]
}

// DirtyPurge clean queue
// DirtyPurge is not thread-safe
func (queue *SafeQueue) DirtyPurge() {
	queue.shards = [][]interface{}{make([]interface{}, queue.shardSize)}
	queue.tailIdx = 0
	queue.tail = queue.shards[queue.tailIdx]
	queue.headIdx = 0
	queue.head = queue.shards[queue.headIdx]
	queue.length = 0
}

// Purge clean queue
func (queue *SafeQueue) Purge() {
	queue.Lock()
	defer queue.Unlock()
	queue.DirtyPurge()
}
