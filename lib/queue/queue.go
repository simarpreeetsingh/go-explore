// Package queue provides a generic, concurrency-safe FIFO queue.
//
// A queue must be created with InitQueue; every method panics on a zero value.
// Individual methods are safe for concurrent use, but sequences of them are not:
// a node returned by one call may be dequeued before the next call acts on it.
package queue

import "sync"

// QueueNode holds one queued element. A node returned by Dequeue is detached
// from the queue; one returned by Peek, Tail, or Lookup stays live only while it
// remains in it.
type QueueNode[Data any] struct {
	data Data
	next *QueueNode[Data]
}

// Data returns the value held by this node.
func (n *QueueNode[Data]) Data() Data {
	return n.data
}

// Queue is a first-in-first-out queue. Elements are added at the tail by Enqueue
// and removed from the head by Dequeue.
type Queue[Data any] struct {
	size        int
	mu          sync.Mutex
	head        *QueueNode[Data]
	tail        *QueueNode[Data]
	initialised bool
}

// InitQueue creates a queue holding data as its only element.
func InitQueue[Data any](data Data) *Queue[Data] {
	q := &Queue[Data]{
		size: 1, head: &QueueNode[Data]{data: data}, initialised: true,
	}
	q.tail = q.head
	return q
}

// Enqueue adds data at the tail of the queue and returns the new size.
func (q *Queue[Data]) Enqueue(data Data) int { // Return bool if this function stops panicking (thinking: enqueuing needs confirmation, size validations are usually put before enqueue)
	q.mustBeInitialised()
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.head == nil {
		q.head = &QueueNode[Data]{data: data}
		q.tail = q.head
		q.size++
		return q.size
	}
	q.tail.next = &QueueNode[Data]{data: data}
	q.tail = q.tail.next
	q.size++

	return q.size
}

// Dequeue removes the element at the head and returns it, or nil if the queue is
// empty. The returned node is detached and cannot reach the rest of the queue.
func (q *Queue[Data]) Dequeue() *QueueNode[Data] {
	q.mustBeInitialised()
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isEmpty() {
		return nil
	}

	dequeuedNode := q.head
	q.head = q.head.next
	dequeuedNode.next = nil
	q.size--

	if q.isEmpty() {
		q.tail = nil
	}

	return dequeuedNode
}

// Peek returns the head element without removing it, or nil if the queue is empty.
func (q *Queue[Data]) Peek() *QueueNode[Data] {
	q.mustBeInitialised()
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.head
}

// Tail returns the element most recently enqueued, or nil if the queue is empty.
func (q *Queue[Data]) Tail() *QueueNode[Data] {
	q.mustBeInitialised()
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.tail
}

// Lookup walks the queue from head to tail, calling fn for each node, and returns
// the first node for which fn returns true, or nil if none does. fn must not call
// any method on this queue: Lookup holds the queue's lock while fn runs, and the
// lock is not reentrant, so doing so deadlocks. Reading a node's Data is safe.
func (q *Queue[Data]) Lookup(fn func(*QueueNode[Data]) bool) *QueueNode[Data] {
	q.mustBeInitialised()
	q.mu.Lock()
	defer q.mu.Unlock()

	for n := q.head; n != nil; n = n.next {
		if fn(n) {
			return n
		}
	}

	return nil
}

// Size returns the number of elements currently in the queue.
func (q *Queue[Data]) Size() int {
	q.mustBeInitialised()
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.size
}

// mustBeInitialised panics unless this queue was built by InitQueue.
func (q *Queue[Data]) mustBeInitialised() {
	if q == nil || !q.initialised {
		panic("Queue is not initialised. Please use InitQueue to create a new queue.")
	}
}

// isEmpty reports whether the queue holds no elements. Callers must hold q.mu.
func (q *Queue[Data]) isEmpty() bool {
	return q == nil || q.size <= 0 // Should never be less than 0 though
}
