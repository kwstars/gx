package rwmap

import (
	"sync"
	"testing"
)

func TestRWMapStoreLoad(t *testing.T) {
	t.Parallel()

	m := New[string, int]()

	if _, ok := m.Load("missing"); ok {
		t.Fatalf("expected missing key to return ok=false")
	}

	m.Store("foo", 1)
	if got, ok := m.Load("foo"); !ok || got != 1 {
		t.Fatalf("expected foo=1, got %v ok=%v", got, ok)
	}

	if gotLen := m.Len(); gotLen != 1 {
		t.Fatalf("expected len=1, got %d", gotLen)
	}

	m.Store("foo", 2)
	if got, _ := m.Load("foo"); got != 2 {
		t.Fatalf("expected foo=2 after overwrite, got %d", got)
	}

	if gotLen := m.Len(); gotLen != 1 {
		t.Fatalf("expected len to remain 1 after overwrite, got %d", gotLen)
	}
}

func TestRWMapLoadOrStoreAndDelete(t *testing.T) {
	t.Parallel()

	m := New[int, string]()

	if actual, loaded := m.LoadOrStore(1, "a"); loaded || actual != "a" {
		t.Fatalf("expected store to insert new value, got %q loaded=%v", actual, loaded)
	}

	if actual, loaded := m.LoadOrStore(1, "b"); !loaded || actual != "a" {
		t.Fatalf("expected load of existing value, got %q loaded=%v", actual, loaded)
	}

	if gotLen := m.Len(); gotLen != 1 {
		t.Fatalf("expected len=1, got %d", gotLen)
	}

	if val, loaded := m.LoadAndDelete(1); !loaded || val != "a" {
		t.Fatalf("expected delete to return stored value, got %q loaded=%v", val, loaded)
	}

	if _, loaded := m.LoadAndDelete(1); loaded {
		t.Fatalf("expected second delete to report loaded=false")
	}

	if gotLen := m.Len(); gotLen != 0 {
		t.Fatalf("expected len=0 after delete, got %d", gotLen)
	}
}

func TestRWMapRangeAndConcurrency(t *testing.T) {
	t.Parallel()

	m := New[int, int]()
	const total = 128

	for i := 0; i < total; i++ {
		m.Store(i, i*i)
	}

	seen := make(map[int]int)
	m.Range(func(k, v int) bool {
		seen[k] = v
		return len(seen) < 10
	})

	if len(seen) != 10 {
		t.Fatalf("expected range to stop after 10 iterations, got %d", len(seen))
	}

	done := make(chan struct{})
	go func() {
		m.Range(func(k, v int) bool {
			m.Store(total+k, v)
			return true
		})
		close(done)
	}()

	<-done
}

func TestRWMapConcurrentAccess(t *testing.T) {
	t.Parallel()

	m := New[int, int]()
	const total = 256

	var wg sync.WaitGroup
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Store(i, i)
		}(i)
	}

	wg.Wait()

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Load(i)
		}(i)
	}

	wg.Wait()

	if got := m.Len(); got != total {
		t.Fatalf("expected len=%d after concurrent writes, got %d", total, got)
	}
}
