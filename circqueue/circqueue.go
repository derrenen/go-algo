package circqueue

import (
	"sync"
)

//Error definition
//+--------------------------------------------------------------------------------------------------------------------+

type NoFreeCapErr struct {
	err string
}

func newNoFreeCapErr() error {
	return &NoFreeCapErr{err: "Queue is full!"}
}

func (e *NoFreeCapErr) Error() string {
	return e.err
}

type QueueEmptyErr struct {
	err string
}

func newQueueEmptyErr() error {
	return &QueueEmptyErr{err: "Queue is empty!"}
}

func (e *QueueEmptyErr) Error() string {
	return e.err
}

//The thread safe implementation of circular queue (https://en.wikipedia.org/wiki/Circular_buffer) data structure
//+--------------------------------------------------------------------------------------------------------------------+

type circularQueue struct {
	mu     *sync.Mutex
	data   []int
	cap    int
	curCap int
	head   int
}

func NewCircularQueue(queueCap int) *circularQueue {
	return &circularQueue{
		mu:     &sync.Mutex{},
		data:   make([]int, queueCap),
		cap:    queueCap,
		curCap: 0,
		head:   0,
	}
}

func (cq *circularQueue) EnQueue(val int) error {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if cq.curCap != cq.cap {
		cq.data[(cq.head+cq.curCap)%cq.cap] = val
		cq.curCap += 1
		return nil
	}
	return newNoFreeCapErr()
}

func (cq *circularQueue) DeQueue() error {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if cq.curCap != 0 {
		cq.head = (cq.head + 1) % cq.cap
		cq.curCap -= 1
		return nil
	}
	return newQueueEmptyErr()
}

func (cq *circularQueue) Front() (int, error) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if cq.curCap != 0 {
		return cq.data[cq.head], nil
	}
	return 0, newQueueEmptyErr()
}

func (cq *circularQueue) Rear() (int, error) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if cq.curCap != 0 {
		return cq.data[(cq.head+cq.curCap-1)%cq.cap], nil
	}
	return 0, newQueueEmptyErr()
}

func (cq *circularQueue) IsEmpty() bool {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	return cq.curCap == 0
}

func (cq *circularQueue) IsFull() bool {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	return cq.curCap == cq.cap
}
