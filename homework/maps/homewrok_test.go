package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	key        int
	value      int
	leftChild  *Node
	rightChild *Node
}

type OrderedMap struct {
	rootNode *Node
	size     int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.rootNode == nil {
		m.rootNode = &Node{
			key:   key,
			value: value,
		}
		m.size++
		return
	}
	if search(m.rootNode, key) {
		panic("Key already exists")
	}
	m.rootNode = insert(m.rootNode, key, value)
	m.size++
}

func (m *OrderedMap) Erase(key int) {
	if search(m.rootNode, key) {
		m.rootNode = deleteNode(m.rootNode, key)
		m.size--
	}
}

func (m *OrderedMap) Contains(key int) bool {
	return search(m.rootNode, key)
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	traverseInOrderWithAction(m.rootNode, action)
}

func traverseInOrderWithAction(root *Node, action func(int, int)) {
	if root == nil {
		return
	}
	traverseInOrderWithAction(root.leftChild, action)
	action(root.key, root.value)
	traverseInOrderWithAction(root.rightChild, action)
}

func search(root *Node, target int) bool {
	if root == nil {
		return false
	}
	if target > root.key {
		return search(root.rightChild, target)
	} else if target < root.key {
		return search(root.leftChild, target)
	}
	return true
}

func minValueNode(subtreeRoot *Node) *Node {
	curr := subtreeRoot
	for curr != nil && curr.leftChild != nil {
		curr = curr.leftChild
	}
	return curr
}

func insert(root *Node, targetKey int, targetVal int) *Node {
	if targetKey > root.key && root.rightChild == nil {
		root.rightChild = &Node{
			key:   targetKey,
			value: targetVal,
		}
	} else if targetKey < root.key && root.leftChild == nil {
		root.leftChild = &Node{
			key:   targetKey,
			value: targetVal,
		}
	} else {
		if targetKey > root.key {
			root.rightChild = insert(root.rightChild, targetKey, targetVal)
		} else if targetKey < root.key {
			root.leftChild = insert(root.leftChild, targetKey, targetVal)
		}
	}
	return root
}

func deleteNode(root *Node, targetKey int) *Node {
	if root == nil {
		return nil
	}

	if targetKey > root.key {
		root.rightChild = deleteNode(root.rightChild, targetKey)
	} else if targetKey < root.key {
		root.leftChild = deleteNode(root.leftChild, targetKey)
	} else {
		if root.leftChild == nil {
			return root.rightChild
		} else if root.rightChild == nil {
			return root.leftChild
		} else {
			minNode := minValueNode(root.rightChild)
			root.key = minNode.key
			root.value = minNode.value
			root.rightChild = deleteNode(root.rightChild, minNode.key)
		}
	}

	return root
}

func TestOrderedMap(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())
	data.Insert(10, 10)
	data.Insert(5, 5)

	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

func TestOrderedMap_Empty(t *testing.T) {
	data := NewOrderedMap()

	assert.Equal(t, 0, data.Size())
	assert.False(t, data.Contains(10))

	data.Erase(10)

	assert.Equal(t, 0, data.Size())

	var keys []int
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.Empty(t, keys)
}

func TestOrderedMap_InsertDuplicate(t *testing.T) {
	data := NewOrderedMap()

	data.Insert(10, 100)

	assert.Panics(t, func() {
		data.Insert(10, 200)
	})

	assert.Equal(t, 1, data.Size())
	assert.True(t, data.Contains(10))

	var values []int
	data.ForEach(func(_, value int) {
		values = append(values, value)
	})

	assert.Equal(t, []int{100}, values)
}

func TestOrderedMap_EraseNonExisting(t *testing.T) {
	data := NewOrderedMap()

	data.Insert(10, 100)
	data.Insert(5, 50)
	data.Insert(15, 150)

	data.Erase(999)

	assert.Equal(t, 3, data.Size())

	var keys []int
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.Equal(t, []int{5, 10, 15}, keys)
}

func TestOrderedMap_EraseSingleElement(t *testing.T) {
	data := NewOrderedMap()

	data.Insert(10, 100)
	assert.Equal(t, 1, data.Size())

	data.Erase(10)

	assert.Equal(t, 0, data.Size())
	assert.False(t, data.Contains(10))

	var keys []int
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.Empty(t, keys)
}

func TestOrderedMap_EraseRoot(t *testing.T) {
	data := NewOrderedMap()

	data.Insert(10, 100)
	data.Insert(5, 50)
	data.Insert(15, 150)

	data.Erase(10)

	assert.Equal(t, 2, data.Size())
	assert.False(t, data.Contains(10))
	assert.True(t, data.Contains(5))
	assert.True(t, data.Contains(15))

	var keys []int
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.Equal(t, []int{5, 15}, keys)
}

func TestOrderedMap_EraseNodeWithTwoChildren(t *testing.T) {
	data := NewOrderedMap()

	data.Insert(10, 100)
	data.Insert(5, 50)
	data.Insert(15, 150)
	data.Insert(12, 120)
	data.Insert(20, 200)

	data.Erase(15)

	assert.Equal(t, 4, data.Size())

	assert.True(t, data.Contains(10))
	assert.True(t, data.Contains(5))
	assert.True(t, data.Contains(12))
	assert.True(t, data.Contains(20))
	assert.False(t, data.Contains(15))

	var keys []int
	var values []int

	data.ForEach(func(key, value int) {
		keys = append(keys, key)
		values = append(values, value)
	})

	assert.Equal(t, []int{5, 10, 12, 20}, keys)
	assert.Equal(t, []int{50, 100, 120, 200}, values)
}

func TestOrderedMap_ForEachValues(t *testing.T) {
	data := NewOrderedMap()

	data.Insert(10, 1000)
	data.Insert(5, 500)
	data.Insert(15, 1500)
	data.Insert(2, 200)

	var keys []int
	var values []int

	data.ForEach(func(key, value int) {
		keys = append(keys, key)
		values = append(values, value)
	})

	assert.Equal(t, []int{2, 5, 10, 15}, keys)
	assert.Equal(t, []int{200, 500, 1000, 1500}, values)
}
