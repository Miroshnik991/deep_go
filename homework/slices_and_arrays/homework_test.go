package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularQueue struct {
	values []int
	head   int
	tail   int
}

func NewCircularQueue(size int) CircularQueue {
	return CircularQueue{
		values: make([]int, size),
		head:   -1,
		tail:   -1,
	}
}

func (q *CircularQueue) Push(value int) bool {
	if q.Full() {
		return false
	} else if q.Empty() {
		q.values[0] = value
		q.head, q.tail = 0, 0
		return true
	}
	if q.reachedTheEnd(q.tail) {
		q.tail = 0
	} else {
		q.tail++
	}
	q.values[q.tail] = value
	return true
}

func (q *CircularQueue) Pop() bool {
	if q.Empty() {
		return false
	}
	if q.lastElement() {
		q.resetPointers()
	}
	if q.reachedTheEnd(q.head) {
		q.head = 0
		return true
	}
	q.head++
	return true
}

func (q *CircularQueue) Front() int {
	if q.Empty() {
		return -1
	}
	return q.values[q.head]
}

func (q *CircularQueue) Back() int {
	if q.Empty() {
		return -1
	}
	return q.values[q.tail]
}

func (q *CircularQueue) Empty() bool {
	return q.head == -1 && q.tail == -1
}

func (q *CircularQueue) Full() bool {
	if (q.head == 0 && q.tail == len(q.values)-1) || (q.head == q.tail+1) {
		return true
	}
	return false
}

func (q *CircularQueue) lastElement() bool {
	if q.head == q.tail {
		return true
	}
	return false
}

func (q *CircularQueue) resetPointers() {
	q.head, q.tail = -1, -1
}

func (q *CircularQueue) reachedTheEnd(index int) bool {
	if index+1 == len(q.values) {
		return true
	}
	return false
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
