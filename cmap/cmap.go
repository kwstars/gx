package cmap

// Map describes the common operations for a generic concurrent-safe map.
type Map[K comparable, V any] interface {
	// Load retrieves the value for key, returning ok=false if the key is absent.
	Load(key K) (value V, ok bool)
	// Store sets the value for key, replacing any existing entry.
	Store(key K, value V)
	// LoadOrStore returns the existing value if present; otherwise, it stores and returns the given value.
	LoadOrStore(key K, value V) (actual V, loaded bool)
	// LoadAndDelete removes the key and returns its previous value if it existed.
	LoadAndDelete(key K) (value V, loaded bool)
	// Delete removes the key without returning the previous value.
	Delete(key K)
	// Range iterates over all key/value pairs until the provided function returns false.
	Range(func(key K, value V) bool)
	// Len reports the number of key/value pairs currently in the map.
	Len() int
}
