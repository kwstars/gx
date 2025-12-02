// Package randx provides utilities for cryptographically secure random
// operations. It includes helpers for generating random integers of generic
// integer types and utilities for weighted random selection.
package randx

import (
	"crypto/rand"
	"math/big"
	"sort"
)

// Picker is a generic weighted random picker.
//
// Picker holds the original items, a slice of prefix sums for fast selection,
// and the totalWeight of all items computed at construction time. Weights are
// integer values provided by the caller via weightFunc.
type Picker[T any] struct {
	items       []T
	prefixSums  []int
	totalWeight int
}

// New constructs a Picker for the provided items using weightFunc to obtain
// a non-negative integer weight for each item.
//
// At least one item must have a positive weight for selection to succeed;
// otherwise Pick will return an error. New does not validate negative weights
// — callers should ensure weightFunc returns non-negative values.
func New[T any](items []T, weightFunc func(T) int) *Picker[T] {
	p := &Picker[T]{
		items:      items,
		prefixSums: make([]int, len(items)),
	}

	sum := 0
	for i, item := range items {
		weight := weightFunc(item)
		sum += weight
		p.prefixSums[i] = sum
	}
	p.totalWeight = sum

	return p
}

// Pick returns a randomly selected item according to the configured weights.
//
// The selection is proportional to each item's weight. If the Picker contains
// no items Pick returns ErrEmptyPicker. If the sum of weights is zero, the
// behavior will result in an error from the random source; callers should
// ensure at least one positive weight exists.
func (p *Picker[T]) Pick() (T, error) {
	if len(p.items) == 0 {
		var zero T
		return zero, &ErrEmptyPicker{}
	}

	// Generate a random number between 1 and totalWeight
	n, err := rand.Int(rand.Reader, big.NewInt(int64(p.totalWeight)))
	if err != nil {
		var zero T
		return zero, err
	}
	x := int(n.Int64()) + 1

	index := sort.SearchInts(p.prefixSums, x)
	return p.items[index], nil
}

// ErrEmptyPicker is returned when attempting to pick from an empty Picker.
type ErrEmptyPicker struct{}

func (e *ErrEmptyPicker) Error() string {
	return "picker is empty"
}
