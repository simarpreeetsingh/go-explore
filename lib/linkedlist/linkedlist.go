// Package linkedlist provides generic, concurrency-safe linked lists in which
// every node is also a handle to the list it belongs to, so list-wide operations
// such as Size, Head, and Tail can be called from any node.
//
// A list must be created with InitSinglyLinkedList or InitDoublyLinkedList; the
// zero value is unusable. Every method panics if the receiver was not built by a
// constructor, or has since been removed from its list.
package linkedlist

import (
	"sync"
)

// Write now, these linked lists do not share locks amongst chained calls. This will lead to
// panics if deletion is performed on nodes who's data were fetched in a seperate goroutine.
// The outcome of one call always depicts the truth, but once that truth is expected to be
// shared amongst chained calls of functions, TOCTOU might come into play

// SinglyLinkedList is a node in a forward-linked list. Nodes share the list's
// tail, size, and mutex, and can only be traversed from head towards tail.
type SinglyLinkedList[Data any] struct {
	next        *SinglyLinkedList[Data]
	head        *SinglyLinkedList[Data]
	tail        **SinglyLinkedList[Data]
	size        *int
	initialised *bool
	listMu      *sync.Mutex
	data        Data
	isDeleted   bool
}

// DoublyLinkedList is a node in a bidirectionally-linked list. Nodes share the
// list's head, tail, size, and mutex, and can be traversed in either direction.
type DoublyLinkedList[Data any] struct {
	next        *DoublyLinkedList[Data]
	prev        *DoublyLinkedList[Data]
	head        **DoublyLinkedList[Data]
	tail        **DoublyLinkedList[Data]
	size        *int
	initialised *bool
	listMu      *sync.Mutex
	data        Data
	isDeleted   bool
}
