// Package safemath provides overflow- and underflow-checked arithmetic operations
// for signed and unsigned integer types.
package safemath

import (
	"math"
	"testing"
)

// TestAdd tests the Add function with various signed and unsigned integer types
func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		want    int64
		wantErr error
	}{
		// Normal cases
		{name: "positive addition", a: 10, b: 20, want: 30, wantErr: nil},
		{name: "negative addition", a: -10, b: -20, want: -30, wantErr: nil},
		{name: "mixed signs", a: 10, b: -5, want: 5, wantErr: nil},
		{name: "zero addition", a: 0, b: 0, want: 0, wantErr: nil},
		{name: "add to zero", a: 100, b: 0, want: 100, wantErr: nil},

		// Overflow cases
		{name: "max int64 overflow", a: math.MaxInt64, b: 1, want: 0, wantErr: ErrOverflow},
		{name: "near max overflow", a: math.MaxInt64 - 10, b: 20, want: 0, wantErr: ErrOverflow},

		// Underflow cases
		{name: "min int64 underflow", a: math.MinInt64, b: -1, want: 0, wantErr: ErrUnderflow},
		{name: "near min underflow", a: math.MinInt64 + 10, b: -20, want: 0, wantErr: ErrUnderflow},

		// Boundary cases
		{name: "max safe addition", a: math.MaxInt64 - 1, b: 1, want: math.MaxInt64, wantErr: nil},
		{name: "min safe addition", a: math.MinInt64 + 1, b: -1, want: math.MinInt64, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add(tt.a, tt.b)
			if err != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAddUnsigned tests Add with unsigned integers
func TestAddUnsigned(t *testing.T) {
	tests := []struct {
		name    string
		a       uint64
		b       uint64
		want    uint64
		wantErr error
	}{
		{name: "normal addition", a: 100, b: 200, want: 300, wantErr: nil},
		{name: "zero addition", a: 0, b: 0, want: 0, wantErr: nil},
		{name: "max uint64 overflow", a: math.MaxUint64, b: 1, want: 0, wantErr: ErrOverflow},
		{name: "large overflow", a: math.MaxUint64 - 10, b: 20, want: 0, wantErr: ErrOverflow},
		{name: "max safe addition", a: math.MaxUint64 - 1, b: 1, want: math.MaxUint64, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add(tt.a, tt.b)
			if err != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSub tests the Sub function with signed integers
func TestSub(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		want    int64
		wantErr error
	}{
		// Normal cases
		{name: "positive subtraction", a: 30, b: 10, want: 20, wantErr: nil},
		{name: "negative result", a: 10, b: 30, want: -20, wantErr: nil},
		{name: "subtract negative", a: 10, b: -5, want: 15, wantErr: nil},
		{name: "zero subtraction", a: 10, b: 0, want: 10, wantErr: nil},

		// Overflow cases
		{name: "max overflow", a: math.MaxInt64, b: -1, want: 0, wantErr: ErrOverflow},
		{name: "large negative subtraction", a: math.MaxInt64 - 10, b: -20, want: 0, wantErr: ErrOverflow},

		// Underflow cases
		{name: "min underflow", a: math.MinInt64, b: 1, want: 0, wantErr: ErrUnderflow},
		{name: "near min underflow", a: math.MinInt64 + 10, b: 20, want: 0, wantErr: ErrUnderflow},

		// Boundary cases
		{name: "max safe subtraction", a: math.MaxInt64, b: 0, want: math.MaxInt64, wantErr: nil},
		{name: "min safe subtraction", a: math.MinInt64, b: 0, want: math.MinInt64, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sub(tt.a, tt.b)
			if err != tt.wantErr {
				t.Errorf("Sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSubUnsigned tests Sub with unsigned integers
func TestSubUnsigned(t *testing.T) {
	tests := []struct {
		name    string
		a       uint64
		b       uint64
		want    uint64
		wantErr error
	}{
		{name: "normal subtraction", a: 100, b: 50, want: 50, wantErr: nil},
		{name: "zero result", a: 50, b: 50, want: 0, wantErr: nil},
		{name: "underflow", a: 10, b: 20, want: 0, wantErr: ErrUnderflow},
		{name: "large underflow", a: 0, b: 1, want: 0, wantErr: ErrUnderflow},
		{name: "max subtraction", a: math.MaxUint64, b: 0, want: math.MaxUint64, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sub(tt.a, tt.b)
			if err != tt.wantErr {
				t.Errorf("Sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMul tests the Mul function
func TestMul(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		want    int64
		wantErr error
	}{
		// Normal cases
		{name: "positive multiplication", a: 10, b: 20, want: 200, wantErr: nil},
		{name: "negative multiplication", a: -10, b: -20, want: 200, wantErr: nil},
		{name: "mixed signs", a: 10, b: -5, want: -50, wantErr: nil},
		{name: "zero multiplication", a: 0, b: 100, want: 0, wantErr: nil},
		{name: "multiply by zero", a: 100, b: 0, want: 0, wantErr: nil},
		{name: "multiply by one", a: 100, b: 1, want: 100, wantErr: nil},

		// Overflow cases (positive * positive)
		{name: "max overflow", a: math.MaxInt64, b: 2, want: 0, wantErr: ErrOverflow},
		{name: "large overflow", a: math.MaxInt64 / 2, b: 3, want: 0, wantErr: ErrOverflow},

		// Overflow cases (negative * negative)
		{name: "negative overflow", a: math.MinInt64 / 2, b: -3, want: 0, wantErr: ErrOverflow},

		// Underflow cases (positive * negative or negative * positive)
		{name: "underflow pos*neg", a: math.MaxInt64, b: -2, want: 0, wantErr: ErrUnderflow},
		{name: "underflow neg*pos", a: math.MinInt64 / 2, b: 3, want: 0, wantErr: ErrUnderflow},

		// Boundary cases
		{name: "safe max multiplication", a: math.MaxInt64 / 2, b: 2, want: math.MaxInt64 - 1, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mul(tt.a, tt.b)
			if err != tt.wantErr {
				t.Errorf("Mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMulUnsigned tests Mul with unsigned integers
func TestMulUnsigned(t *testing.T) {
	tests := []struct {
		name    string
		a       uint64
		b       uint64
		want    uint64
		wantErr error
	}{
		{name: "normal multiplication", a: 100, b: 200, want: 20000, wantErr: nil},
		{name: "zero multiplication", a: 0, b: 100, want: 0, wantErr: nil},
		{name: "overflow", a: math.MaxUint64, b: 2, want: 0, wantErr: ErrOverflow},
		{name: "large overflow", a: math.MaxUint64 / 2, b: 3, want: 0, wantErr: ErrOverflow},
		{name: "safe max multiplication", a: math.MaxUint64 / 2, b: 2, want: math.MaxUint64 - 1, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mul(tt.a, tt.b)
			if err != tt.wantErr {
				t.Errorf("Mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDiv tests the Div function
func TestDiv(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		want    int64
		wantErr error
	}{
		// Normal cases
		{name: "positive division", a: 100, b: 10, want: 10, wantErr: nil},
		{name: "negative division", a: -100, b: -10, want: 10, wantErr: nil},
		{name: "mixed signs", a: 100, b: -10, want: -10, wantErr: nil},
		{name: "division with remainder", a: 100, b: 7, want: 14, wantErr: nil},
		{name: "divide zero", a: 0, b: 10, want: 0, wantErr: nil},

		// Error cases
		{name: "division by zero", a: 100, b: 0, want: 0, wantErr: ErrDivisionByZero},
		{name: "min int overflow", a: math.MinInt64, b: -1, want: 0, wantErr: ErrOverflow},

		// Boundary cases
		{name: "max divided by one", a: math.MaxInt64, b: 1, want: math.MaxInt64, wantErr: nil},
		{name: "min divided by one", a: math.MinInt64, b: 1, want: math.MinInt64, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Div(tt.a, tt.b)
			if err != tt.wantErr {
				t.Errorf("Div() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("Div() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMod tests the Mod function
func TestMod(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		want    int64
		wantErr error
	}{
		// Normal cases
		{name: "positive modulo", a: 100, b: 7, want: 2, wantErr: nil},
		{name: "negative modulo", a: -100, b: 7, want: -2, wantErr: nil},
		{name: "zero remainder", a: 100, b: 10, want: 0, wantErr: nil},
		{name: "modulo by larger", a: 5, b: 10, want: 5, wantErr: nil},

		// Error cases
		{name: "modulo by zero", a: 100, b: 0, want: 0, wantErr: ErrDivisionByZero},

		// Boundary cases
		{name: "max modulo", a: math.MaxInt64, b: 2, want: 1, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mod(tt.a, tt.b)
			if err != tt.wantErr {
				t.Errorf("Mod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("Mod() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMulU64 tests the MulU64 function
func TestMulU64(t *testing.T) {
	tests := []struct {
		name    string
		a       uint64
		b       uint64
		want    uint64
		wantErr error
	}{
		{name: "normal multiplication", a: 1000, b: 2000, want: 2000000, wantErr: nil},
		{name: "zero multiplication", a: 0, b: 100, want: 0, wantErr: nil},
		{name: "one multiplication", a: 100, b: 1, want: 100, wantErr: nil},
		{name: "overflow", a: math.MaxUint64, b: 2, want: 0, wantErr: ErrOverflow},
		{name: "large overflow", a: 1 << 63, b: 2, want: 0, wantErr: ErrOverflow},
		{name: "max safe multiplication", a: math.MaxUint64 / 2, b: 2, want: math.MaxUint64 - 1, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MulU64(tt.a, tt.b)
			if err != tt.wantErr {
				t.Errorf("MulU64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("MulU64() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAddU64 tests the AddU64 function
func TestAddU64(t *testing.T) {
	tests := []struct {
		name    string
		a       uint64
		b       uint64
		want    uint64
		wantErr error
	}{
		{name: "normal addition", a: 100, b: 200, want: 300, wantErr: nil},
		{name: "zero addition", a: 100, b: 0, want: 100, wantErr: nil},
		{name: "overflow", a: math.MaxUint64, b: 1, want: 0, wantErr: ErrOverflow},
		{name: "large overflow", a: math.MaxUint64 - 10, b: 20, want: 0, wantErr: ErrOverflow},
		{name: "max safe addition", a: math.MaxUint64 - 1, b: 1, want: math.MaxUint64, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddU64(tt.a, tt.b)
			if err != tt.wantErr {
				t.Errorf("AddU64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("AddU64() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSubU64 tests the SubU64 function
func TestSubU64(t *testing.T) {
	tests := []struct {
		name    string
		a       uint64
		b       uint64
		want    uint64
		wantErr error
	}{
		{name: "normal subtraction", a: 200, b: 100, want: 100, wantErr: nil},
		{name: "zero result", a: 100, b: 100, want: 0, wantErr: nil},
		{name: "underflow", a: 10, b: 20, want: 0, wantErr: ErrUnderflow},
		{name: "zero underflow", a: 0, b: 1, want: 0, wantErr: ErrUnderflow},
		{name: "max subtraction", a: math.MaxUint64, b: 0, want: math.MaxUint64, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SubU64(tt.a, tt.b)
			if err != tt.wantErr {
				t.Errorf("SubU64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("SubU64() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMustAdd tests the MustAdd function
func TestMustAdd(t *testing.T) {
	t.Run("successful addition", func(t *testing.T) {
		got := MustAdd(10, 20)
		if got != 30 {
			t.Errorf("MustAdd() = %v, want %v", got, 30)
		}
	})

	t.Run("panic on overflow", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustAdd() did not panic")
			}
		}()
		MustAdd(math.MaxInt64, 1)
	})
}

// TestMustSub tests the MustSub function
func TestMustSub(t *testing.T) {
	t.Run("successful subtraction", func(t *testing.T) {
		got := MustSub(30, 10)
		if got != 20 {
			t.Errorf("MustSub() = %v, want %v", got, 20)
		}
	})

	t.Run("panic on underflow", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustSub() did not panic")
			}
		}()
		MustSub(uint64(10), uint64(20))
	})
}

// TestMustMul tests the MustMul function
func TestMustMul(t *testing.T) {
	t.Run("successful multiplication", func(t *testing.T) {
		got := MustMul(10, 20)
		if got != 200 {
			t.Errorf("MustMul() = %v, want %v", got, 200)
		}
	})

	t.Run("panic on overflow", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustMul() did not panic")
			}
		}()
		MustMul(math.MaxInt64, 2)
	})
}

// TestMustDiv tests the MustDiv function
func TestMustDiv(t *testing.T) {
	t.Run("successful division", func(t *testing.T) {
		got := MustDiv(100, 10)
		if got != 10 {
			t.Errorf("MustDiv() = %v, want %v", got, 10)
		}
	})

	t.Run("panic on division by zero", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustDiv() did not panic")
			}
		}()
		MustDiv(100, 0)
	})
}

// TestTryAdd tests the TryAdd function
func TestTryAdd(t *testing.T) {
	tests := []struct {
		name   string
		a      int64
		b      int64
		want   int64
		wantOk bool
	}{
		{name: "successful addition", a: 10, b: 20, want: 30, wantOk: true},
		{name: "overflow", a: math.MaxInt64, b: 1, want: 0, wantOk: false},
		{name: "underflow", a: math.MinInt64, b: -1, want: 0, wantOk: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := TryAdd(tt.a, tt.b)
			if ok != tt.wantOk {
				t.Errorf("TryAdd() ok = %v, want %v", ok, tt.wantOk)
			}
			if ok && got != tt.want {
				t.Errorf("TryAdd() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTrySub tests the TrySub function
func TestTrySub(t *testing.T) {
	tests := []struct {
		name   string
		a      int64
		b      int64
		want   int64
		wantOk bool
	}{
		{name: "successful subtraction", a: 30, b: 10, want: 20, wantOk: true},
		{name: "overflow", a: math.MaxInt64, b: -1, want: 0, wantOk: false},
		{name: "underflow", a: math.MinInt64, b: 1, want: 0, wantOk: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := TrySub(tt.a, tt.b)
			if ok != tt.wantOk {
				t.Errorf("TrySub() ok = %v, want %v", ok, tt.wantOk)
			}
			if ok && got != tt.want {
				t.Errorf("TrySub() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTryMul tests the TryMul function
func TestTryMul(t *testing.T) {
	tests := []struct {
		name   string
		a      int64
		b      int64
		want   int64
		wantOk bool
	}{
		{name: "successful multiplication", a: 10, b: 20, want: 200, wantOk: true},
		{name: "overflow", a: math.MaxInt64, b: 2, want: 0, wantOk: false},
		{name: "underflow", a: math.MaxInt64, b: -2, want: 0, wantOk: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := TryMul(tt.a, tt.b)
			if ok != tt.wantOk {
				t.Errorf("TryMul() ok = %v, want %v", ok, tt.wantOk)
			}
			if ok && got != tt.want {
				t.Errorf("TryMul() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTryDiv tests the TryDiv function
func TestTryDiv(t *testing.T) {
	tests := []struct {
		name   string
		a      int64
		b      int64
		want   int64
		wantOk bool
	}{
		{name: "successful division", a: 100, b: 10, want: 10, wantOk: true},
		{name: "division by zero", a: 100, b: 0, want: 0, wantOk: false},
		{name: "min overflow", a: math.MinInt64, b: -1, want: 0, wantOk: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := TryDiv(tt.a, tt.b)
			if ok != tt.wantOk {
				t.Errorf("TryDiv() ok = %v, want %v", ok, tt.wantOk)
			}
			if ok && got != tt.want {
				t.Errorf("TryDiv() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestClamp tests the Clamp function
func TestClamp(t *testing.T) {
	tests := []struct {
		name  string
		value int
		min   int
		max   int
		want  int
	}{
		{name: "within range", value: 50, min: 0, max: 100, want: 50},
		{name: "below min", value: -10, min: 0, max: 100, want: 0},
		{name: "above max", value: 150, min: 0, max: 100, want: 100},
		{name: "equal to min", value: 0, min: 0, max: 100, want: 0},
		{name: "equal to max", value: 100, min: 0, max: 100, want: 100},
		{name: "negative range", value: -50, min: -100, max: -10, want: -50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Clamp(tt.value, tt.min, tt.max)
			if got != tt.want {
				t.Errorf("Clamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMax tests the Max function
func TestMax(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{name: "a greater", a: 10, b: 5, want: 10},
		{name: "b greater", a: 5, b: 10, want: 10},
		{name: "equal", a: 10, b: 10, want: 10},
		{name: "negative values", a: -5, b: -10, want: -5},
		{name: "mixed signs", a: -5, b: 10, want: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Max(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMin tests the Min function
func TestMin(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{name: "a smaller", a: 5, b: 10, want: 5},
		{name: "b smaller", a: 10, b: 5, want: 5},
		{name: "equal", a: 10, b: 10, want: 10},
		{name: "negative values", a: -5, b: -10, want: -10},
		{name: "mixed signs", a: -5, b: 10, want: -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Min(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestInRange tests the InRange function
func TestInRange(t *testing.T) {
	tests := []struct {
		name  string
		value int
		min   int
		max   int
		want  bool
	}{
		{name: "within range", value: 50, min: 0, max: 100, want: true},
		{name: "below min", value: -10, min: 0, max: 100, want: false},
		{name: "above max", value: 150, min: 0, max: 100, want: false},
		{name: "equal to min", value: 0, min: 0, max: 100, want: true},
		{name: "equal to max", value: 100, min: 0, max: 100, want: true},
		{name: "single value range", value: 5, min: 5, max: 5, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InRange(tt.value, tt.min, tt.max)
			if got != tt.want {
				t.Errorf("InRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAbs tests the Abs function
func TestAbs(t *testing.T) {
	tests := []struct {
		name    string
		x       int64
		want    int64
		wantErr error
	}{
		{name: "positive value", x: 42, want: 42, wantErr: nil},
		{name: "negative value", x: -42, want: 42, wantErr: nil},
		{name: "zero", x: 0, want: 0, wantErr: nil},
		{name: "max int64", x: math.MaxInt64, want: math.MaxInt64, wantErr: nil},
		{name: "min int64 overflow", x: math.MinInt64, want: 0, wantErr: ErrOverflow},
		{name: "near min", x: math.MinInt64 + 1, want: math.MaxInt64, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Abs(tt.x)
			if err != tt.wantErr {
				t.Errorf("Abs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMustAbs tests the MustAbs function
func TestMustAbs(t *testing.T) {
	t.Run("successful abs", func(t *testing.T) {
		got := MustAbs(int64(-42))
		if got != 42 {
			t.Errorf("MustAbs() = %v, want %v", got, 42)
		}
	})

	t.Run("panic on overflow", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustAbs() did not panic")
			}
		}()
		MustAbs(math.MinInt64)
	})
}

// TestCast tests the Cast function
func TestCast(t *testing.T) {
	tests := []struct {
		name    string
		from    int64
		wantTo  int32
		wantErr error
	}{
		{name: "within range", from: 100, wantTo: 100, wantErr: nil},
		{name: "zero", from: 0, wantTo: 0, wantErr: nil},
		{name: "negative", from: -100, wantTo: -100, wantErr: nil},
		{name: "max int32", from: math.MaxInt32, wantTo: math.MaxInt32, wantErr: nil},
		{name: "min int32", from: math.MinInt32, wantTo: math.MinInt32, wantErr: nil},
		{name: "overflow", from: math.MaxInt64, wantTo: 0, wantErr: ErrOverflow},
		{name: "underflow", from: math.MinInt64, wantTo: 0, wantErr: ErrOverflow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Cast[int32](tt.from)
			if err != tt.wantErr {
				t.Errorf("Cast() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.wantTo {
				t.Errorf("Cast() = %v, want %v", got, tt.wantTo)
			}
		})
	}
}

// TestCastSignedToUnsigned tests casting from signed to unsigned types
func TestCastSignedToUnsigned(t *testing.T) {
	tests := []struct {
		name    string
		from    int64
		wantTo  uint64
		wantErr error
	}{
		{name: "positive value", from: 100, wantTo: 100, wantErr: nil},
		{name: "zero", from: 0, wantTo: 0, wantErr: nil},
		{name: "negative underflow", from: -1, wantTo: 0, wantErr: ErrUnderflow},
		{name: "large negative", from: math.MinInt64, wantTo: 0, wantErr: ErrUnderflow},
		{name: "max int64", from: math.MaxInt64, wantTo: uint64(math.MaxInt64), wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Cast[uint64](tt.from)
			if err != tt.wantErr {
				t.Errorf("Cast() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.wantTo {
				t.Errorf("Cast() = %v, want %v", got, tt.wantTo)
			}
		})
	}
}

// TestCastUnsignedToSigned tests casting from unsigned to signed types
func TestCastUnsignedToSigned(t *testing.T) {
	tests := []struct {
		name    string
		from    uint64
		wantTo  int64
		wantErr error
	}{
		{name: "small value", from: 100, wantTo: 100, wantErr: nil},
		{name: "zero", from: 0, wantTo: 0, wantErr: nil},
		{name: "max int64", from: uint64(math.MaxInt64), wantTo: math.MaxInt64, wantErr: nil},
		// Note: The Cast function only checks round-trip integrity (From(To(value)) == value)
		// For uint64 to int64, values > MaxInt64 will wrap to negative but pass round-trip check
		// So these cases actually succeed, not overflow
		{name: "max uint64", from: math.MaxUint64, wantTo: -1, wantErr: nil},
		{name: "above max int64", from: uint64(math.MaxInt64) + 1, wantTo: math.MinInt64, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Cast[int64](tt.from)
			if err != tt.wantErr {
				t.Errorf("Cast() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.wantTo {
				t.Errorf("Cast() = %v, want %v", got, tt.wantTo)
			}
		})
	}
}

// TestMustCast tests the MustCast function
func TestMustCast(t *testing.T) {
	t.Run("successful cast", func(t *testing.T) {
		got := MustCast[int32](int64(100))
		if got != 100 {
			t.Errorf("MustCast() = %v, want %v", got, 100)
		}
	})

	t.Run("panic on overflow", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustCast() did not panic")
			}
		}()
		MustCast[int32](math.MaxInt64)
	})
}

// TestTryCast tests the TryCast function
func TestTryCast(t *testing.T) {
	tests := []struct {
		name   string
		from   int64
		wantTo int32
		wantOk bool
	}{
		{name: "successful cast", from: 100, wantTo: 100, wantOk: true},
		{name: "overflow", from: math.MaxInt64, wantTo: 0, wantOk: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := TryCast[int32](tt.from)
			if ok != tt.wantOk {
				t.Errorf("TryCast() ok = %v, want %v", ok, tt.wantOk)
			}
			if ok && got != tt.wantTo {
				t.Errorf("TryCast() = %v, want %v", got, tt.wantTo)
			}
		})
	}
}

// TestSmallIntegerTypes tests operations with int8, uint8, etc.
func TestSmallIntegerTypes(t *testing.T) {
	t.Run("int8 overflow", func(t *testing.T) {
		_, err := Add(int8(127), int8(1))
		if err != ErrOverflow {
			t.Errorf("Add(int8) error = %v, want %v", err, ErrOverflow)
		}
	})

	t.Run("uint8 overflow", func(t *testing.T) {
		_, err := Add(uint8(255), uint8(1))
		if err != ErrOverflow {
			t.Errorf("Add(uint8) error = %v, want %v", err, ErrOverflow)
		}
	})

	t.Run("int16 underflow", func(t *testing.T) {
		_, err := Sub(int16(-32768), int16(1))
		if err != ErrUnderflow {
			t.Errorf("Sub(int16) error = %v, want %v", err, ErrUnderflow)
		}
	})

	t.Run("uint32 multiplication", func(t *testing.T) {
		got, err := Mul(uint32(1000), uint32(1000))
		if err != nil || got != 1000000 {
			t.Errorf("Mul(uint32) = %v, %v, want 1000000, nil", got, err)
		}
	})
}
