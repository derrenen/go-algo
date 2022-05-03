package circqueue

import "sync"

//Error definition
//+--------------------------------------------------------------------------------------------------------------------+

type NoFreeCapErr struct {
	err string
}

func NewNoFreeCapErr() error {
	return &NoFreeCapErr{err: "Queue is full!"}
}

func (e *NoFreeCapErr) Error() string {
	return e.err
}

type QueueEmptyErr struct {
	err string
}

func NewQueueEmptyErr() error {
	return &QueueEmptyErr{err: "Queue is empty!"}
}

func (e *QueueEmptyErr) Error() string {
	return e.err
}

//The thread safe implimentation of circular queue (https://en.wikipedia.org/wiki/Circular_buffer) data structure
//+--------------------------------------------------------------------------------------------------------------------+

type CircularQueue struct {
	mu     *sync.Mutex
	data   []int
	cap    int
	curCap int
	head   int
}

func NewCircularQueue(queueCap int) *CircularQueue {
	return &CircularQueue{
		mu:     &sync.Mutex{},
		data:   make([]int, queueCap),
		cap:    queueCap,
		curCap: 0,
		head:   0,
	}
}

func (cq *CircularQueue) EnQueue(val int) error {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if cq.curCap != cq.cap {
		cq.data[cq.head+cq.curCap] = val
		cq.curCap += 1
		return nil
	}
	return NewNoFreeCapErr()
}

func (cq *CircularQueue) DeQueue() error {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if cq.curCap != 0 {
		cq.head = cq.head + 1
		cq.curCap -= 1
		return nil
	}
	return NewQueueEmptyErr()
}

func (cq *CircularQueue) Front() (int, error) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if cq.curCap != 0 {
		return cq.data[cq.head], nil
	}
	return 0, NewQueueEmptyErr()
}

func (cq *CircularQueue) Rear() (int, error) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if cq.curCap != 0 {
		return cq.data[cq.head+cq.curCap-1], nil
	}
	return 0, NewQueueEmptyErr()
}

func (cq *CircularQueue) IsEmpty() bool {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	return cq.curCap == 0
}

func (cq *CircularQueue) IsFull() bool {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	return cq.curCap == cq.cap
}
