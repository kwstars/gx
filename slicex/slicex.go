// Package slicex provides generic slice manipulation utilities.
// https://go.dev/wiki/SliceTricks
// All operations follow zero-allocation principles where possible.
package slicex

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

var (
	ErrIndexOutOfRange = errors.New("index out of range")
	ErrInvalidArgument = errors.New("invalid argument")
)

// Append appends elements to slice
func Append[T any](s []T, elems ...T) []T {
	return append(s, elems...)
}

// Copy creates a deep copy of slice
func Copy[T any](s []T) []T {
	if s == nil {
		return nil
	}
	b := make([]T, len(s))
	copy(b, s)
	return b
}

// Cut removes elements from index i to j
func Cut[T any](s []T, i, j int) ([]T, error) {
	if i < 0 || j > len(s) || i > j {
		return s, fmt.Errorf("%w: Cut(i=%d, j=%d) on slice of length %d", ErrIndexOutOfRange, i, j, len(s))
	}
	copy(s[i:], s[j:])
	var zero T
	for k, n := len(s)-j+i, len(s); k < n; k++ {
		s[k] = zero
	}
	return s[:len(s)-(j-i)], nil
}

// Delete removes element at index i
func Delete[T any](s []T, i int) ([]T, error) {
	if i < 0 || i >= len(s) {
		return s, fmt.Errorf("%w: Delete(i=%d) on slice of length %d", ErrIndexOutOfRange, i, len(s))
	}
	copy(s[i:], s[i+1:])
	var zero T
	s[len(s)-1] = zero
	return s[:len(s)-1], nil
}

// DeleteFast removes element at index i without preserving order
func DeleteFast[T any](s []T, i int) ([]T, error) {
	if i < 0 || i >= len(s) {
		return s, fmt.Errorf("%w: DeleteFast(i=%d) on slice of length %d", ErrIndexOutOfRange, i, len(s))
	}
	s[i] = s[len(s)-1]
	var zero T
	s[len(s)-1] = zero
	return s[:len(s)-1], nil
}

// Expand inserts n zero-value elements at position i
func Expand[T any](s []T, i, n int) ([]T, error) {
	if i < 0 || i > len(s) {
		return s, fmt.Errorf("%w: Expand(i=%d) on slice of length %d", ErrIndexOutOfRange, i, len(s))
	}
	if n < 0 {
		return s, fmt.Errorf("%w: Expand(n=%d) must be non-negative", ErrInvalidArgument, n)
	}
	if n == 0 {
		return s, nil
	}
	return append(s[:i], append(make([]T, n), s[i:]...)...), nil
}

// Extend appends n zero-value elements
func Extend[T any](s []T, n int) []T {
	return append(s, make([]T, n)...)
}

// ExtendCap ensures capacity for n additional elements
func ExtendCap[T any](s []T, n int) []T {
	if cap(s)-len(s) < n {
		return append(make([]T, 0, len(s)+n), s...)
	}
	return s
}

// Filter filters elements in place using predicate function
func Filter[T any](s []T, keep func(T) bool) []T {
	n := 0
	for _, x := range s {
		if keep(x) {
			s[n] = x
			n++
		}
	}
	return s[:n]
}

// Insert inserts element x at position i
func Insert[T any](s []T, i int, x T) ([]T, error) {
	if i < 0 || i > len(s) {
		return s, fmt.Errorf("%w: Insert(i=%d) on slice of length %d", ErrIndexOutOfRange, i, len(s))
	}
	s = append(s, *new(T))
	copy(s[i+1:], s[i:])
	s[i] = x
	return s, nil
}

// InsertSlice inserts slice vs at position i
func InsertSlice[T any](s []T, i int, vs ...T) ([]T, error) {
	if i < 0 || i > len(s) {
		return s, fmt.Errorf("%w: InsertSlice(i=%d) on slice of length %d", ErrIndexOutOfRange, i, len(s))
	}
	if len(vs) == 0 {
		return s, nil
	}

	if n := len(s) + len(vs); n <= cap(s) {
		s2 := s[:n]
		copy(s2[i+len(vs):], s[i:])
		copy(s2[i:], vs)
		return s2, nil
	}

	s2 := make([]T, len(s)+len(vs))
	copy(s2, s[:i])
	copy(s2[i:], vs)
	copy(s2[i+len(vs):], s[i:])
	return s2, nil
}

// Push appends element to end
func Push[T any](s []T, x T) []T {
	return append(s, x)
}

// Pop removes and returns last element
func Pop[T any](s []T) (elem T, result []T, err error) {
	if len(s) == 0 {
		return elem, s, fmt.Errorf("%w: Pop on empty slice", ErrIndexOutOfRange)
	}
	elem = s[len(s)-1]
	var zero T
	s[len(s)-1] = zero
	return elem, s[:len(s)-1], nil
}

// Unshift prepends element to front
func Unshift[T any](s []T, x T) []T {
	return append([]T{x}, s...)
}

// Shift removes and returns first element
func Shift[T any](s []T) (elem T, result []T, err error) {
	if len(s) == 0 {
		return elem, s, fmt.Errorf("%w: Shift on empty slice", ErrIndexOutOfRange)
	}
	return s[0], s[1:], nil
}

// Reverse reverses slice in place
func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Shuffle randomizes slice order using Fisher-Yates
func Shuffle[T any](s []T) {
	for i := len(s) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

// Batch splits slice into batches of given size
func Batch[T any](s []T, size int) [][]T {
	if size <= 0 || len(s) == 0 {
		return nil
	}
	batches := make([][]T, 0, (len(s)+size-1)/size)
	for size < len(s) {
		s, batches = s[size:], append(batches, s[:size:size])
	}
	return append(batches, s)
}

// Dedupe removes duplicates from sorted comparable slice in place
func Dedupe[T comparable](s []T) []T {
	if len(s) < 2 {
		return s
	}
	j := 0
	for i := 1; i < len(s); i++ {
		if s[j] == s[i] {
			continue
		}
		j++
		s[j] = s[i]
	}
	return s[:j+1]
}

// DedupeFunc removes duplicates using custom equality function
func DedupeFunc[T any](s []T, eq func(T, T) bool) []T {
	if len(s) < 2 {
		return s
	}
	j := 0
	for i := 1; i < len(s); i++ {
		if eq(s[j], s[i]) {
			continue
		}
		j++
		s[j] = s[i]
	}
	return s[:j+1]
}

// MoveToFront moves needle to front, prepends if not present
func MoveToFront[T comparable](needle T, haystack []T) []T {
	if len(haystack) != 0 && haystack[0] == needle {
		return haystack
	}
	prev := needle
	for i, elem := range haystack {
		switch {
		case i == 0:
			haystack[0] = needle
			prev = elem
		case elem == needle:
			haystack[i] = prev
			return haystack
		default:
			haystack[i] = prev
			prev = elem
		}
	}
	return append(haystack, prev)
}

// SlidingWindow creates sliding windows of given size
func SlidingWindow[T any](s []T, size int) [][]T {
	if size <= 0 || len(s) < size {
		return nil
	}
	r := make([][]T, 0, len(s)-size+1)
	for i, j := 0, size; j <= len(s); i, j = i+1, j+1 {
		r = append(r, s[i:j])
	}
	return r
}

// Contains checks if slice contains element
func Contains[T comparable](s []T, x T) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}

// Index returns first index of element, -1 if not found
func Index[T comparable](s []T, x T) int {
	for i, v := range s {
		if v == x {
			return i
		}
	}
	return -1
}

// Map applies function to each element
func Map[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

// Reduce reduces slice to single value
func Reduce[T, U any](s []T, init U, f func(U, T) U) U {
	acc := init
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}

// All checks if all elements satisfy predicate
func All[T any](s []T, f func(T) bool) bool {
	for _, v := range s {
		if !f(v) {
			return false
		}
	}
	return true
}

// Any checks if any element satisfies predicate
func Any[T any](s []T, f func(T) bool) bool {
	for _, v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

// Sort sorts comparable slice in place
func Sort[T sort.Interface](s T) {
	sort.Sort(s)
}

// Clear sets all elements to zero value
func Clear[T any](s []T) {
	var zero T
	for i := range s {
		s[i] = zero
	}
}
