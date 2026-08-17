package linkedlist

import (
	"fmt"
	"sync"
	"time"
)

// InitDoublyLinkedList creates a one-node list holding data and returns that
// node, which is both the head and the tail.
func InitDoublyLinkedList[Data any](data Data) *DoublyLinkedList[Data] {
	ll := &DoublyLinkedList[Data]{
		next:        nil,
		initialised: new(bool),
		data:        data,
		size:        new(int),
		listMu:      &sync.Mutex{},
	}
	*ll.size = 1
	*ll.initialised = true
	// If two pointers point to the same address A1 which in turn points to address A2, then changing the
	// value at A2 changes it in both the variables. To prevent this, we need to have different addresses
	// A1.h and A1.t pointing to A2. Assigning current node as head and tail but via different variables does
	// exactly that, since h and t have different addresses. Had the code below been `ll.head, ll.tail = &ll, &ll`
	// instead, this would've been the case of ll.head -> A1 and ll.tail -> A1. Now, ll.head -> A1.h and
	// ll.tail -> A1.t. This is because each and every variable has a different address in memory.
	h, t := ll, ll
	ll.head, ll.tail = &h, &t
	return ll
}

// Add inserts data immediately after this node and returns the new node,
// advancing the tail if this node was the last one.
func (ll *DoublyLinkedList[Data]) Add(data Data) *DoublyLinkedList[Data] {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("DoublyLinkedList is not initialised. Please use InitDoublyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("DoublyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	newNode := &DoublyLinkedList[Data]{
		next:        ll.next,
		initialised: ll.initialised,
		size:        ll.size,
		data:        data,
		listMu:      ll.listMu,
		head:        ll.head,
		tail:        ll.tail,
		prev:        ll,
	}
	ll.next = newNode
	*ll.size += 1
	if ll == *ll.tail {
		*ll.tail = newNode
	} else { // If not tail, then the `prev` of the next node should also point this node.
		newNode.next.prev = newNode
	}

	return newNode
}

// Append adds data as a new tail node and returns it. It can be called from any
// node in the list and is atomic with respect to concurrent appends.
func (ll *DoublyLinkedList[Data]) Append(data Data) *DoublyLinkedList[Data] {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("DoublyLinkedList is not initialised. Please use InitDoublyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("DoublyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	newNode := &DoublyLinkedList[Data]{
		next:        nil,
		initialised: ll.initialised,
		data:        data,
		size:        ll.size,
		listMu:      ll.listMu,
		head:        ll.head,
		tail:        ll.tail,
		prev:        *ll.tail,
	}
	(*ll.tail).next, *ll.tail = newNode, newNode
	*ll.size += 1

	return newNode
}

// Next returns the node after this one, or nil if this node is the tail.
func (ll *DoublyLinkedList[Data]) Next() *DoublyLinkedList[Data] {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("DoublyLinkedList is not initialised. Please use InitDoublyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("DoublyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	return ll.next
}

// Prev returns the node before this one, or nil if this node is the head.
func (ll *DoublyLinkedList[Data]) Prev() *DoublyLinkedList[Data] {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("DoublyLinkedList is not initialised. Please use InitDoublyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("DoublyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	return ll.prev
}

// Delete unlinks this node from the list and reports whether it was removed; it
// returns false when this is the only node left. A deleted node is permanently
// unusable afterwards.
func (ll *DoublyLinkedList[Data]) Delete() bool {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("DoublyLinkedList is not initialised. Please use InitDoublyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("DoublyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	// Capture the next node to be deleted
	nextNode, prevNode := ll.next, ll.prev
	switch {
	case prevNode == nil && nextNode == nil: // Do nothing
		return false
	case ll == *ll.tail: // If tail, update tail pointer everywhere to the previous node
		*ll.tail = prevNode
		prevNode.next = nil
	case ll == *ll.head: // If head, update head pointer everywhere to the next node
		*ll.head = nextNode
		nextNode.prev = nil
	// If somewhere in the middle, update the next and prev pointers of adjacent nodes
	default:
		prevNode.next, nextNode.prev = nextNode, prevNode
	}

	// Decrementing the size of the list before severing pointers
	*ll.size -= 1
	// Update the pointers of the deleted node to nil
	ll.prev, ll.next, ll.head, ll.size, ll.tail = nil, nil, nil, nil, nil
	// Mark the next node as deleted
	ll.isDeleted = true

	return true
}

// Data returns the value held by this node.
func (ll *DoublyLinkedList[Data]) Data() Data {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("DoublyLinkedList is not initialised. Please use InitDoublyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("DoublyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	return ll.data
}

// Head returns the first node of the list.
func (ll *DoublyLinkedList[Data]) Head() *DoublyLinkedList[Data] {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("DoublyLinkedList is not initialised. Please use InitDoublyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("DoublyLinkedList node has been deleted. Please do not use this node anymore.")
	}
	return *ll.head
}

// Size returns the number of nodes currently in the list.
func (ll *DoublyLinkedList[Data]) Size() int {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("DoublyLinkedList is not initialised. Please use InitDoublyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("DoublyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	return *ll.size
}

// Tail returns the last node of the list.
func (ll *DoublyLinkedList[Data]) Tail() *DoublyLinkedList[Data] {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("DoublyLinkedList is not initialised. Please use InitDoublyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("DoublyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	return *ll.tail
}

// PrintList writes every node from head to tail to stdout, followed by how long
// building and printing the output took.
func (ll *DoublyLinkedList[Data]) PrintList() {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("DoublyLinkedList is not initialised. Please use InitDoublyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("DoublyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	// TODO: See if bufio and string builders can be used to make this faster.
	start := time.Now()
	current, idx := *ll.head, 0
	outputList := make([]string, *ll.size) // Allocate size upfront to avoid resizing and slowing down the printing process
	for current != nil {
		outputList[idx] = fmt.Sprintf("Node %d: %v\n", idx, current.data)
		// ^Won't be concatenating to one single string, that makes things slow. Hence the outputList array.
		current = current.next
		idx++
	}

	strtPrint := time.Now()
	fmt.Println(outputList)

	fmt.Println("Time taken to form the output:", time.Since(start).Seconds())
	fmt.Println("Time taken to print the list:", time.Since(strtPrint).Seconds())
}

// Print writes this node's value and the list's current size to stdout.
func (ll *DoublyLinkedList[Data]) Print() {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("DoublyLinkedList is not initialised. Please use InitDoublyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("DoublyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	fmt.Println("Data:", ll.data, "Size:", *ll.size)
}
