package utils

import "sync"

// thread safe Queue
type WrapQueue[T any] struct {
	q  *Queue[T]
	mu *sync.Mutex
}

func NewQueue[T any]() *WrapQueue[T] {
	return &WrapQueue[T]{
		q:  New[T](),
		mu: &sync.Mutex{},
	}
}

func (q *WrapQueue[T]) Peak() T {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.q.Peek()
}

func (q *WrapQueue[T]) Len() int32 {
	q.mu.Lock()
	defer q.mu.Unlock()
	return int32(q.q.Length())
}

func (q *WrapQueue[T]) Put(data T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.q.Add(data)
}

func (q *WrapQueue[T]) PeakAndTake() T {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.q.Remove()
}
