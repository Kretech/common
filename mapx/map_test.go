package mapx

import (
	"fmt"
	"testing"
)

var size = 8
var m = MakeHashMap(size)
var s = MakeSliceMap(size)
var name = make([]string, size)

func TestMain(tm *testing.M) {
	for i := 0; i < size; i++ {
		m.Store(i%size, i+1)
		s.Store(i%size, i+1)
		name[i%size] = fmt.Sprint(i + 1)
	}
	tm.Run()
}

func BenchmarkMapLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m.Load(name[i%size])
	}
}

func BenchmarkSliceLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s.Load(name[i%size])
	}
}

func BenchmarkMapStore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m.Store(name[i%size], i%size)
	}
}

func BenchmarkSliceStore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s.Store(name[i%size], i%size)
	}
}
