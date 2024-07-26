package cache

import "sync"

type Queue struct {
	head *Node
	tail *Node
	mx   sync.Mutex
	cap  int
	size int
}

type Node struct {
	Val         any
	left, right *Node
}

func NewQueue(cap int) *Queue {
	if cap < 1 {
		panic("capacity must be greater than 0")
	}
	return &Queue{cap: cap}
}

func (q *Queue) MoveToTail(e *Node) (poped *Node) {
	q.mx.Lock()
	defer q.mx.Unlock()
	if q.tail == nil {
		q.append(e)
		return nil
	}
	if q.tail == e {
		return nil
	}
	if e.left != nil || e.right != nil {
		q.delete(e)
	}
	q.append(e)
	return q.shrink()
}

func (q *Queue) append(e *Node) {
	if q.tail == nil {
		q.tail, q.head = e, e
		e.left, e.right = nil, nil
		q.size = 1
		return
	}
	q.tail.right, e.left = e, q.tail
	q.tail, e.right = e, nil
	q.size++
}

func (q *Queue) delete(e *Node) {
	if q.head == e {
		q.head = e.right
	}
	if e.left != nil {
		e.left.right = e.right
	}
	if e.right != nil {
		e.right.left = e.left
	}
	e.right, e.left = nil, nil
	q.size--
}

func (q *Queue) shrink() *Node {
	if q.size > q.cap {
		return q.pop()
	}
	return nil
}

func (q *Queue) pop() (poped *Node) {
	if q.head == nil {
		return nil
	}
	poped, q.head = q.head, q.head.right
	q.head.left = nil
	q.size--
	return
}
