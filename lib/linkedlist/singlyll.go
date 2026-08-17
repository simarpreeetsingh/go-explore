package linkedlist

import (
	"fmt"
	"sync"
	"time"
)

// InitSinglyLinkedList creates a one-node list holding data and returns that
// node, which is both the head and the tail.
func InitSinglyLinkedList[Data any](data Data) *SinglyLinkedList[Data] {
	ll := &SinglyLinkedList[Data]{
		next:        nil,
		initialised: new(bool),
		data:        data,
		size:        new(int),
		listMu:      &sync.Mutex{},
	}
	*ll.size = 1
	*ll.initialised = true
	ll.head, ll.tail = ll, &ll
	return ll
}

// Append adds data as a new tail node and returns it. It can be called from any
// node in the list and is atomic with respect to concurrent appends.
func (ll *SinglyLinkedList[Data]) Append(data Data) *SinglyLinkedList[Data] {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("SinglyLinkedList is not initialised. Please use InitSinglyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("SinglyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	newNode := &SinglyLinkedList[Data]{
		next:        nil,
		initialised: ll.initialised,
		data:        data,
		size:        ll.size,
		listMu:      ll.listMu,
		head:        ll.head,
		tail:        ll.tail,
	}
	(*ll.tail).next, *ll.tail = newNode, newNode
	*ll.size += 1

	return newNode
}

// Next returns the node after this one, or nil if this node is the tail.
func (ll *SinglyLinkedList[Data]) Next() *SinglyLinkedList[Data] {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("SinglyLinkedList is not initialised. Please use InitSinglyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("SinglyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	return ll.next
}

// DeleteNext removes the node immediately after this one and reports whether a
// node was removed; it returns false when this node is the tail. The removed node
// is permanently unusable afterwards.
func (ll *SinglyLinkedList[Data]) DeleteNext() bool {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("SinglyLinkedList is not initialised. Please use InitSinglyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("SinglyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	// Tail node, nothing present next
	if ll == *ll.tail {
		return false
	}

	if ll.next == *ll.tail {
		*ll.tail = ll
	}
	// Capture the next node to be deleted
	deletedNode := ll.next
	// Update the next pointer of the current node to skip the deleted node
	ll.next = ll.next.next
	// Update the pointers of the deleted node to nil
	deletedNode.next, deletedNode.head, deletedNode.size, deletedNode.tail = nil, nil, nil, nil
	// Mark the next node as deleted
	deletedNode.isDeleted = true

	// Decrement the size of the list
	*ll.size -= 1

	return true
}

// Data returns the value held by this node.
func (ll *SinglyLinkedList[Data]) Data() Data {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("SinglyLinkedList is not initialised. Please use InitSinglyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("SinglyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	return ll.data
}

// Head returns the first node of the list.
func (ll *SinglyLinkedList[Data]) Head() *SinglyLinkedList[Data] {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("SinglyLinkedList is not initialised. Please use InitSinglyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("SinglyLinkedList node has been deleted. Please do not use this node anymore.")
	}
	return ll.head
}

// Size returns the number of nodes currently in the list.
func (ll *SinglyLinkedList[Data]) Size() int {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("SinglyLinkedList is not initialised. Please use InitSinglyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("SinglyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	return *ll.size
}

// Tail returns the last node of the list.
func (ll *SinglyLinkedList[Data]) Tail() *SinglyLinkedList[Data] {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("SinglyLinkedList is not initialised. Please use InitSinglyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("SinglyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	return *ll.tail
}

// PrintList writes every node from head to tail to stdout, followed by how long
// building and printing the output took.
func (ll *SinglyLinkedList[Data]) PrintList() {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("SinglyLinkedList is not initialised. Please use InitSinglyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("SinglyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	// TODO: See if bufio and string builders can be used to make this faster.
	start := time.Now()
	current, idx := ll.head, 0
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
func (ll *SinglyLinkedList[Data]) Print() {
	if ll == nil || ll.initialised == nil || !*ll.initialised {
		panic("SinglyLinkedList is not initialised. Please use InitSinglyLinkedList to create a new list.")
	}
	ll.listMu.Lock()
	defer ll.listMu.Unlock()

	if ll.isDeleted {
		panic("SinglyLinkedList node has been deleted. Please do not use this node anymore.")
	}

	fmt.Println("Data:", ll.data, "Size:", *ll.size)
}
