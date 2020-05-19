package falsesharing

import (
	"sync"
	"testing"
)

const M = 1_000_000
const CacheLinePadSize = 64

type CacheLinePad struct{ _ [CacheLinePadSize]byte }

type A struct {
	n int
}

type PaddedA struct {
	n int
	_ CacheLinePad
}

type B struct {
	n int
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
				structA.n += j
			}
			wg.Done()
		}()
		go func() {
			for j := 0; j < M; j++ {
				structB.n += j
			}
			wg.Done()
		}()
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkPadding(b *testing.B) {
	structA := PaddedA{}
	structB := B{}
	wg := sync.WaitGroup{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go func() {
			for j := 0; j < M; j++ {
				structA.n += j
			}
			wg.Done()
		}()
		go func() {
			for j := 0; j < M; j++ {
				structB.n += j
			}
			wg.Done()
		}()
		wg.Wait()
	}
	b.StopTimer()
}
