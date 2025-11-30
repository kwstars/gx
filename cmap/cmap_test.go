package cmap_test

import (
	"sync/atomic"
	"testing"

	"github.com/kwstars/gx/cmap"
	"github.com/kwstars/gx/cmap/rwmap"
	"github.com/kwstars/gx/cmap/syncmap"
)

const benchKeySpace = 1024

var benchSink atomic.Int64

var benchFactories = []struct {
	name    string
	factory func() cmap.Map[int, int]
}{
	{
		name: "rwmap",
		factory: func() cmap.Map[int, int] {
			return rwmap.New[int, int]()
		},
	},
	{
		name: "syncmap",
		factory: func() cmap.Map[int, int] {
			return syncmap.New[int, int]()
		},
	},
}

func BenchmarkMapStore(b *testing.B) {
	for _, tc := range benchFactories {
		b.Run(tc.name, func(b *testing.B) {
			benchmarkStore(b, tc.factory)
		})
	}
}

func BenchmarkMapLoad(b *testing.B) {
	for _, tc := range benchFactories {
		b.Run(tc.name, func(b *testing.B) {
			benchmarkLoad(b, tc.factory)
		})
	}
}

func BenchmarkMapLoadOrStore(b *testing.B) {
	for _, tc := range benchFactories {
		b.Run(tc.name, func(b *testing.B) {
			benchmarkLoadOrStore(b, tc.factory)
		})
	}
}

func BenchmarkMapRange(b *testing.B) {
	for _, tc := range benchFactories {
		b.Run(tc.name, func(b *testing.B) {
			benchmarkRange(b, tc.factory)
		})
	}
}

func benchmarkStore(b *testing.B, factory func() cmap.Map[int, int]) {
	b.Helper()
	m := factory()
	mask := benchKeySpace - 1

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		local := 0
		for pb.Next() {
			key := local & mask
			m.Store(key, local)
			local++
		}
	})
}

func benchmarkLoad(b *testing.B, factory func() cmap.Map[int, int]) {
	b.Helper()
	m := factory()
	for i := 0; i < benchKeySpace; i++ {
		m.Store(i, i)
	}
	mask := benchKeySpace - 1

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		local := 0
		for pb.Next() {
			key := local & mask
			val, _ := m.Load(key)
			benchSink.Add(int64(val))
			local++
		}
	})
}

func benchmarkLoadOrStore(b *testing.B, factory func() cmap.Map[int, int]) {
	b.Helper()
	m := factory()
	mask := benchKeySpace - 1

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		local := 0
		for pb.Next() {
			key := local & mask
			val, _ := m.LoadOrStore(key, local)
			benchSink.Add(int64(val))
			local++
		}
	})
}

func benchmarkRange(b *testing.B, factory func() cmap.Map[int, int]) {
	b.Helper()
	m := factory()
	for i := 0; i < benchKeySpace; i++ {
		m.Store(i, i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Range(func(key, value int) bool {
				benchSink.Add(int64(value))
				return true
			})
		}
	})
}
