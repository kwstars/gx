package rwmap

import (
	"sync"

	"github.com/kwstars/gx/cmap"
)

// rwMap implements cmap.Map using sync.RWMutex + built-in map.
// All methods assume m != nil. Calling methods on nil *rwMap will panic,
// consistent with Go standard library conventions (e.g., sync.Map).
type rwMap[K comparable, V any] struct {
	mu    sync.RWMutex
	store map[K]V // Guaranteed non-nil by newMap()
}

// Ensure rwMap obeys cmap.Map interface at compile time.
var _ cmap.Map[int, int] = (*rwMap[int, int])(nil)

// New returns a cmap.Map implementation backed by RWMutex.
func New[K comparable, V any]() cmap.Map[K, V] {
	return newMap[K, V]()
}

// newMap exposes concrete type for callers needing assertions in tests.
func newMap[K comparable, V any]() *rwMap[K, V] {
	return &rwMap[K, V]{
		store: make(map[K]V), // Always initialized; never nil
	}
}

// Load retrieves the value for key, returning ok=false when missing.
func (m *rwMap[K, V]) Load(key K) (value V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok = m.store[key]
	return value, ok
}

// Store sets the value for key, overwriting any existing value.
func (m *rwMap[K, V]) Store(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[key] = value
}

// LoadOrStore returns the existing value if present, storing otherwise.
// Uses write lock to ensure atomic "check-then-act" semantics.
func (m *rwMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, ok := m.store[key]; ok {
		return existing, true
	}
	m.store[key] = value
	return value, false
}

// LoadAndDelete removes key and returns prior value if it existed.
func (m *rwMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, loaded = m.store[key]
	if loaded {
		delete(m.store, key)
	}
	return value, loaded
}

// Delete removes the key without reporting previous value.
func (m *rwMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, key)
}

// Range iterates over a snapshot of the map entries until fn returns false.
// Snapshot semantics:
//   - Entries seen are those present when Range acquired the read lock
//   - Modifications during iteration do not affect the snapshot
//   - Read lock is held only during snapshot creation (not during fn execution)
func (m *rwMap[K, V]) Range(fn func(key K, value V) bool) {
	if fn == nil {
		return
	}

	m.mu.RLock()
	// Create snapshot under read lock to avoid long-term lock contention
	snapshot := make([]struct {
		key K
		val V
	}, 0, len(m.store))
	for k, v := range m.store {
		snapshot = append(snapshot, struct {
			key K
			val V
		}{k, v})
	}
	m.mu.RUnlock()

	// Execute user callback without holding any locks
	for _, item := range snapshot {
		if !fn(item.key, item.val) {
			return
		}
	}
}

// Len reports the number of key/value pairs in the map.
func (m *rwMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.store)
}
