package orderedmap

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEmpty(t *testing.T) {
	m := New[int, string]()
	checkAll(t, m, []Item[int, string]{})

	m.Clear()
	checkAll(t, m, []Item[int, string]{})
}

func TestClear(t *testing.T) {
	m := newFromItems(t, []Item[int, string]{{1, "one"}, {2, "two"}})
	checkAll(t, m, []Item[int, string]{{1, "one"}, {2, "two"}})

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

func TestGet(t *testing.T) {
	cases := []struct {
		name      string
		items     []Item[int, string]
		key       int
		wantValue string
		ok        bool
	}{
		{
			name:  "empty",
			items: []Item[int, string]{},
			key:   1,
		},
		{
			name:      "existing key",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}},
			key:       1,
			wantValue: "one",
			ok:        true,
		},
		{
			name:  "missing key",
			items: []Item[int, string]{{1, "one"}, {2, "two"}},
			key:   3,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.items)
			gotValue, ok := m.Get(c.key)
			if ok != c.ok {
				t.Fatalf("unexpected ok: want: %t, got %t", c.ok, ok)
			}
			if gotValue != c.wantValue {
				t.Fatalf("unexpected value: want: %v, got %v", c.wantValue, gotValue)
			}
			// validate that the map was not modified by the Get operation
			checkAll(t, m, c.items)
		})
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name  string
		items []Item[int, string]
		key   int
		value string
		want  []Item[int, string]
		err   error
	}{
		{
			name:  "empty",
			items: []Item[int, string]{},
			key:   1,
			want:  []Item[int, string]{},
			err:   ErrKeyMissing,
		},
		{
			name:  "update key",
			items: []Item[int, string]{{1, "one"}, {2, "two"}},
			key:   1,
			value: "newone",
			want:  []Item[int, string]{{1, "newone"}, {2, "two"}},
		},
		{
			name:  "no-op",
			items: []Item[int, string]{{1, "one"}, {2, "two"}},
			key:   1,
			value: "one",
			want:  []Item[int, string]{{1, "one"}, {2, "two"}},
		},
		{
			name:  "missing key",
			items: []Item[int, string]{{1, "one"}, {2, "two"}},
			key:   3,
			want:  []Item[int, string]{{1, "one"}, {2, "two"}},
			err:   ErrKeyMissing,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.items)
			wantOldValue, _ := m.Get(c.key)
			gotOldValue, err := m.Update(c.key, c.value)
			if !errors.Is(err, c.err) {
				t.Fatalf("unexpected err: want: %v, got %v", c.err, err)
			}
			if gotOldValue != wantOldValue {
				t.Fatalf("unexpected old value: want: %v, got %v", wantOldValue, gotOldValue)
			}
			checkAll(t, m, c.want)
		})
	}
}

func TestPushBack(t *testing.T) {
	cases := []struct {
		name       string
		items      []Item[int, string]
		itemToPush Item[int, string]
		want       []Item[int, string]
		err        error
	}{
		{
			name:       "empty",
			itemToPush: Item[int, string]{1, "one"},
			want:       []Item[int, string]{{1, "one"}},
		},
		{
			name:       "existing key",
			items:      []Item[int, string]{{1, "one"}, {2, "two"}},
			itemToPush: Item[int, string]{1, "one"},
			want:       []Item[int, string]{{1, "one"}, {2, "two"}},
			err:        ErrKeyAlreadyPresent,
		},
		{
			name:       "push new key",
			items:      []Item[int, string]{{1, "one"}, {2, "two"}},
			itemToPush: Item[int, string]{3, "three"},
			want:       []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.items)
			if err := m.PushBack(c.itemToPush.Key, c.itemToPush.Value); !errors.Is(err, c.err) {
				t.Fatalf("unexpected error: want: %v, got %v", c.err, err)
			}
			checkAll(t, m, c.want)
		})
	}
}

func TestPushFront(t *testing.T) {
	cases := []struct {
		name       string
		items      []Item[int, string]
		itemToPush Item[int, string]
		want       []Item[int, string]
		err        error
	}{
		{
			name:       "empty",
			itemToPush: Item[int, string]{1, "one"},
			want:       []Item[int, string]{{1, "one"}},
		},
		{
			name:       "existing key",
			items:      []Item[int, string]{{1, "one"}, {2, "two"}},
			itemToPush: Item[int, string]{1, "one"},
			want:       []Item[int, string]{{1, "one"}, {2, "two"}},
			err:        ErrKeyAlreadyPresent,
		},
		{
			name:       "push new key",
			items:      []Item[int, string]{{2, "two"}, {3, "three"}},
			itemToPush: Item[int, string]{1, "one"},
			want:       []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.items)
			if err := m.PushFront(c.itemToPush.Key, c.itemToPush.Value); !errors.Is(err, c.err) {
				t.Fatalf("unexpected error: want: %v, got %v", c.err, err)
			}
			checkAll(t, m, c.want)
		})
	}
}

func TestInsertAfter(t *testing.T) {
	cases := []struct {
		name         string
		items        []Item[int, string]
		itemToInsert Item[int, string]
		mark         int
		want         []Item[int, string]
		err          error
	}{
		{
			name:         "empty",
			itemToInsert: Item[int, string]{1, "one"},
			mark:         2,
			want:         []Item[int, string]{},
			err:          ErrMarkKeyMissing,
		},
		{
			name:         "empty and mark equals key",
			itemToInsert: Item[int, string]{2, "two"},
			mark:         2,
			want:         []Item[int, string]{},
			err:          ErrMarkKeyMissing,
		},
		{
			name:         "one item",
			items:        []Item[int, string]{{1, "one"}},
			itemToInsert: Item[int, string]{2, "two"},
			mark:         1,
			want:         []Item[int, string]{{1, "one"}, {2, "two"}},
		},
		{
			name:         "insert at back",
			items:        []Item[int, string]{{1, "one"}, {2, "two"}},
			itemToInsert: Item[int, string]{3, "three"},
			mark:         2,
			want:         []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
		},
		{
			name:         "insert in the middle",
			items:        []Item[int, string]{{1, "one"}, {3, "three"}},
			itemToInsert: Item[int, string]{2, "two"},
			mark:         1,
			want:         []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
		},
		{
			name:         "mark missing",
			items:        []Item[int, string]{{1, "one"}, {2, "two"}},
			itemToInsert: Item[int, string]{3, "three"},
			mark:         4,
			want:         []Item[int, string]{{1, "one"}, {2, "two"}},
			err:          ErrMarkKeyMissing,
		},
		{
			name:         "key already present",
			items:        []Item[int, string]{{1, "one"}, {2, "two"}},
			itemToInsert: Item[int, string]{2, "two"},
			mark:         1,
			want:         []Item[int, string]{{1, "one"}, {2, "two"}},
			err:          ErrKeyAlreadyPresent,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.items)
			if err := m.InsertAfter(c.itemToInsert.Key, c.itemToInsert.Value, c.mark); !errors.Is(err, c.err) {
				t.Fatalf("unexpected error: want: %v, got %v", c.err, err)
			}
			checkAll(t, m, c.want)
		})
	}
}

func TestInsertBefore(t *testing.T) {
	cases := []struct {
		name         string
		items        []Item[int, string]
		itemToInsert Item[int, string]
		mark         int
		want         []Item[int, string]
		err          error
	}{
		{
			name:         "empty",
			itemToInsert: Item[int, string]{1, "one"},
			mark:         2,
			want:         []Item[int, string]{},
			err:          ErrMarkKeyMissing,
		},
		{
			name:         "empty and mark equals key",
			itemToInsert: Item[int, string]{2, "two"},
			mark:         2,
			want:         []Item[int, string]{},
			err:          ErrMarkKeyMissing,
		},
		{
			name:         "one item",
			items:        []Item[int, string]{{2, "two"}},
			itemToInsert: Item[int, string]{1, "one"},
			mark:         2,
			want:         []Item[int, string]{{1, "one"}, {2, "two"}},
		},
		{
			name:         "insert at front",
			items:        []Item[int, string]{{2, "two"}, {3, "three"}},
			itemToInsert: Item[int, string]{1, "one"},
			mark:         2,
			want:         []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
		},
		{
			name:         "insert in the middle",
			items:        []Item[int, string]{{1, "one"}, {3, "three"}},
			itemToInsert: Item[int, string]{2, "two"},
			mark:         3,
			want:         []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
		},
		{
			name:         "mark missing",
			items:        []Item[int, string]{{1, "one"}, {2, "two"}},
			itemToInsert: Item[int, string]{3, "three"},
			mark:         4,
			want:         []Item[int, string]{{1, "one"}, {2, "two"}},
			err:          ErrMarkKeyMissing,
		},
		{
			name:         "key already present",
			items:        []Item[int, string]{{1, "one"}, {2, "two"}},
			itemToInsert: Item[int, string]{1, "one"},
			mark:         2,
			want:         []Item[int, string]{{1, "one"}, {2, "two"}},
			err:          ErrKeyAlreadyPresent,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.items)
			if err := m.InsertBefore(c.itemToInsert.Key, c.itemToInsert.Value, c.mark); !errors.Is(err, c.err) {
				t.Fatalf("unexpected error: want: %v, got %v", c.err, err)
			}
			checkAll(t, m, c.want)
		})
	}
}

func TestMoveToFront(t *testing.T) {
	cases := []struct {
		name      string
		items     []Item[int, string]
		keyToMove int
		want      []Item[int, string]
		err       error
	}{
		{
			name:      "empty",
			keyToMove: 1,
			want:      []Item[int, string]{},
			err:       ErrKeyMissing,
		},
		{
			name:      "no-op",
			items:     []Item[int, string]{{1, "one"}},
			keyToMove: 1,
			want:      []Item[int, string]{{1, "one"}},
		},
		{
			name:      "move from front",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 1,
			want:      []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
		},
		{
			name:      "move from middle",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 2,
			want:      []Item[int, string]{{2, "two"}, {1, "one"}, {3, "three"}},
		},
		{
			name:      "move from back",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 3,
			want:      []Item[int, string]{{3, "three"}, {1, "one"}, {2, "two"}},
		},
		{
			name:      "missing key",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 4,
			want:      []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			err:       ErrKeyMissing,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.items)
			if err := m.MoveToFront(c.keyToMove); !errors.Is(err, c.err) {
				t.Fatalf("unexpected error: want: %v, got %v", c.err, err)
			}
			checkAll(t, m, c.want)
		})
	}
}

func TestMoveToBack(t *testing.T) {
	cases := []struct {
		name      string
		items     []Item[int, string]
		keyToMove int
		want      []Item[int, string]
		err       error
	}{
		{
			name:      "empty",
			keyToMove: 1,
			want:      []Item[int, string]{},
			err:       ErrKeyMissing,
		},
		{
			name:      "no-op",
			items:     []Item[int, string]{{1, "one"}},
			keyToMove: 1,
			want:      []Item[int, string]{{1, "one"}},
		},
		{
			name:      "move from front",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 1,
			want:      []Item[int, string]{{2, "two"}, {3, "three"}, {1, "one"}},
		},
		{
			name:      "move from middle",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 2,
			want:      []Item[int, string]{{1, "one"}, {3, "three"}, {2, "two"}},
		},
		{
			name:      "move from back",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 3,
			want:      []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
		},
		{
			name:      "missing key",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 4,
			want:      []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			err:       ErrKeyMissing,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.items)
			if err := m.MoveToBack(c.keyToMove); !errors.Is(err, c.err) {
				t.Fatalf("unexpected error: want: %v, got %v", c.err, err)
			}
			checkAll(t, m, c.want)
		})
	}
}

func TestMoveAfter(t *testing.T) {
	cases := []struct {
		name      string
		items     []Item[int, string]
		keyToMove int
		mark      int
		want      []Item[int, string]
		err       error
	}{
		{
			name:      "empty",
			keyToMove: 1,
			mark:      2,
			want:      []Item[int, string]{},
			err:       ErrKeyMissing,
		},
		{
			name:      "one item",
			items:     []Item[int, string]{{1, "one"}},
			keyToMove: 1,
			mark:      1,
			want:      []Item[int, string]{{1, "one"}},
		},
		{
			name:      "mark equals key",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 2,
			mark:      2,
			want:      []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
		},
		{
			name:      "move front to middle",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 1,
			mark:      2,
			want:      []Item[int, string]{{2, "two"}, {1, "one"}, {3, "three"}},
		},
		{
			name:      "move front to back",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 1,
			mark:      3,
			want:      []Item[int, string]{{2, "two"}, {3, "three"}, {1, "one"}},
		},
		{
			name:      "move middle to back",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 2,
			mark:      3,
			want:      []Item[int, string]{{1, "one"}, {3, "three"}, {2, "two"}},
		},
		{
			name:      "move back to middle",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 3,
			mark:      1,
			want:      []Item[int, string]{{1, "one"}, {3, "three"}, {2, "two"}},
		},
		{
			name:      "move back to back",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 3,
			mark:      2,
			want:      []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
		},
		{
			name:      "missing key",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 4,
			mark:      2,
			want:      []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			err:       ErrKeyMissing,
		},
		{
			name:      "missing mark",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 2,
			mark:      4,
			want:      []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			err:       ErrKeyMissing,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.items)
			if err := m.MoveAfter(c.keyToMove, c.mark); !errors.Is(err, c.err) {
				t.Fatalf("unexpected error: want: %v, got %v", c.err, err)
			}
			checkAll(t, m, c.want)
		})
	}
}

func TestMoveBefore(t *testing.T) {
	cases := []struct {
		name      string
		items     []Item[int, string]
		keyToMove int
		mark      int
		want      []Item[int, string]
		err       error
	}{
		{
			name:      "empty",
			keyToMove: 1,
			mark:      2,
			want:      []Item[int, string]{},
			err:       ErrKeyMissing,
		},
		{
			name:      "one item",
			items:     []Item[int, string]{{1, "one"}},
			keyToMove: 1,
			mark:      1,
			want:      []Item[int, string]{{1, "one"}},
		},
		{
			name:      "mark equals key",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 2,
			mark:      2,
			want:      []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
		},
		{
			name:      "move back to middle",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 3,
			mark:      2,
			want:      []Item[int, string]{{1, "one"}, {3, "three"}, {2, "two"}},
		},
		{
			name:      "move back to front",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 3,
			mark:      1,
			want:      []Item[int, string]{{3, "three"}, {1, "one"}, {2, "two"}},
		},
		{
			name:      "move middle to front",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 2,
			mark:      1,
			want:      []Item[int, string]{{2, "two"}, {1, "one"}, {3, "three"}},
		},
		{
			name:      "move front to middle",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 1,
			mark:      3,
			want:      []Item[int, string]{{2, "two"}, {1, "one"}, {3, "three"}},
		},
		{
			name:      "move front to front",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 1,
			mark:      2,
			want:      []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
		},
		{
			name:      "missing key",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 4,
			mark:      2,
			want:      []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			err:       ErrKeyMissing,
		},
		{
			name:      "missing mark",
			items:     []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToMove: 2,
			mark:      4,
			want:      []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			err:       ErrKeyMissing,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.items)
			if err := m.MoveBefore(c.keyToMove, c.mark); !errors.Is(err, c.err) {
				t.Fatalf("unexpected error: want: %v, got %v", c.err, err)
			}
			checkAll(t, m, c.want)
		})
	}
}

func TestFilter(t *testing.T) {
	cases := []struct {
		name   string
		in     []Item[int, string]
		filter func(int, string) bool
		want   []Item[int, string]
	}{
		{
			name: "no filter",
			in:   []Item[int, string]{{1, "once"}, {2, "two"}},
			want: []Item[int, string]{{1, "once"}, {2, "two"}},
		},
		{
			name:   "exclude all",
			in:     []Item[int, string]{{1, "once"}, {2, "two"}},
			filter: func(_ int, _ string) bool { return false },
			want:   []Item[int, string]{},
		},
		{
			name:   "include all",
			in:     []Item[int, string]{{1, "once"}, {2, "two"}},
			filter: func(_ int, _ string) bool { return true },
			want:   []Item[int, string]{{1, "once"}, {2, "two"}},
		},
		{
			name:   "full test",
			in:     []Item[int, string]{{1, "once"}, {2, "two"}, {3, "three"}},
			filter: func(k int, _ string) bool { return k%2 == 0 },
			want:   []Item[int, string]{{2, "two"}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.in)
			got := m.Filter(c.filter)
			checkAll(t, got, c.want)
		})
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		name        string
		items       []Item[int, string]
		keyToRemove int
		value       string
		ok          bool
		want        []Item[int, string]
	}{
		{
			name:        "empty",
			keyToRemove: 1,
			value:       "",
			want:        []Item[int, string]{},
		},
		{
			name:        "one item",
			items:       []Item[int, string]{{1, "one"}},
			keyToRemove: 1,
			value:       "one",
			ok:          true,
			want:        []Item[int, string]{},
		},
		{
			name:        "front",
			items:       []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToRemove: 1,
			value:       "one",
			ok:          true,
			want:        []Item[int, string]{{2, "two"}, {3, "three"}},
		},
		{
			name:        "middle",
			items:       []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToRemove: 2,
			value:       "two",
			ok:          true,
			want:        []Item[int, string]{{1, "one"}, {3, "three"}},
		},
		{
			name:        "back",
			items:       []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToRemove: 3,
			value:       "three",
			ok:          true,
			want:        []Item[int, string]{{1, "one"}, {2, "two"}},
		},
		{
			name:        "no-op",
			items:       []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
			keyToRemove: 4,
			want:        []Item[int, string]{{1, "one"}, {2, "two"}, {3, "three"}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.items)
			val, ok := m.Remove(c.keyToRemove)
			if val != c.value {
				t.Fatalf("unexpected value: want: %s, got: %s", c.value, val)
			}
			if ok != c.ok {
				t.Fatalf("unexpected ok: want: %t, got: %t", c.ok, ok)
			}
			checkAll(t, m, c.want)
		})
	}
}

func TestPrev(t *testing.T) {
	cases := []struct {
		name  string
		items []Item[int, string]
		key   int
		want  Item[int, string]
		ok    bool
	}{
		{
			name:  "empty",
			items: []Item[int, string]{},
			key:   1,
			want:  Item[int, string]{},
		},
		{
			name:  "single element",
			items: []Item[int, string]{{1, "one"}},
			key:   1,
			want:  Item[int, string]{},
		},
		{
			name:  "two elements",
			items: []Item[int, string]{{1, "one"}, {2, "two"}},
			key:   2,
			want:  Item[int, string]{1, "one"},
			ok:    true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.items)
			got, ok := m.Prev(c.key)
			if ok != c.ok {
				t.Fatalf("unexpected ok: want: %t, got %t", c.ok, ok)
			}
			if got != c.want {
				t.Fatalf("unexpected value: want: %v, got %v", c.want, got)
			}
			// validate that the map was not modified by the Get operation
			checkAll(t, m, c.items)
		})
	}
}

func TestNext(t *testing.T) {
	cases := []struct {
		name  string
		items []Item[int, string]
		key   int
		want  Item[int, string]
		ok    bool
	}{
		{
			name:  "empty",
			items: []Item[int, string]{},
			key:   1,
			want:  Item[int, string]{},
		},
		{
			name:  "single element",
			items: []Item[int, string]{{1, "one"}},
			key:   1,
			want:  Item[int, string]{},
		},
		{
			name:  "two elements",
			items: []Item[int, string]{{1, "one"}, {2, "two"}},
			key:   1,
			want:  Item[int, string]{2, "two"},
			ok:    true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newFromItems(t, c.items)
			got, ok := m.Next(c.key)
			if ok != c.ok {
				t.Fatalf("unexpected ok: want: %t, got %t", c.ok, ok)
			}
			if got != c.want {
				t.Fatalf("unexpected value: want: %v, got %v", c.want, got)
			}
			// validate that the map was not modified by the Get operation
			checkAll(t, m, c.items)
		})
	}
}

func newFromItems[K comparable, V any](t *testing.T, items []Item[K, V]) *OrderedMap[K, V] {
	m := New[K, V]()
	for _, item := range items {
		if err := m.PushBack(item.Key, item.Value); err != nil {
			t.Fatalf("error inserting key %v: %v", item.Key, err)
		}
	}
	return m
}

// checkAll runs all correctness checks on an orderedmap without altering its internal state
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
	if _, ok := om.Prev(gotFront.Key); ok {
		t.Fatal("front element has previous element")
	}
	if _, ok := om.Next(gotBack.Key); ok {
		t.Fatal("back element has next element")
	}
}

// checkKeys checks the correctness of the Keys() method
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

// checkKeys checks the correctness of the Prev() and Next() methods
func checkPrevNext[K comparable, V any](t *testing.T, om *OrderedMap[K, V], items []Item[K, V]) {
	t.Helper()

	// iterate front to back
	{
		got := make([]Item[K, V], 0, len(items))
		for e, ok := om.Front(); ok; e, ok = om.Next(e.Key) {
			got = append(got, e)
		}

		if diff := cmp.Diff(items, got); diff != "" {
			t.Fatalf("unexpected items (-want +got):\n%s", diff)
		}
	}

	// iterate back to front
	{
		got := make([]Item[K, V], 0, len(items))
		for e, ok := om.Back(); ok; e, ok = om.Prev(e.Key) {
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
