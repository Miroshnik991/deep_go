package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type SharedBuffer struct {
	data []byte
	refs int
}

type COWBuffer struct {
	sharedData *SharedBuffer
}

func NewCOWBuffer(data []byte) COWBuffer {
	return COWBuffer{
		&SharedBuffer{
			data: data,
			refs: 1,
		},
	}
}

func (b *COWBuffer) Clone() COWBuffer {
	if b.sharedData == nil {
		panic("clone of closed buffer")
	}
	b.sharedData.refs++
	return COWBuffer{
		b.sharedData,
	}
}

func (b *COWBuffer) Close() {
	if b.sharedData == nil {
		return
	}
	b.sharedData.refs--
	b.sharedData = nil
}

func (b *COWBuffer) Update(index int, value byte) bool {
	if b.sharedData == nil {
		return false
	}
	if index < 0 || index >= len(b.sharedData.data) {
		return false
	}
	if b.sharedData.refs > 1 {
		copyData := make([]byte, len(b.sharedData.data))
		copy(copyData, b.sharedData.data)
		copyData[index] = value
		b.sharedData.refs--
		b.sharedData = &SharedBuffer{
			data: copyData,
			refs: 1,
		}
		return true
	}
	b.sharedData.data[index] = value
	return true
}

func (b *COWBuffer) String() string {
	if b.sharedData == nil {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b.sharedData.data), len(b.sharedData.data))
}

func TestCOWBuffer(t *testing.T) {
	original := []byte{'a', 'b', 'c', 'd'}

	buffer := NewCOWBuffer(original)
	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	defer buffer.Close()
	defer copy1.Close()
	defer copy2.Close()

	assert.Same(t, buffer.sharedData, copy1.sharedData)
	assert.Same(t, copy1.sharedData, copy2.sharedData)

	assert.Equal(t, 3, buffer.sharedData.refs)

	assert.Equal(t,
		unsafe.SliceData(buffer.sharedData.data),
		unsafe.SliceData(copy1.sharedData.data),
	)
	assert.Equal(t,
		unsafe.SliceData(copy1.sharedData.data),
		unsafe.SliceData(copy2.sharedData.data),
	)

	assert.Same(t,
		unsafe.StringData(buffer.String()),
		unsafe.SliceData(buffer.sharedData.data),
	)
	assert.Same(t,
		unsafe.StringData(copy1.String()),
		unsafe.SliceData(copy1.sharedData.data),
	)
	assert.Same(t,
		unsafe.StringData(copy2.String()),
		unsafe.SliceData(copy2.sharedData.data),
	)

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	assert.NotSame(t, buffer.sharedData, copy1.sharedData)
	assert.Same(t, copy1.sharedData, copy2.sharedData)

	assert.Equal(t, 1, buffer.sharedData.refs)
	assert.Equal(t, 2, copy1.sharedData.refs)

	assert.Equal(t, []byte{'g', 'b', 'c', 'd'}, buffer.sharedData.data)
	assert.Equal(t, []byte{'a', 'b', 'c', 'd'}, copy1.sharedData.data)
	assert.Equal(t, []byte{'a', 'b', 'c', 'd'}, copy2.sharedData.data)

	assert.NotEqual(t,
		unsafe.SliceData(buffer.sharedData.data),
		unsafe.SliceData(copy1.sharedData.data),
	)

	assert.Equal(t,
		unsafe.SliceData(copy1.sharedData.data),
		unsafe.SliceData(copy2.sharedData.data),
	)

	copy1.Close()

	assert.Equal(t, 1, copy2.sharedData.refs)

	before := unsafe.SliceData(copy2.sharedData.data)

	assert.True(t, copy2.Update(0, 'f'))

	after := unsafe.SliceData(copy2.sharedData.data)

	assert.Equal(t, before, after)
	assert.Equal(t, []byte{'f', 'b', 'c', 'd'}, copy2.sharedData.data)
}
