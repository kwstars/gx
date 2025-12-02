package randx

import (
	"testing"
)

// TestPicker tests the weighted random picker.
func TestPicker(t *testing.T) {
	t.Parallel()
	type item struct {
		value  int
		weight int
	}

	items := []item{
		{value: 1, weight: 1},
		{value: 2, weight: 2},
		{value: 3, weight: 3},
		{value: 4, weight: 4},
	}

	// Create a new weighted random picker
	picker := New(items, func(i item) int { return i.weight })

	counts := make(map[int]int)
	const iterations = 1000000 // Number of iterations

	// Perform multiple picks and count how often each value is chosen
	for i := 0; i < iterations; i++ {
		picked, err := picker.Pick()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		counts[picked.value]++
	}

	totalWeight := 0
	// Calculate total weight
	for _, item := range items {
		totalWeight += item.weight
	}

	// Verify that each value's actual frequency is close to expected frequency
	for _, item := range items {
		expected := float64(item.weight) / float64(totalWeight)
		actual := float64(counts[item.value]) / float64(iterations)

		if diff := abs(expected - actual); diff > 0.01 {
			t.Errorf("Value %d: expected frequency %.4f, got %.4f", item.value, expected, actual)
		}
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
