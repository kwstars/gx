// Package randx provides cryptographically secure random utilities.
//
// The package contains helpers to generate secure random integers over generic
// integer types and helpers for weighted random selection (in other file).
// Comments follow GoDoc conventions: every exported identifier is documented.
package randx

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// Integer is a type constraint that permits any built-in integer type.
//
// This constraint can be used with generics to write functions that accept
// any signed or unsigned integer type.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// RandIntRange generates a cryptographically secure random integer of type T
// in the closed interval [min, max].
//
// The function returns an error if min > max or if the underlying random
// source fails. For unsigned integer types the implementation uses uint64
// arithmetic internally; for signed types it uses int64. If min == max the
// function returns min without consuming randomness.
func RandIntRange[T Integer](min, max T) (T, error) {
	if min > max {
		return 0, errors.New("min cannot be greater than max")
	}

	if min == max {
		return min, nil
	}

	var rangeSize *big.Int

	var isUnsigned bool
	var zeroVal T
	isUnsigned = (zeroVal - 1) > zeroVal

	if isUnsigned {
		uMin := uint64(min)
		uMax := uint64(max)
		rangeSizeU64 := uMax - uMin + 1
		rangeSize = new(big.Int).SetUint64(rangeSizeU64)
	} else {
		bigMin := big.NewInt(int64(min))
		bigMax := big.NewInt(int64(max))
		rangeSize = new(big.Int).Sub(bigMax, bigMin)
		rangeSize.Add(rangeSize, big.NewInt(1))
	}

	randomOffset, err := rand.Int(rand.Reader, rangeSize)
	if err != nil {
		return 0, fmt.Errorf("randx: failed to generate random number: %w", err)
	}

	var result T
	if isUnsigned {
		uMin := uint64(min)
		offset := randomOffset.Uint64()
		result = T(uMin + offset)
	} else {
		iMin := int64(min)
		offset := randomOffset.Int64()
		result = T(iMin + offset)
	}

	return result, nil
}
