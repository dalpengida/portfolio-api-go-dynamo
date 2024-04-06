package common

import (
	"strings"
	"testing"
)

var result []T

const size = 10000

type T int

// BenchmarkCopy copy를 이용한 slice 복사 benchmark
func BenchmarkCopy(b *testing.B) {
	orig := make([]T, size)

	for n := 0; n < b.N; n++ {
		cpy := make([]T, len(orig))
		copy(cpy, orig)
		orig = cpy
	}
	result = orig
}

// BenchmarkAppend 빈 배열 구조체만으로 append 하여 복사 benchmark
func BenchmarkAppend(b *testing.B) {
	orig := make([]T, size)

	for n := 0; n < b.N; n++ {
		cpy := append([]T{}, orig...)
		orig = cpy
	}
	result = orig
}

// BenchmarkAppendPreCapped 빈 배열에 capacity 값을 옮겨 담을 경우 benchmark
func BenchmarkAppendPreCapped(b *testing.B) {
	orig := make([]T, size)
	for n := 0; n < b.N; n++ {
		cpy := append(make([]T, 0, len(orig)), orig...)
		orig = cpy
	}
	result = orig
}

// BenchmarkAppendNil nil 을 append 하여 값 복사 benchmark
func BenchmarkAppendNil(b *testing.B) {
	orig := make([]T, size)
	for n := 0; n < b.N; n++ {
		cpy := append([]T(nil), orig...)
		orig = cpy
	}
	result = orig
}

// joinWithPlus 무식하게 += 에 의한 구현
func joinWithPlus(strs ...string) string {
	var ret string
	for _, str := range strs {
		ret += str
	}
	return ret
}

// joinWithBuilder Builder.Builder 에 의한 구현
func joinWithBuilder(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

// BenchmarkPlus 무식하게  += 에 의한 구현 벤치마크
func BenchmarkPlus(b *testing.B) {
	strs := []string{"aaa", "bbb", "ccc", "ddd", "eee", "fff", "ggg", "hhh", "iii"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = joinWithPlus(strs...)
	}
}

// BenchmarkBuilder Builder.Builder 에 의한 구현 벤치마크
func BenchmarkBuilder(b *testing.B) {
	strs := []string{"aaa", "bbb", "ccc", "ddd", "eee", "fff", "ggg", "hhh", "iii"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = joinWithBuilder(strs...)
	}
}
