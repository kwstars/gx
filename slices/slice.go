package slices

import (
	"math/rand"
	"sort"
)

// Slice represents a generic slice container with built-in operations
type Slice[T any] struct {
	data []T
}

// NewSlice creates a new Slice instance from an existing slice
// Creates a copy to prevent external modifications from affecting internal data
func NewSlice[T any](data []T) *Slice[T] {
	copied := make([]T, len(data))
	copy(copied, data)
	return &Slice[T]{copied}
}

// From creates a Slice from variadic arguments
func From[T any](items ...T) *Slice[T] {
	copied := make([]T, len(items))
	copy(copied, items)
	return &Slice[T]{copied}
}

// Data returns a copy of the underlying slice to prevent external modification
func (s *Slice[T]) Data() []T {
	result := make([]T, len(s.data))
	copy(result, s.data)
	return result
}

// DataUnsafe returns the underlying slice directly without copying
// Use with caution as modifications will affect the internal state
func (s *Slice[T]) DataUnsafe() []T {
	return s.data
}

// Len returns the length of the slice
func (s *Slice[T]) Len() int {
	return len(s.data)
}

// Cap returns the capacity of the slice
func (s *Slice[T]) Cap() int {
	return cap(s.data)
}

// IsEmpty checks if the slice is empty
func (s *Slice[T]) IsEmpty() bool {
	return len(s.data) == 0
}

// Get retrieves an element at the specified index
// Returns the element and true if index is valid, zero value and false otherwise
func (s *Slice[T]) Get(index int) (T, bool) {
	if index < 0 || index >= len(s.data) {
		var zero T
		return zero, false
	}
	return s.data[index], true
}

// Set sets an element at the specified index
// Returns true if successful, false if index is out of bounds
func (s *Slice[T]) Set(index int, value T) bool {
	if index < 0 || index >= len(s.data) {
		return false
	}
	s.data[index] = value
	return true
}

// AppendSlice appends all elements from another Slice
// Renamed from AppendVector for clearer semantics
func (s *Slice[T]) AppendSlice(other *Slice[T]) *Slice[T] {
	if other == nil {
		return s
	}
	otherData := other.Data() // Use copy to ensure safety
	s.data = append(s.data, otherData...)
	return s
}

// Append adds elements to the end
func (s *Slice[T]) Append(values ...T) *Slice[T] {
	s.data = append(s.data, values...)
	return s
}

// Copy creates a deep copy of the current slice
func (s *Slice[T]) Copy() *Slice[T] {
	newData := make([]T, len(s.data))
	copy(newData, s.data)
	return &Slice[T]{newData}
}

// Cut removes elements from index i to j (exclusive: [i, j))
// Returns the modified slice for chaining
func (s *Slice[T]) Cut(i, j int) *Slice[T] {
	if i < 0 || j > len(s.data) || i > j {
		return s
	}
	s.data = append(s.data[:i], s.data[j:]...)
	return s
}

// CutWithCleanup removes elements from index i to j and zeros out references for GC
// This is useful for slices containing pointers to prevent memory leaks
func (s *Slice[T]) CutWithCleanup(i, j int) *Slice[T] {
	if i < 0 || j > len(s.data) || i > j {
		return s
	}

	copy(s.data[i:], s.data[j:])
	// Zero out the tail elements to help GC
	for k := len(s.data) - (j - i); k < len(s.data); k++ {
		var zero T
		s.data[k] = zero
	}
	s.data = s.data[:len(s.data)-(j-i)]
	return s
}

// Delete removes the element at index i
func (s *Slice[T]) Delete(i int) *Slice[T] {
	if i < 0 || i >= len(s.data) {
		return s
	}
	s.data = append(s.data[:i], s.data[i+1:]...)
	return s
}

// DeleteWithCleanup removes element at index i and zeros the reference for GC
func (s *Slice[T]) DeleteWithCleanup(i int) *Slice[T] {
	if i < 0 || i >= len(s.data) {
		return s
	}

	copy(s.data[i:], s.data[i+1:])
	var zero T
	s.data[len(s.data)-1] = zero
	s.data = s.data[:len(s.data)-1]
	return s
}

// DeleteUnordered deletes element at index i without preserving order (faster than Delete)
// The last element is moved to position i
func (s *Slice[T]) DeleteUnordered(i int) *Slice[T] {
	n := len(s.data)
	if n == 0 || i < 0 || i >= n {
		return s
	}
	s.data[i] = s.data[n-1]
	s.data = s.data[:n-1]
	return s
}

// DeleteUnorderedWithCleanup deletes without preserving order and zeros the reference
func (s *Slice[T]) DeleteUnorderedWithCleanup(i int) *Slice[T] {
	n := len(s.data)
	if n == 0 || i < 0 || i >= n {
		return s
	}
	s.data[i] = s.data[n-1]
	var zero T
	s.data[n-1] = zero
	s.data = s.data[:n-1]
	return s
}

// Expand inserts n zero-value elements at position i
func (s *Slice[T]) Expand(i, n int) *Slice[T] {
	if i < 0 || i > len(s.data) || n < 0 {
		return s
	}
	if n == 0 {
		return s
	}
	s.data = append(s.data[:i], append(make([]T, n), s.data[i:]...)...)
	return s
}

// Extend appends n zero-value elements to the end
func (s *Slice[T]) Extend(n int) *Slice[T] {
	if n <= 0 {
		return s
	}
	s.data = append(s.data, make([]T, n)...)
	return s
}

// Filter filters elements in place based on predicate function
// Elements that satisfy keep() are retained
func (s *Slice[T]) Filter(keep func(T) bool) *Slice[T] {
	n := 0
	for _, x := range s.data {
		if keep(x) {
			s.data[n] = x
			n++
		}
	}
	// Zero out the tail to help GC for pointer types
	for i := n; i < len(s.data); i++ {
		var zero T
		s.data[i] = zero
	}
	s.data = s.data[:n]
	return s
}

// FilterWithoutAllocating filters without allocating new memory
// This is an alias for Filter with different implementation strategy
func (s *Slice[T]) FilterWithoutAllocating(f func(T) bool) *Slice[T] {
	b := s.data[:0]
	for _, x := range s.data {
		if f(x) {
			b = append(b, x)
		}
	}
	// Zero out unused elements
	for i := len(b); i < len(s.data); i++ {
		var zero T
		s.data[i] = zero
	}
	s.data = b
	return s
}

// Insert inserts an element at index i
func (s *Slice[T]) Insert(i int, x T) *Slice[T] {
	if i < 0 || i > len(s.data) {
		return s
	}
	s.data = append(s.data[:i], append([]T{x}, s.data[i:]...)...)
	return s
}

// InsertNoAlloc inserts an element at index i with minimal allocations
// More efficient than Insert for reducing intermediate slice allocations
func (s *Slice[T]) InsertNoAlloc(i int, x T) *Slice[T] {
	if i < 0 || i > len(s.data) {
		return s
	}
	var zero T
	s.data = append(s.data, zero)
	copy(s.data[i+1:], s.data[i:])
	s.data[i] = x
	return s
}

// Push adds an element to the end of the slice
func (s *Slice[T]) Push(x T) *Slice[T] {
	s.data = append(s.data, x)
	return s
}

// Pop removes and returns the last element
// Returns zero value and false if slice is empty
func (s *Slice[T]) Pop() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	x := s.data[len(s.data)-1]
	var zero T
	s.data[len(s.data)-1] = zero // Help GC
	s.data = s.data[:len(s.data)-1]
	return x, true
}

// PushFront adds an element to the beginning of the slice
func (s *Slice[T]) PushFront(x T) *Slice[T] {
	s.data = append([]T{x}, s.data...)
	return s
}

// PopFront removes and returns the first element
// Returns zero value and false if slice is empty
func (s *Slice[T]) PopFront() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	x := s.data[0]
	// Note: This doesn't zero out the removed element as the slice is re-sliced
	// If GC is critical, consider using copy and zeroing
	s.data = s.data[1:]
	return x, true
}

// Reverse reverses the slice in place
func (s *Slice[T]) Reverse() *Slice[T] {
	for left, right := 0, len(s.data)-1; left < right; left, right = left+1, right-1 {
		s.data[left], s.data[right] = s.data[right], s.data[left]
	}
	return s
}

// Shuffle randomizes the order of elements using Fisher-Yates algorithm
// Note: Uses global rand, consider passing *rand.Rand for better control
func (s *Slice[T]) Shuffle() *Slice[T] {
	for i := len(s.data) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s.data[i], s.data[j] = s.data[j], s.data[i]
	}
	return s
}

// ShuffleWithRand shuffles using a provided random generator for reproducibility
func (s *Slice[T]) ShuffleWithRand(r *rand.Rand) *Slice[T] {
	if r == nil {
		return s.Shuffle()
	}
	for i := len(s.data) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		s.data[i], s.data[j] = s.data[j], s.data[i]
	}
	return s
}

// Clear clears the slice and optionally zeros out elements for GC
func (s *Slice[T]) Clear(zeroOut bool) *Slice[T] {
	if zeroOut {
		for i := range s.data {
			var zero T
			s.data[i] = zero
		}
	}
	s.data = s.data[:0]
	return s
}

// ForEach applies a function to each element with its index
func (s *Slice[T]) ForEach(fn func(value T, index int)) {
	for i, v := range s.data {
		fn(v, i)
	}
}

// Map transforms each element in place using the provided function
func (s *Slice[T]) Map(fn func(T) T) *Slice[T] {
	for i, v := range s.data {
		s.data[i] = fn(v)
	}
	return s
}

// MapToNew creates a new slice by transforming each element
// The original slice remains unchanged
func (s *Slice[T]) MapToNew(fn func(T) T) *Slice[T] {
	result := make([]T, len(s.data))
	for i, v := range s.data {
		result[i] = fn(v)
	}
	return &Slice[T]{data: result}
}

// Reduce reduces the slice to a single value using an accumulator function
func (s *Slice[T]) Reduce(initial T, fn func(accumulator T, current T) T) T {
	result := initial
	for _, v := range s.data {
		result = fn(result, v)
	}
	return result
}

// Find finds the first element matching the predicate
// Returns the element and true if found, zero value and false otherwise
func (s *Slice[T]) Find(predicate func(T) bool) (T, bool) {
	for _, v := range s.data {
		if predicate(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex finds the index of the first element matching the predicate
// Returns -1 if not found
func (s *Slice[T]) FindIndex(predicate func(T) bool) int {
	for i, v := range s.data {
		if predicate(v) {
			return i
		}
	}
	return -1
}

// IndexOf finds the index of the first occurrence of a value using a custom equality function
// Returns -1 if not found
func (s *Slice[T]) IndexOf(value T, equal func(a, b T) bool) int {
	for i, v := range s.data {
		if equal(v, value) {
			return i
		}
	}
	return -1
}

// Contains checks if the slice contains a specific value using a custom equality function
func (s *Slice[T]) Contains(value T, equal func(a, b T) bool) bool {
	return s.IndexOf(value, equal) != -1
}

// All checks if all elements satisfy the predicate
func (s *Slice[T]) All(predicate func(T) bool) bool {
	for _, v := range s.data {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Any checks if any element satisfies the predicate
func (s *Slice[T]) Any(predicate func(T) bool) bool {
	for _, v := range s.data {
		if predicate(v) {
			return true
		}
	}
	return false
}

// Deduplicate removes duplicates using a custom comparator
// The slice is sorted as a side effect
// comparator should return: negative if a < b, zero if a == b, positive if a > b
func (s *Slice[T]) Deduplicate(comparator func(a, b T) int) *Slice[T] {
	if len(s.data) <= 1 {
		return s
	}

	// Sort the slice in place
	sort.Slice(s.data, func(i, j int) bool {
		return comparator(s.data[i], s.data[j]) < 0
	})

	// Remove duplicates
	j := 0
	for i := 1; i < len(s.data); i++ {
		if comparator(s.data[j], s.data[i]) != 0 {
			j++
			s.data[j] = s.data[i]
		}
	}

	// Zero out unused elements
	for k := j + 1; k < len(s.data); k++ {
		var zero T
		s.data[k] = zero
	}

	s.data = s.data[:j+1]
	return s
}

// DeduplicateStable removes duplicates while preserving order
// Uses a map-based approach (requires comparable types via the equal function)
func (s *Slice[T]) DeduplicateStable(equal func(a, b T) bool) *Slice[T] {
	if len(s.data) <= 1 {
		return s
	}

	j := 0

	for i := 0; i < len(s.data); i++ {
		isDuplicate := false
		for k := 0; k < j; k++ {
			if equal(s.data[k], s.data[i]) {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			s.data[j] = s.data[i]
			j++
		}
	}

	// Zero out unused elements
	for k := j; k < len(s.data); k++ {
		var zero T
		s.data[k] = zero
	}

	s.data = s.data[:j]
	return s
}

// Batch divides the slice into batches of specified size
// Returns a slice of Slice pointers, each containing a batch
func (s *Slice[T]) Batch(batchSize int) []*Slice[T] {
	if len(s.data) == 0 {
		return []*Slice[T]{}
	}

	if batchSize <= 0 {
		// Invalid batch size, return one batch with all elements
		return []*Slice[T]{s.Copy()}
	}

	var batches []*Slice[T]
	for i := 0; i < len(s.data); i += batchSize {
		end := i + batchSize
		if end > len(s.data) {
			end = len(s.data)
		}
		// Create a copy of the batch to avoid sharing underlying array
		batchData := make([]T, end-i)
		copy(batchData, s.data[i:end])
		batches = append(batches, &Slice[T]{batchData})
	}
	return batches
}

// SlidingWindow creates sliding windows of specified size
// Returns overlapping windows advancing one element at a time
func (s *Slice[T]) SlidingWindow(size int) []*Slice[T] {
	if size <= 0 || len(s.data) == 0 {
		return []*Slice[T]{}
	}

	if size > len(s.data) {
		return []*Slice[T]{s.Copy()}
	}

	var result []*Slice[T]
	for i := 0; i <= len(s.data)-size; i++ {
		// Create a copy of the window to avoid sharing underlying array
		windowData := make([]T, size)
		copy(windowData, s.data[i:i+size])
		result = append(result, &Slice[T]{windowData})
	}
	return result
}

// Chunk is an alias for Batch for convenience
func (s *Slice[T]) Chunk(chunkSize int) []*Slice[T] {
	return s.Batch(chunkSize)
}

// ToArray returns the underlying slice as a standard Go slice
// Always returns a copy to maintain encapsulation
func (s *Slice[T]) ToArray() []T {
	result := make([]T, len(s.data))
	copy(result, s.data)
	return result
}

// ToSlice is an alias for ToArray
func (s *Slice[T]) ToSlice() []T {
	return s.ToArray()
}

// Chain allows method chaining by returning the slice itself
func (s *Slice[T]) Chain() *Slice[T] {
	return s
}

// Pipe allows functional-style composition
// Applies the provided function and returns the result
func (s *Slice[T]) Pipe(fn func(*Slice[T]) *Slice[T]) *Slice[T] {
	return fn(s)
}

// Clone is an alias for Copy for clarity
func (s *Slice[T]) Clone() *Slice[T] {
	return s.Copy()
}

// First returns the first element if it exists
func (s *Slice[T]) First() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	return s.data[0], true
}

// Last returns the last element if it exists
func (s *Slice[T]) Last() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

// Take returns a new Slice containing the first n elements
func (s *Slice[T]) Take(n int) *Slice[T] {
	if n <= 0 {
		return &Slice[T]{data: []T{}}
	}
	if n >= len(s.data) {
		return s.Copy()
	}
	result := make([]T, n)
	copy(result, s.data[:n])
	return &Slice[T]{data: result}
}

// Skip returns a new Slice skipping the first n elements
func (s *Slice[T]) Skip(n int) *Slice[T] {
	if n <= 0 {
		return s.Copy()
	}
	if n >= len(s.data) {
		return &Slice[T]{data: []T{}}
	}
	result := make([]T, len(s.data)-n)
	copy(result, s.data[n:])
	return &Slice[T]{data: result}
}

// Slice returns a new Slice containing elements from index i to j (exclusive)
func (s *Slice[T]) Slice(i, j int) *Slice[T] {
	if i < 0 {
		i = 0
	}
	if j > len(s.data) {
		j = len(s.data)
	}
	if i >= j {
		return &Slice[T]{data: []T{}}
	}
	result := make([]T, j-i)
	copy(result, s.data[i:j])
	return &Slice[T]{data: result}
}
