package orderedmap

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEmpty(t *testing.T) {
	m := New[int, string]()
	checkAll(t, m, []Item[int, string]{})

	m.Clear()
	checkAll(t, m, []Item[int, string]{})
}

func TestPointers(t *testing.T) {
	type s struct {
		A string
		B int
	}

	m := New[int, *s]()
	checkAll(t, m, []Item[int, *s]{})

	want := &s{A: "one", B: 1}
	m.PushFront(1, want)

	{
		got, ok := m.Get(1)
		if !ok {
			t.Fatalf("key not found")
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("unexpected value (-want +got):\n%s", diff)
		}
	}

	// modify value, expected value stored by the map to change too
	{
		want.A = "two"

		got, ok := m.Get(1)
		if !ok {
			t.Fatalf("key not found")
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("unexpected value (-want +got):\n%s", diff)
		}
	}
}

func TestPushBack(t *testing.T) {
	m := New[int, string]()
	if err := m.PushBack(1, "one"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{1, "one"}})
	if err := m.PushBack(2, "two"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{1, "one"}, {2, "two"}})
	m.Clear()
	checkAll(t, m, []Item[int, string]{})
}

func TestPushFront(t *testing.T) {
	m := New[int, string]()
	if err := m.PushFront(1, "one"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{1, "one"}})
	if err := m.PushFront(2, "two"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{2, "two"}, {1, "one"}})
	m.Clear()
	checkAll(t, m, []Item[int, string]{})
}

func TestInsertAfter(t *testing.T) {
	m := New[int, string]()
	if err := m.InsertAfter(1, "one", 0); err == nil {
		t.Fatal("expected error, none got")
	}
	if err := m.PushFront(0, "zero"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := m.InsertAfter(2, "two", 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{0, "zero"}, {2, "two"}})
	if err := m.InsertAfter(1, "one", 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{0, "zero"}, {1, "one"}, {2, "two"}})
	if err := m.InsertAfter(3, "three", 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{0, "zero"}, {1, "one"}, {2, "two"}, {3, "three"}})
}

func TestInsertBefore(t *testing.T) {
	m := New[int, string]()
	if err := m.InsertBefore(1, "one", 0); err != ErrMarkKeyMissing {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := m.PushFront(0, "zero"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := m.InsertBefore(2, "two", 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{2, "two"}, {0, "zero"}})
	if err := m.InsertBefore(1, "one", 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{2, "two"}, {1, "one"}, {0, "zero"}})
	if err := m.InsertBefore(3, "three", 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{3, "three"}, {2, "two"}, {1, "one"}, {0, "zero"}})
}

func TestMoveToFront(t *testing.T) {
	m := New[int, string]()

	if err := m.MoveToFront(1); err != ErrKeyMissing {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := m.PushFront(1, "one"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := m.MoveToFront(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{1, "one"}})

	if err := m.PushFront(2, "two"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{2, "two"}, {1, "one"}})

	if err := m.MoveToFront(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{1, "one"}, {2, "two"}})

	if err := m.PushFront(0, "zero"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{0, "zero"}, {1, "one"}, {2, "two"}})

	if err := m.MoveToFront(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{1, "one"}, {0, "zero"}, {2, "two"}})

	if err := m.MoveToFront(2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{2, "two"}, {1, "one"}, {0, "zero"}})

	if err := m.MoveToFront(3); err != ErrKeyMissing {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMoveToBaack(t *testing.T) {
	m := New[int, string]()

	if err := m.MoveToBack(1); err != ErrKeyMissing {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := m.PushBack(1, "one"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := m.MoveToBack(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{1, "one"}})

	if err := m.PushBack(2, "two"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{1, "one"}, {2, "two"}})

	if err := m.MoveToBack(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{2, "two"}, {1, "one"}})

	if err := m.PushBack(0, "zero"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{2, "two"}, {1, "one"}, {0, "zero"}})

	if err := m.MoveToBack(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{2, "two"}, {0, "zero"}, {1, "one"}})

	if err := m.MoveToBack(2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkAll(t, m, []Item[int, string]{{0, "zero"}, {1, "one"}, {2, "two"}})

	if err := m.MoveToBack(3); err != ErrKeyMissing {
		t.Fatalf("unexpected error: %v", err)
	}
}

func checkAll[K comparable, V any](t *testing.T, om *OrderedMap[K, V], items []Item[K, V]) {
	t.Helper()

	if want, got := len(items), om.Len(); want != got {
		t.Fatalf("incorrect length: want: %d, got: %d", want, got)
	}
	if diff := cmp.Diff(items, om.Items()); diff != "" {
		t.Fatalf("unexpected keys (-want +got):\n%s", diff)
	}

	checkFrontBack(t, om, items)
	checkMapGet(t, om, items)
	checkKeys(t, om, items)
	checkPrevNext(t, om, items)
}

// checkMapGet converts items to map and validate all entries are present and return the correct value
func checkMapGet[K comparable, V any](t *testing.T, om *OrderedMap[K, V], items []Item[K, V]) {
	t.Helper()

	m := make(map[K]V, len(items))
	for _, item := range items {
		m[item.Key] = item.Value
	}
	if diff := cmp.Diff(m, om.Map()); diff != "" {
		t.Fatalf("unexpected map (-want +got):\n%s", diff)
	}
	for k, want := range m {
		got, ok := om.Get(k)
		if !ok {
			t.Fatalf("key %v not found", k)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("unexpected value for key %v (-want +got):\n%s", k, diff)
		}
	}
}

// checkFrontBack checks front and back correctness
func checkFrontBack[K comparable, V any](t *testing.T, om *OrderedMap[K, V], items []Item[K, V]) {
	t.Helper()

	if len(items) == 0 {
		if front, ok := om.Front(); ok {
			t.Fatalf("unexpected front: %v", front)
		}
		if back, ok := om.Back(); ok {
			t.Fatalf("unexpected back: %v", back)
		}
		return
	}

	wantFront := items[0]
	wantBack := items[len(items)-1]
	gotFront, okFront := om.Front()
	gotBack, okBack := om.Back()

	if !okFront {
		t.Fatalf("front item not present")
	}
	if !okBack {
		t.Fatalf("back item not present")
	}
	if diff := cmp.Diff(wantFront, gotFront); diff != "" {
		t.Fatalf("unexpected front (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(wantBack, gotBack); diff != "" {
		t.Fatalf("unexpected back (-want +got):\n%s", diff)
	}
	if _, ok := om.Prev(gotFront); ok {
		t.Fatal("front element has previous element")
	}
	if _, ok := om.Next(gotBack); ok {
		t.Fatal("back element has next element")
	}
}

func checkKeys[K comparable, V any](t *testing.T, om *OrderedMap[K, V], items []Item[K, V]) {
	t.Helper()

	keys := make([]K, 0, len(items))
	for _, item := range items {
		keys = append(keys, item.Key)
	}
	if diff := cmp.Diff(keys, om.Keys()); diff != "" {
		t.Fatalf("unexpected keys (-want +got):\n%s", diff)
	}
}

func checkPrevNext[K comparable, V any](t *testing.T, om *OrderedMap[K, V], items []Item[K, V]) {
	t.Helper()

	// iterate front to back
	{
		got := make([]Item[K, V], 0, len(items))
		for e, ok := om.Front(); ok; e, ok = om.Next(e) {
			got = append(got, e)
		}

		if diff := cmp.Diff(items, got); diff != "" {
			t.Fatalf("unexpected items (-want +got):\n%s", diff)
		}
	}

	// iterate back to front
	{
		got := make([]Item[K, V], 0, len(items))
		for e, ok := om.Back(); ok; e, ok = om.Prev(e) {
			got = append(got, e)
		}

		// reverse slice
		for i, j := 0, len(got)-1; i < j; i, j = i+1, j-1 {
			got[i], got[j] = got[j], got[i]
		}

		if diff := cmp.Diff(items, got); diff != "" {
			t.Fatalf("unexpected items (-want +got):\n%s", diff)
		}
	}
}
