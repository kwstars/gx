package slicex

import (
	"errors"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

// TestAppend verifies the Append function.
func TestAppend(t *testing.T) {
	tests := []struct {
		name  string
		s     []int
		elems []int
		want  []int
	}{
		{"nil slice", nil, []int{1}, []int{1}},
		{"empty slice", []int{}, []int{1, 2}, []int{1, 2}},
		{"append multiple", []int{1}, []int{2, 3}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Append(tt.s, tt.elems...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Append() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCut verifies Cut, including boundary checks and zeroing of underlying arrays.
func TestCut(t *testing.T) {
	tests := []struct {
		name    string
		s       []int
		i, j    int
		want    []int
		wantErr bool
	}{
		{"normal cut", []int{1, 2, 3, 4, 5}, 1, 3, []int{1, 4, 5}, false},
		{"cut all", []int{1, 2, 3}, 0, 3, []int{}, false},
		{"cut none", []int{1, 2, 3}, 1, 1, []int{1, 2, 3}, false},
		{"cut prefix", []int{1, 2, 3}, 0, 1, []int{2, 3}, false},
		{"cut suffix", []int{1, 2, 3}, 2, 3, []int{1, 2}, false},
		{"out of bounds i negative", []int{1}, -1, 0, []int{1}, true},
		{"out of bounds j too large", []int{1}, 0, 2, []int{1}, true},
		{"invalid range i > j", []int{1}, 1, 0, []int{1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Cut(tt.s, tt.i, tt.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cut() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cut() = %v, want %v", got, tt.want)
			}
			// Verify zeroing for memory leak prevention
			if !tt.wantErr && cap(got) > len(got) {
				full := got[:cap(got)]
				for k := len(got); k < cap(got); k++ {
					if full[k] != 0 {
						t.Errorf("Cut() memory leak: index %d not zeroed", k)
					}
				}
			}
		})
	}
}

// TestDelete verifies Delete and zeroing behavior.
func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		s       []int
		i       int
		want    []int
		wantErr bool
	}{
		{"delete middle", []int{1, 2, 3}, 1, []int{1, 3}, false},
		{"delete first", []int{1, 2, 3}, 0, []int{2, 3}, false},
		{"delete last", []int{1, 2, 3}, 2, []int{1, 2}, false},
		{"out of bounds negative", []int{1}, -1, []int{1}, true},
		{"out of bounds too large", []int{1}, 1, []int{1}, true},
		{"empty slice", []int{}, 0, []int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Delete(tt.s, tt.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
			// Verify zeroing
			if !tt.wantErr {
				full := got[:cap(got)]
				if full[len(got)] != 0 {
					t.Errorf("Delete() memory leak: tail not zeroed")
				}
			}
		})
	}
}

// TestDeleteFast verifies O(1) deletion and zeroing.
func TestDeleteFast(t *testing.T) {
	tests := []struct {
		name    string
		s       []int
		i       int
		want    []int // Note: order changes
		wantErr bool
	}{
		{"delete middle", []int{1, 2, 3, 4}, 1, []int{1, 4, 3}, false}, // 2 replaced by 4
		{"delete last", []int{1, 2, 3}, 2, []int{1, 2}, false},
		{"delete first", []int{1, 2, 3}, 0, []int{3, 2}, false},
		{"out of bounds", []int{1}, 1, []int{1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteFast(tt.s, tt.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFast() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteFast() = %v, want %v", got, tt.want)
			}
			// Verify zeroing
			if !tt.wantErr {
				full := got[:cap(got)]
				if full[len(got)] != 0 {
					t.Errorf("DeleteFast() memory leak: tail not zeroed")
				}
			}
		})
	}
}

// TestInsertSlice verifies InsertSlice logic.
func TestInsertSlice(t *testing.T) {
	tests := []struct {
		name    string
		s       []int
		i       int
		vs      []int
		want    []int
		wantErr bool
	}{
		{"insert middle", []int{1, 4}, 1, []int{2, 3}, []int{1, 2, 3, 4}, false},
		{"insert at start", []int{3}, 0, []int{1, 2}, []int{1, 2, 3}, false},
		{"insert at end", []int{1}, 1, []int{2}, []int{1, 2}, false},
		{"insert empty", []int{1, 2}, 1, []int{}, []int{1, 2}, false},
		{"insert into nil", nil, 0, []int{1}, []int{1}, false},
		{"out of bounds", []int{1}, 2, []int{2}, []int{1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InsertSlice(tt.s, tt.i, tt.vs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestExpand verifies Expand functionality.
func TestExpand(t *testing.T) {
	tests := []struct {
		name    string
		s       []int
		i, n    int
		want    []int
		wantErr bool
	}{
		{"expand middle", []int{1, 2}, 1, 2, []int{1, 0, 0, 2}, false},
		{"expand start", []int{1}, 0, 1, []int{0, 1}, false},
		{"expand end", []int{1}, 1, 1, []int{1, 0}, false},
		{"expand zero n", []int{1}, 0, 0, []int{1}, false},
		{"negative n", []int{1}, 0, -1, []int{1}, true},
		{"out of bounds i", []int{1}, 2, 1, []int{1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Expand(tt.s, tt.i, tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expand() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFilter verifies in-place filtering.
func TestFilter(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		keep func(int) bool
		want []int
	}{
		{"keep evens", []int{1, 2, 3, 4}, func(x int) bool { return x%2 == 0 }, []int{2, 4}},
		{"keep none", []int{1, 3}, func(x int) bool { return x%2 == 0 }, []int{}},
		{"keep all", []int{2, 4}, func(x int) bool { return x%2 == 0 }, []int{2, 4}},
		{"nil slice", nil, func(x int) bool { return true }, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.s, tt.keep)
			// For nil slice, got might be nil, reflect.DeepEqual handles nil vs empty slice specifically
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestPop verifies Pop operation and zeroing.
func TestPop(t *testing.T) {
	tests := []struct {
		name     string
		s        []int
		wantElem int
		wantRem  []int
		wantErr  bool
	}{
		{"normal pop", []int{1, 2, 3}, 3, []int{1, 2}, false},
		{"pop single", []int{1}, 1, []int{}, false},
		{"empty slice", []int{}, 0, []int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elem, result, err := Pop(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if elem != tt.wantElem {
					t.Errorf("Pop() elem = %v, want %v", elem, tt.wantElem)
				}
				if !reflect.DeepEqual(result, tt.wantRem) {
					t.Errorf("Pop() result = %v, want %v", result, tt.wantRem)
				}
				// Verify zeroing
				full := result[:cap(result)]
				if full[len(result)] != 0 {
					t.Errorf("Pop() memory leak: tail not zeroed")
				}
			}
		})
	}
}

// TestSlidingWindow verifies window creation logic.
func TestSlidingWindow(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		size int
		want [][]int
	}{
		{"normal window", []int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {2, 3}, {3, 4}}},
		{"window size equals len", []int{1, 2}, 2, [][]int{{1, 2}}},
		{"window size larger than len", []int{1, 2}, 3, nil},
		{"invalid size", []int{1, 2}, 0, nil},
		{"empty slice", []int{}, 1, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SlidingWindow(tt.s, tt.size)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SlidingWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDedupe verifies in-place deduplication.
func TestDedupe(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want []int
	}{
		{"duplicates sorted", []int{1, 1, 2, 3, 3, 3}, []int{1, 2, 3}},
		{"no duplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"all same", []int{1, 1, 1}, []int{1}},
		{"empty", []int{}, []int{}},
		{"single", []int{1}, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Dedupe(tt.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dedupe() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMoveToFront verifies element moving logic.
func TestMoveToFront(t *testing.T) {
	tests := []struct {
		name     string
		needle   int
		haystack []int
		want     []int
	}{
		{"move middle", 3, []int{1, 2, 3, 4}, []int{3, 1, 2, 4}},
		{"already front", 1, []int{1, 2, 3}, []int{1, 2, 3}},
		{"move last", 3, []int{1, 2, 3}, []int{3, 1, 2}},
		{"not present (prepend)", 0, []int{1, 2}, []int{0, 1, 2}},
		{"empty", 1, []int{}, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MoveToFront(tt.needle, tt.haystack)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MoveToFront() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBatch(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		size int
		want [][]int
	}{
		{"exact split", []int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {3, 4}}},
		{"uneven split", []int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}},
		{"size larger than slice", []int{1, 2}, 5, [][]int{{1, 2}}},
		{"invalid size", []int{1, 2}, 0, nil},
		{"empty slice", []int{}, 2, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Batch(tt.s, tt.size)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Batch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	s := []int{1, 3}
	got, err := Insert(s, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("Insert result %v", got)
	}
	_, err = Insert(s, 5, 9)
	if err == nil || !errors.Is(err, ErrIndexOutOfRange) {
		t.Errorf("expected ErrIndexOutOfRange")
	}
}

func TestExtend(t *testing.T) {
	s := []int{1, 2}
	got := Extend(s, 3)
	if len(got) != 5 || !reflect.DeepEqual(got[:2], []int{1, 2}) {
		t.Errorf("Extend failed %v", got)
	}
}

func TestExtendCap(t *testing.T) {
	s := make([]int, 2, 3)
	got := ExtendCap(s, 2) // need +2 capacity
	if cap(got) < len(s)+2 {
		t.Errorf("ExtendCap insufficient cap: %d", cap(got))
	}
	// When enough capacity no reallocation
	s2 := make([]int, 2, 10)
	got2 := ExtendCap(s2, 3)
	if &got2[0] != &s2[0] {
		t.Errorf("ExtendCap should not allocate when capacity sufficient")
	}
}

func TestCopy(t *testing.T) {
	s := []int{1, 2, 3}
	c := Copy(s)
	if !reflect.DeepEqual(s, c) {
		t.Errorf("Copy mismatch")
	}
	c[0] = 9
	if s[0] == 9 {
		t.Errorf("Copy not deep")
	}
	if Copy[int](nil) != nil {
		t.Errorf("Copy(nil) should be nil")
	}
}

func TestShiftUnshift(t *testing.T) {
	s := []int{2, 3}
	s = Unshift(s, 1)
	if !reflect.DeepEqual(s, []int{1, 2, 3}) {
		t.Errorf("Unshift failed %v", s)
	}
	elem, rest, err := Shift(s)
	if err != nil || elem != 1 || !reflect.DeepEqual(rest, []int{2, 3}) {
		t.Errorf("Shift failed elem=%v rest=%v err=%v", elem, rest, err)
	}
	_, _, err = Shift([]int{})
	if err == nil || !errors.Is(err, ErrIndexOutOfRange) {
		t.Errorf("Shift empty missing error")
	}
}

func TestReverse(t *testing.T) {
	s := []int{1, 2, 3, 4}
	Reverse(s)
	if !reflect.DeepEqual(s, []int{4, 3, 2, 1}) {
		t.Errorf("Reverse failed %v", s)
	}
}

func TestContainsIndex(t *testing.T) {
	s := []int{1, 2, 3}
	if !Contains(s, 2) || Contains(s, 9) {
		t.Errorf("Contains incorrect")
	}
	if Index(s, 3) != 2 || Index(s, 9) != -1 {
		t.Errorf("Index incorrect")
	}
}

func TestMapReduce(t *testing.T) {
	s := []int{1, 2, 3}
	m := Map(s, func(x int) int { return x * x })
	if !reflect.DeepEqual(m, []int{1, 4, 9}) {
		t.Errorf("Map failed %v", m)
	}
	sum := Reduce(s, 0, func(acc, x int) int { return acc + x })
	if sum != 6 {
		t.Errorf("Reduce expected 6 got %d", sum)
	}
}

func TestAllAny(t *testing.T) {
	s := []int{2, 4, 6}
	if !All(s, func(x int) bool { return x%2 == 0 }) {
		t.Errorf("All failed")
	}
	if Any(s, func(x int) bool { return x%2 == 1 }) {
		t.Errorf("Any failed")
	}
}

func TestErrorSentinels(t *testing.T) {
	_, err := Delete([]int{1}, 5)
	if err == nil || !errors.Is(err, ErrIndexOutOfRange) {
		t.Errorf("expected ErrIndexOutOfRange from Delete")
	}
	_, err = Expand([]int{1}, 0, -2)
	if err == nil || !errors.Is(err, ErrInvalidArgument) {
		t.Errorf("expected ErrInvalidArgument from Expand")
	}
}

func TestClear(t *testing.T) {
	s := []int{1, 2, 3}
	Clear(s)
	for i, v := range s {
		if v != 0 {
			t.Errorf("Clear() index %d = %v, want 0", i, v)
		}
	}
}

func TestPush(t *testing.T) {
	s := []int{1, 2}
	s2 := Push(s, 3)
	if !reflect.DeepEqual(s2, []int{1, 2, 3}) {
		t.Errorf("Push failed %v", s2)
	}
	var nilSlice []int
	nilSlice = Push(nilSlice, 1)
	if !reflect.DeepEqual(nilSlice, []int{1}) {
		t.Errorf("Push nil failed %v", nilSlice)
	}
}

func TestShuffle(t *testing.T) {
	rand.Seed(42)
	orig := []int{1, 2, 3, 4, 5, 6}
	s := make([]int, len(orig))
	copy(s, orig)
	Shuffle(s)
	if len(s) != len(orig) {
		t.Fatalf("Shuffle length changed")
	}
	count := map[int]int{}
	for _, v := range s {
		count[v]++
	}
	for _, v := range orig {
		if count[v] != 1 {
			t.Errorf("Shuffle element count broken for %d", v)
		}
	}
	// 可能极小概率保持原顺序，允许重试一次
	if reflect.DeepEqual(s, orig) {
		rand.Seed(43)
		copy(s, orig)
		Shuffle(s)
		if reflect.DeepEqual(s, orig) {
			t.Logf("Shuffle produced original order twice (acceptable but unlikely)")
		}
	}
}

func TestDedupeFunc(t *testing.T) {
	eq := func(a, b int) bool { return a == b }
	s := []int{1, 1, 2, 2, 2, 3}
	got := DedupeFunc(s, eq)
	if !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("DedupeFunc basic failed %v", got)
	}
	// 自定义等价：同奇偶视为相等（相邻同奇偶才折叠）
	eqParity := func(a, b int) bool { return a%2 == b%2 }
	s2 := []int{1, 3, 5, 2, 4, 6, 7}
	// 相邻奇数折叠 -> 1,2(偶组),7(最后奇组)
	got2 := DedupeFunc(s2, eqParity)
	if !reflect.DeepEqual(got2, []int{1, 2, 7}) {
		t.Errorf("DedupeFunc parity failed %v", got2)
	}
}

func TestSortGeneric(t *testing.T) {
	data := sort.IntSlice{5, 4, 2, 9, 1}
	Sort(data)
	if !reflect.DeepEqual([]int(data), []int{1, 2, 4, 5, 9}) {
		t.Errorf("Sort failed %v", data)
	}
}

func TestContainsIndexEmpty(t *testing.T) {
	var s []int
	if Contains(s, 1) {
		t.Errorf("Contains empty should be false")
	}
	if Index(s, 1) != -1 {
		t.Errorf("Index empty should be -1")
	}
}
