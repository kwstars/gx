package syncmap

import (
	"sync"
	"sync/atomic"

	"github.com/kwstars/gx/cmap"
)

// syncMap implements cmap.syncMap by wrapping sync.syncMap with typed helpers.
type syncMap[K comparable, V any] struct {
	store sync.Map
	len   atomic.Int64
}

// Ensure Map satisfies the cmap.Map interface at compile time.
var _ cmap.Map[int, int] = (*syncMap[int, int])(nil)

// New returns a cmap.Map backed by sync.Map.
func New[K comparable, V any]() cmap.Map[K, V] {
	return newMap[K, V]()
}

// newMap returns the concrete *Map for callers that need type assertions (e.g. tests).
func newMap[K comparable, V any]() *syncMap[K, V] {
	return &syncMap[K, V]{}
}

// Load retrieves the value for key.
func (m *syncMap[K, V]) Load(key K) (value V, ok bool) {
	if m == nil {
		var zero V
		return zero, false
	}
	if raw, exists := m.store.Load(key); exists {
		return raw.(V), true
	}
	var zero V
	return zero, false
}

// Store sets the value for key, replacing any existing entry.
func (m *syncMap[K, V]) Store(key K, value V) {
	if m == nil {
		return
	}
	if _, loaded := m.store.LoadOrStore(key, value); loaded {
		m.store.Store(key, value)
		return
	}
	m.len.Add(1)
}

// LoadOrStore returns the existing value if present; otherwise stores and returns value.
func (m *syncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	if m == nil {
		return value, false
	}
	raw, ok := m.store.LoadOrStore(key, value)
	if !ok {
		m.len.Add(1)
	}
	return raw.(V), ok
}

// LoadAndDelete removes the key and returns its previous value.
func (m *syncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	if m == nil {
		var zero V
		return zero, false
	}
	raw, ok := m.store.LoadAndDelete(key)
	if ok {
		m.len.Add(-1)
		return raw.(V), true
	}
	var zero V
	return zero, false
}

// Delete removes the key without returning the previous value.
func (m *syncMap[K, V]) Delete(key K) {
	if m == nil {
		return
	}
	if _, ok := m.store.LoadAndDelete(key); ok {
		m.len.Add(-1)
	}
}

// Range iterates over the map until the provided function returns false.
func (m *syncMap[K, V]) Range(fn func(key K, value V) bool) {
	if m == nil {
		return
	}
	m.store.Range(func(k, v any) bool {
		return fn(k.(K), v.(V))
	})
}

// Len reports the number of key/value pairs currently in the map.
func (m *syncMap[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return int(m.len.Load())
}
