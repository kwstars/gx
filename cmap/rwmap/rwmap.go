package rwmap

import (
	"sync"

	"github.com/kwstars/gx/cmap"
)

// rwMap implements cmap.Map using sync.RWMutex + built-in map.
type rwMap[K comparable, V any] struct {
	mu    sync.RWMutex
	store map[K]V
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
		store: make(map[K]V),
	}
}

func (m *rwMap[K, V]) ensureStore() {
	if m.store == nil {
		m.store = make(map[K]V)
	}
}

// Load retrieves the value for key, returning ok=false when missing.
func (m *rwMap[K, V]) Load(key K) (value V, ok bool) {
	if m == nil {
		var zero V
		return zero, false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.store == nil {
		var zero V
		return zero, false
	}
	value, ok = m.store[key]
	return value, ok
}

// Store sets the value for key, overwriting any existing value.
func (m *rwMap[K, V]) Store(key K, value V) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ensureStore()
	m.store[key] = value
}

// LoadOrStore returns the existing value if present, storing otherwise.
func (m *rwMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	if m == nil {
		return value, false
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ensureStore()
	if existing, ok := m.store[key]; ok {
		return existing, true
	}
	m.store[key] = value
	return value, false
}

// LoadAndDelete removes key and returns prior value if it existed.
func (m *rwMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	if m == nil {
		var zero V
		return zero, false
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.store == nil {
		var zero V
		return zero, false
	}
	value, loaded = m.store[key]
	if loaded {
		delete(m.store, key)
	}
	return value, loaded
}

// Delete removes the key without reporting previous value.
func (m *rwMap[K, V]) Delete(key K) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.store != nil {
		delete(m.store, key)
	}
}

// Range iterates over entries until fn returns false.
func (m *rwMap[K, V]) Range(fn func(key K, value V) bool) {
	if m == nil || fn == nil {
		return
	}

	m.mu.RLock()
	if len(m.store) == 0 {
		m.mu.RUnlock()
		return
	}

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

	for _, item := range snapshot {
		if !fn(item.key, item.val) {
			return
		}
	}
}

// Len reports the number of key/value pairs in the map.
func (m *rwMap[K, V]) Len() int {
	if m == nil {
		return 0
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.store)
}
