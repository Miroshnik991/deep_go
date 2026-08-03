package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go
// go test homework_test.go -bench=. -benchmem -run=^$

func ToLittleEndian(number uint32) uint32 {
	return (number >> 24) | (number << 24) | ((number >> 8) & 0xFF00) | ((number << 8) & 0xFF0000)
}

func ToLittleEndian2(number uint32) uint32 {
	var result [4]byte
	numPtr := unsafe.Pointer(&number)
	padding := 3
	for i := 0; i < 4; i++ {
		result[i] = *(*uint8)(unsafe.Add(numPtr, padding))
		padding--
	}
	return *(*uint32)(unsafe.Pointer(&result))
}

func ToLittleEndian3(number uint32) uint32 {
	for i := 0; i < 2; i++ {
		s := unsafe.Add(unsafe.Pointer(&number), i)
		e := unsafe.Add(unsafe.Pointer(&number), (int)(unsafe.Sizeof(number))-1-i)
		*(*int8)(s), *(*int8)(e) = *(*int8)(e), *(*int8)(s)
	}
	return number
}

func TestConversion(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}

// implementation of benchmark aggregation
//func BenchmarkToLittleEndian(b *testing.B) {
//	inputs := getInputsForBenchmark()
//
//	benchmarks := []struct {
//		name string
//		f    func(uint32) uint32
//	}{
//		{"V1", ToLittleEndian},
//		{"V2", ToLittleEndian2},
//		{"V3", ToLittleEndian3},
//	}
//
//	for _, bmk := range benchmarks {
//		b.Run(bmk.name, func(b *testing.B) {
//			b.ReportAllocs()
//			for i := 0; i < b.N; i++ {
//				for _, num := range inputs {
//					_ = bmk.f(num)
//				}
//			}
//		})
//	}
//}

func BenchmarkToLittleEndian(b *testing.B) {
	inputs := getInputsForBenchmark()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, num := range inputs {
			_ = ToLittleEndian(num)
		}
	}
}

func BenchmarkToLittleEndian2(b *testing.B) {
	inputs := getInputsForBenchmark()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, num := range inputs {
			_ = ToLittleEndian2(num)
		}
	}
}

func BenchmarkToLittleEndian3(b *testing.B) {
	inputs := getInputsForBenchmark()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, num := range inputs {
			_ = ToLittleEndian3(num)
		}
	}
}

func getInputsForBenchmark() []uint32 {
	return []uint32{
		0x00000000,
		0xFFFFFFFF,
		0x00FF00FF,
		0x0000FFFF,
		0x01020304,
		0x12345678,
		0xAABBCCDD,
	}
}
