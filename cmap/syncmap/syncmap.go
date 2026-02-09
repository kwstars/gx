package syncmap

import (
	"sync"

	"github.com/kwstars/gx/cmap"
)

// syncMap implements cmap.Map by wrapping sync.Map with typed helpers.
// Note: Len() is approximate and computed by iteration since sync.Map
// does not provide atomic length tracking.
type syncMap[K comparable, V any] struct {
	store sync.Map
}

// Ensure syncMap satisfies the cmap.Map interface at compile time.
var _ cmap.Map[int, int] = (*syncMap[int, int])(nil)

// New returns a cmap.Map backed by sync.Map.
// The returned map is safe for concurrent use by multiple goroutines.
func New[K comparable, V any]() cmap.Map[K, V] {
	return &syncMap[K, V]{}
}

// Load retrieves the value for key.
// Returns (zeroValue, false) if key does not exist or m is nil.
func (m *syncMap[K, V]) Load(key K) (value V, ok bool) {
	if m == nil {
		var zero V
		return zero, false
	}
	raw, exists := m.store.Load(key)
	if !exists {
		var zero V
		return zero, false
	}
	return raw.(V), true
}

// Store sets the value for key, replacing any existing entry.
// Safe for concurrent use. Does not affect length tracking (not maintained).
func (m *syncMap[K, V]) Store(key K, value V) {
	if m == nil {
		return
	}
	m.store.Store(key, value)
}

// LoadOrStore returns the existing value if present; otherwise stores and returns value.
// Safe for concurrent use. Only one goroutine will store the value for a given key.
func (m *syncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	if m == nil {
		return value, false
	}
	raw, ok := m.store.LoadOrStore(key, value)
	return raw.(V), ok
}

// LoadAndDelete removes the key and returns its previous value.
// Returns (zeroValue, false) if key does not exist or m is nil.
func (m *syncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	if m == nil {
		var zero V
		return zero, false
	}
	raw, ok := m.store.LoadAndDelete(key)
	if !ok {
		var zero V
		return zero, false
	}
	return raw.(V), true
}

// Delete removes the key without returning the previous value.
// Safe for concurrent use. No-op if key does not exist or m is nil.
func (m *syncMap[K, V]) Delete(key K) {
	if m == nil {
		return
	}
	m.store.Delete(key)
}

// Range iterates over the map until the provided function returns false.
// The iteration is safe for concurrent use, but the map may be modified
// during iteration. The function fn must not modify the map.
func (m *syncMap[K, V]) Range(fn func(key K, value V) bool) {
	if m == nil {
		return
	}
	m.store.Range(func(k, v any) bool {
		return fn(k.(K), v.(V))
	})
}

// Len reports an approximate number of key/value pairs in the map.
// This is computed by iterating over the map and may not reflect concurrent modifications.
// For exact counts, use a different data structure (e.g., sharded map with atomic counters).
func (m *syncMap[K, V]) Len() int {
	if m == nil {
		return 0
	}
	count := 0
	m.Range(func(_ K, _ V) bool {
		count++
		return true
	})
	return count
}
