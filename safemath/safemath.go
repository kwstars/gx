// Package safemath provides overflow- and underflow-checked arithmetic operations
// for signed and unsigned integer types.
package safemath

import (
	"cmp"
	"errors"
	"math"
	"math/bits"
)

var (
	// ErrOverflow is returned when an operation would overflow the representable range.
	ErrOverflow = errors.New("safemath: operation would overflow")

	// ErrUnderflow is returned when an operation would underflow the representable range,
	// e.g., subtracting a larger unsigned value from a smaller one.
	ErrUnderflow = errors.New("safemath: operation would underflow")

	// ErrDivisionByZero is returned when division or modulo by zero is attempted.
	ErrDivisionByZero = errors.New("safemath: division by zero")
)

// Signed is a type constraint for all signed integer types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a type constraint for all unsigned integer types.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a type constraint for all integer types,
// i.e., both signed and unsigned.
type Integer interface {
	Signed | Unsigned
}

// Add returns a + b if no overflow or underflow occurs.
func Add[T Integer](a, b T) (T, error) {
	var zero T
	result := a + b

	switch any(a).(type) {
	case int, int8, int16, int32, int64:
		// Signed: overflow if b > 0 and a > max - b; underflow if b < 0 and a < min - b.
		if b > 0 && a > maxValue[T]()-b {
			return zero, ErrOverflow
		}
		if b < 0 && a < minValue[T]()-b {
			return zero, ErrUnderflow
		}
	default:
		// Unsigned: overflow if result < a (since b ≥ 0).
		if result < a {
			return zero, ErrOverflow
		}
	}

	return result, nil
}

// Sub returns a - b if no overflow or underflow occurs.
func Sub[T Integer](a, b T) (T, error) {
	var zero T
	result := a - b

	switch any(a).(type) {
	case int, int8, int16, int32, int64:
		// Signed: overflow if b < 0 and a > max + b; underflow if b > 0 and a < min + b.
		if b < 0 && a > maxValue[T]()+b {
			return zero, ErrOverflow
		}
		if b > 0 && a < minValue[T]()+b {
			return zero, ErrUnderflow
		}
	default:
		// Unsigned: underflow if result > a (due to borrow).
		if result > a {
			return zero, ErrUnderflow
		}
	}

	return result, nil
}

// Mul returns a * b if no overflow or underflow occurs.
func Mul[T Integer](a, b T) (T, error) {
	var zero T
	if a == 0 || b == 0 {
		return 0, nil
	}
	result := a * b

	// Check for overflow: if division recovers original operand, no overflow.
	if result/b != a {
		switch any(a).(type) {
		case int, int8, int16, int32, int64:
			// Signed: same-sign → overflow; opposite-sign → underflow.
			if (a > 0 && b > 0) || (a < 0 && b < 0) {
				return zero, ErrOverflow
			}
			return zero, ErrUnderflow
		default:
			// Unsigned only overflows.
			return zero, ErrOverflow
		}
	}

	return result, nil
}

// Div returns a / b if b is nonzero and no overflow occurs.
// Notably, signed integer division of math.MinInt by -1 overflows.
func Div[T Integer](a, b T) (T, error) {
	var zero T
	if b == 0 {
		return zero, ErrDivisionByZero
	}

	switch any(a).(type) {
	case int, int8, int16, int32, int64:
		// MinInt / -1 overflows (e.g., math.MinInt64 / -1).
		// Construct -1 for signed types to avoid compile-time overflow with unsigned types.
		negativeOne := zero - 1
		if a == minValue[T]() && b == negativeOne {
			return zero, ErrOverflow
		}
	}

	return a / b, nil
}

// Mod returns a % b if b is nonzero.
func Mod[T Integer](a, b T) (T, error) {
	var zero T
	if b == 0 {
		return zero, ErrDivisionByZero
	}
	return a % b, nil
}

// MulU64 uses bits.Mul64 to perform overflow-checked multiplication of uint64 values.
func MulU64(a, b uint64) (uint64, error) {
	hi, lo := bits.Mul64(a, b)
	if hi != 0 {
		return 0, ErrOverflow
	}
	return lo, nil
}

// AddU64 uses bits.Add64 to perform overflow-checked addition of uint64 values.
func AddU64(a, b uint64) (uint64, error) {
	sum, carry := bits.Add64(a, b, 0)
	if carry != 0 {
		return 0, ErrOverflow
	}
	return sum, nil
}

// SubU64 uses bits.Sub64 to perform underflow-checked subtraction of uint64 values.
func SubU64(a, b uint64) (uint64, error) {
	diff, borrow := bits.Sub64(a, b, 0)
	if borrow != 0 {
		return 0, ErrUnderflow
	}
	return diff, nil
}

// MustAdd returns a + b, panicking on overflow or underflow.
func MustAdd[T Integer](a, b T) T {
	result, err := Add(a, b)
	if err != nil {
		panic(err)
	}
	return result
}

// MustSub returns a - b, panicking on overflow or underflow.
func MustSub[T Integer](a, b T) T {
	result, err := Sub(a, b)
	if err != nil {
		panic(err)
	}
	return result
}

// MustMul returns a * b, panicking on overflow or underflow.
func MustMul[T Integer](a, b T) T {
	result, err := Mul(a, b)
	if err != nil {
		panic(err)
	}
	return result
}

// MustDiv returns a / b, panicking on division by zero or overflow.
func MustDiv[T Integer](a, b T) T {
	result, err := Div(a, b)
	if err != nil {
		panic(err)
	}
	return result
}

// TryAdd returns a + b and true if no overflow or underflow occurs; otherwise false.
func TryAdd[T Integer](a, b T) (T, bool) {
	result, err := Add(a, b)
	return result, err == nil
}

// TrySub returns a - b and true if no overflow or underflow occurs; otherwise false.
func TrySub[T Integer](a, b T) (T, bool) {
	result, err := Sub(a, b)
	return result, err == nil
}

// TryMul returns a * b and true if no overflow or underflow occurs; otherwise false.
func TryMul[T Integer](a, b T) (T, bool) {
	result, err := Mul(a, b)
	return result, err == nil
}

// TryDiv returns a / b and true if b is nonzero and no overflow occurs; otherwise false.
func TryDiv[T Integer](a, b T) (T, bool) {
	result, err := Div(a, b)
	return result, err == nil
}

// Clamp restricts value to the range [min, max].
// It is equivalent to Min(Max(value, min), max).
func Clamp[T cmp.Ordered](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Max returns the greater of a or b.
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min returns the lesser of a or b.
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// InRange reports whether value is in the inclusive range [min, max].
func InRange[T cmp.Ordered](value, min, max T) bool {
	return value >= min && value <= max
}

// Abs returns the absolute value of x.
// For signed integers, it returns an error if x is the minimum representable value
// (e.g., math.MinInt64), as its negation would overflow.
func Abs[T Signed](x T) (T, error) {
	if x < 0 {
		if x == minValue[T]() {
			var zero T
			return zero, ErrOverflow
		}
		return -x, nil
	}
	return x, nil
}

// MustAbs returns the absolute value of x, panicking on overflow.
func MustAbs[T Signed](x T) T {
	result, err := Abs(x)
	if err != nil {
		panic(err)
	}
	return result
}

// Cast safely converts value from type From to type To.
// It returns an error if the conversion would lose precision or,
// when converting signed to unsigned, if the value is negative.
func Cast[To Integer, From Integer](value From) (To, error) {
	var zero To
	result := To(value)
	// Check round-trip integrity.
	if From(result) != value {
		return zero, ErrOverflow
	}
	// Disallow negative → unsigned conversion.
	switch any(value).(type) {
	case int, int8, int16, int32, int64:
		if value < 0 {
			switch any(zero).(type) {
			case uint, uint8, uint16, uint32, uint64, uintptr:
				return zero, ErrUnderflow
			}
		}
	}
	return result, nil
}

// MustCast converts value from From to To, panicking on error.
func MustCast[To Integer, From Integer](value From) To {
	result, err := Cast[To](value)
	if err != nil {
		panic(err)
	}
	return result
}

// TryCast converts value from From to To and reports success.
func TryCast[To Integer, From Integer](value From) (To, bool) {
	result, err := Cast[To](value)
	return result, err == nil
}

// maxValue returns the maximum value representable by type T.
func maxValue[T Integer]() T {
	var v T
	switch any(v).(type) {
	case int:
		return any(math.MaxInt).(T)
	case int8:
		return any(int8(math.MaxInt8)).(T)
	case int16:
		return any(int16(math.MaxInt16)).(T)
	case int32:
		return any(int32(math.MaxInt32)).(T)
	case int64:
		return any(int64(math.MaxInt64)).(T)
	case uint:
		return any(uint(math.MaxUint)).(T)
	case uint8:
		return any(uint8(math.MaxUint8)).(T)
	case uint16:
		return any(uint16(math.MaxUint16)).(T)
	case uint32:
		return any(uint32(math.MaxUint32)).(T)
	case uint64:
		return any(uint64(math.MaxUint64)).(T)
	case uintptr:
		return any(^uintptr(0)).(T)
	}
	return v
}

// minValue returns the minimum value representable by type T.
// For unsigned types, it returns zero.
func minValue[T Integer]() T {
	var v T
	switch any(v).(type) {
	case int:
		return any(math.MinInt).(T)
	case int8:
		return any(int8(math.MinInt8)).(T)
	case int16:
		return any(int16(math.MinInt16)).(T)
	case int32:
		return any(int32(math.MinInt32)).(T)
	case int64:
		return any(int64(math.MinInt64)).(T)
	default:
		return 0
	}
}
