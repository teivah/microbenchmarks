package falsesharing

import (
	"sync"
	"testing"
)

const M = 1_000_000
const CacheLinePadSize = 64

type CacheLinePad struct{ _ [CacheLinePadSize]byte }

type A struct {
	foo int
}

type APadded struct {
	foo int
	_   CacheLinePad
}

type B struct {
	foo int
}

// Sink makes sure results are not optimized away by the compiler
var Sink int

var SinkMu sync.Mutex

func BenchmarkFalseSharing(b *testing.B) {
	structA := A{}
	structB := B{}
	wg := sync.WaitGroup{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go func() {
			for j := 0; j < M; j++ {
				structA.foo += j
			}
			wg.Done()
		}()
		go func() {
			for j := 0; j < M; j++ {
				structB.foo += j
			}
			wg.Done()
		}()
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkPadding(b *testing.B) {
	structA := APadded{}
	structB := B{}
	wg := sync.WaitGroup{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go func() {
			for j := 0; j < M; j++ {
				structA.foo += j
			}
			wg.Done()
		}()
		go func() {
			for j := 0; j < M; j++ {
				structB.foo += j
			}
			wg.Done()
		}()
		wg.Wait()
	}
	b.StopTimer()
}
