package stack

import (
	"sync"
)

// StackNode holds one stack element. A node returned by Pop is detached from the
// stack; one returned by Peek or Lookup stays live only while it remains in it.
type StackNode[Data any] struct {
	data Data
	next *StackNode[Data]
}

// Data returns the value held by this node.
func (n *StackNode[Data]) Data() Data {
	return n.data
}

// LifoStack is a last-in-first-out stack. It must be created with InitLifoStack;
// every method panics on a zero value. All methods are safe for concurrent use.
type LifoStack[Data any] struct {
	mu          *sync.Mutex
	top         *StackNode[Data]
	size        *int
	initialised bool
}

// InitLifoStack creates a stack holding data as its only element.
func InitLifoStack[Data any](data Data) *LifoStack[Data] {
	ls := &LifoStack[Data]{
		top:         &StackNode[Data]{data: data, next: nil},
		size:        new(int),
		mu:          &sync.Mutex{},
		initialised: true,
	}
	*ls.size = 1

	return ls
}

// Push puts data on top of the stack and returns the new size.
func (ls *LifoStack[Data]) Push(data Data) int {
	ls.MustBeInitialised()
	ls.mu.Lock()
	defer ls.mu.Unlock()

	nextNode := ls.top
	ls.top = &StackNode[Data]{data: data, next: nextNode}
	*ls.size += 1

	return *ls.size
}

// Pop removes the top element and returns it, or nil if the stack is empty.
// The returned node is detached and cannot reach the rest of the stack.
func (ls *LifoStack[Data]) Pop() *StackNode[Data] {
	ls.MustBeInitialised()
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if ls.size == nil || *ls.size == 0 {
		return nil
	}
	poppedNode := ls.top
	ls.top = ls.top.next
	poppedNode.next = nil
	*ls.size -= 1

	return poppedNode
}

// Peek returns the top element without removing it, or nil if the stack is empty.
func (ls *LifoStack[Data]) Peek() *StackNode[Data] {
	ls.MustBeInitialised()
	ls.mu.Lock()
	defer ls.mu.Unlock()

	return ls.top
}

// Lookup walks the stack from top to bottom, calling fn for each node, and stops
// at the first node for which fn returns true. fn must not call any method on
// this stack: Lookup holds the stack's lock while fn runs, and the lock is not
// reentrant, so doing so deadlocks. Reading a node's Data is safe.
func (ls *LifoStack[Data]) Lookup(fn func(*StackNode[Data]) bool) {
	ls.MustBeInitialised()
	ls.mu.Lock()
	defer ls.mu.Unlock()

	for node := ls.top; node != nil; node = node.next {
		if fn(node) {
			return
		}
	}
}

// Size returns the number of elements currently on the stack.
func (ls *LifoStack[Data]) Size() int {
	ls.MustBeInitialised()
	ls.mu.Lock()
	defer ls.mu.Unlock()

	return *ls.size
}

// MustBeInitialised panics unless this stack was built by InitLifoStack.
func (ls *LifoStack[Data]) MustBeInitialised() {
	if ls == nil || !ls.initialised {
		panic("LifoStack is not initialised. Please use InitLifoStack to create a new stack.")
	}
}
