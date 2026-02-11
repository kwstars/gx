package slices

import (
	"math/rand"
	"testing"
)

// TestCutWithCleanup verifies the bug fix in CutWithCleanup
func TestCutWithCleanup(t *testing.T) {
	s := NewSlice([]int{1, 2, 3, 4, 5, 6, 7, 8})
	s.CutWithCleanup(2, 5) // Remove elements at indices 2, 3, 4

	expected := []int{1, 2, 6, 7, 8}
	result := s.ToArray()

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], result[i])
		}
	}
}

// TestFilter verifies GC cleanup in Filter
func TestFilter(t *testing.T) {
	s := NewSlice([]int{1, 2, 3, 4, 5, 6, 7, 8})
	s.Filter(func(x int) bool { return x%2 == 0 })

	expected := []int{2, 4, 6, 8}
	result := s.ToArray()

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], result[i])
		}
	}
}

// TestPop verifies GC cleanup in Pop
func TestPop(t *testing.T) {
	s := NewSlice([]int{1, 2, 3, 4, 5})

	val, ok := s.Pop()
	if !ok || val != 5 {
		t.Errorf("Expected 5 and true, got %d and %v", val, ok)
	}

	if s.Len() != 4 {
		t.Errorf("Expected length 4, got %d", s.Len())
	}
}

// TestDeduplicate verifies the bug fix in Deduplicate
func TestDeduplicate(t *testing.T) {
	s := NewSlice([]int{3, 1, 2, 1, 3, 2, 4})
	s.Deduplicate(func(a, b int) int { return a - b })

	expected := []int{1, 2, 3, 4}
	result := s.ToArray()

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], result[i])
		}
	}
}

// TestBatch verifies edge cases
func TestBatch(t *testing.T) {
	// Test normal case
	s := NewSlice([]int{1, 2, 3, 4, 5, 6, 7})
	batches := s.Batch(3)

	if len(batches) != 3 {
		t.Errorf("Expected 3 batches, got %d", len(batches))
	}

	// Test empty slice
	empty := NewSlice([]int{})
	emptyBatches := empty.Batch(3)
	if len(emptyBatches) != 0 {
		t.Errorf("Expected 0 batches for empty slice, got %d", len(emptyBatches))
	}

	// Test invalid batch size
	invalidBatches := s.Batch(0)
	if len(invalidBatches) != 1 {
		t.Errorf("Expected 1 batch for invalid size, got %d", len(invalidBatches))
	}
}

// TestSlidingWindow verifies edge cases
func TestSlidingWindow(t *testing.T) {
	s := NewSlice([]int{1, 2, 3, 4, 5})

	// Normal case
	windows := s.SlidingWindow(3)
	if len(windows) != 3 {
		t.Errorf("Expected 3 windows, got %d", len(windows))
	}

	// Window size larger than slice
	largeWindows := s.SlidingWindow(10)
	if len(largeWindows) != 1 {
		t.Errorf("Expected 1 window when size > length, got %d", len(largeWindows))
	}

	// Invalid window size
	invalidWindows := s.SlidingWindow(0)
	if len(invalidWindows) != 0 {
		t.Errorf("Expected 0 windows for size 0, got %d", len(invalidWindows))
	}
}

// TestMethodChaining demonstrates method chaining
func TestMethodChaining(t *testing.T) {
	s := NewSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	result := s.
		Filter(func(x int) bool { return x%2 == 0 }). // Keep even numbers
		Map(func(x int) int { return x * 2 }).        // Double them
		ToArray()

	expected := []int{4, 8, 12, 16, 20}

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], result[i])
		}
	}
}

// TestInsertNoAlloc verifies the fix using zero value
func TestInsertNoAlloc(t *testing.T) {
	s := NewSlice([]int{1, 2, 4, 5})
	s.InsertNoAlloc(2, 3)

	expected := []int{1, 2, 3, 4, 5}
	result := s.ToArray()

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], result[i])
		}
	}
}

// TestAppendSlice verifies nil safety
func TestAppendSlice(t *testing.T) {
	s := NewSlice([]int{1, 2, 3})

	// Should handle nil gracefully
	s.AppendSlice(nil)
	if s.Len() != 3 {
		t.Errorf("Expected length 3 after appending nil, got %d", s.Len())
	}

	// Normal append
	other := NewSlice([]int{4, 5})
	s.AppendSlice(other)
	if s.Len() != 5 {
		t.Errorf("Expected length 5, got %d", s.Len())
	}
}

// TestReduce demonstrates usage
func TestReduce(t *testing.T) {
	s := NewSlice([]int{1, 2, 3, 4, 5})
	sum := s.Reduce(0, func(acc, val int) int { return acc + val })

	if sum != 15 {
		t.Errorf("Expected sum 15, got %d", sum)
	}
}

// TestFind demonstrates usage
func TestFind(t *testing.T) {
	s := NewSlice([]int{1, 2, 3, 4, 5})

	val, ok := s.Find(func(x int) bool { return x > 3 })
	if !ok || val != 4 {
		t.Errorf("Expected 4 and true, got %d and %v", val, ok)
	}

	val2, ok2 := s.Find(func(x int) bool { return x > 10 })
	if ok2 {
		t.Errorf("Expected false for not found, got true with value %d", val2)
	}
}

// TestShuffleWithRand verifies reproducible shuffling
func TestShuffleWithRand(t *testing.T) {
	s1 := NewSlice([]int{1, 2, 3, 4, 5})
	s2 := NewSlice([]int{1, 2, 3, 4, 5})

	// Same seed should produce same shuffle
	r1 := rand.New(rand.NewSource(42))
	r2 := rand.New(rand.NewSource(42))

	s1.ShuffleWithRand(r1)
	s2.ShuffleWithRand(r2)

	result1 := s1.ToArray()
	result2 := s2.ToArray()

	for i := range result1 {
		if result1[i] != result2[i] {
			t.Errorf("Shuffles with same seed should match, but differ at index %d", i)
		}
	}
}

// TestAll and TestAny
func TestAllAny(t *testing.T) {
	s := NewSlice([]int{2, 4, 6, 8})

	if !s.All(func(x int) bool { return x%2 == 0 }) {
		t.Error("Expected all elements to be even")
	}

	if s.Any(func(x int) bool { return x%2 != 0 }) {
		t.Error("Expected no odd elements")
	}
}

// TestTakeSkip verifies new utility methods
func TestTakeSkip(t *testing.T) {
	s := NewSlice([]int{1, 2, 3, 4, 5})

	taken := s.Take(3).ToArray()
	expectedTaken := []int{1, 2, 3}
	for i := range expectedTaken {
		if taken[i] != expectedTaken[i] {
			t.Errorf("Take: at index %d expected %d, got %d", i, expectedTaken[i], taken[i])
		}
	}

	skipped := s.Skip(2).ToArray()
	expectedSkipped := []int{3, 4, 5}
	for i := range expectedSkipped {
		if skipped[i] != expectedSkipped[i] {
			t.Errorf("Skip: at index %d expected %d, got %d", i, expectedSkipped[i], skipped[i])
		}
	}
}

// TestDeduplicateStable verifies order-preserving deduplication
func TestDeduplicateStable(t *testing.T) {
	s := NewSlice([]int{3, 1, 2, 1, 3, 2, 4})
	s.DeduplicateStable(func(a, b int) bool { return a == b })

	expected := []int{3, 1, 2, 4}
	result := s.ToArray()

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], result[i])
		}
	}
}

// BenchmarkFilter benchmarks the Filter operation
func BenchmarkFilter(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := NewSlice(data)
		s.Filter(func(x int) bool { return x%2 == 0 })
	}
}

// BenchmarkInsert benchmarks Insert vs InsertNoAlloc
func BenchmarkInsert(b *testing.B) {
	data := make([]int, 100)
	for i := range data {
		data[i] = i
	}

	b.Run("Insert", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := NewSlice(data)
			s.Insert(50, 999)
		}
	})

	b.Run("InsertNoAlloc", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := NewSlice(data)
			s.InsertNoAlloc(50, 999)
		}
	})
}
