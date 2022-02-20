package orderedmap

import (
	"fmt"
)

func ExampleOrderedMap_iteration() {
	m := New[int, string]()
	m.PushBack(1, "one")
	m.PushBack(2, "two")
	m.PushBack(3, "three")
	for e, ok := m.Front(); ok; e, ok = m.Next(e.Key) {
		fmt.Println(e.Key, e.Value)
	}
	// Output:
	// 1 one
	// 2 two
	// 3 three
}

func ExampleOrderedMap_reverse_iteration() {
	m := New[int, string]()
	m.PushBack(1, "one")
	m.PushBack(2, "two")
	m.PushBack(3, "three")
	for e, ok := m.Back(); ok; e, ok = m.Prev(e.Key) {
		fmt.Println(e.Key, e.Value)
	}
	// Output:
	// 3 three
	// 2 two
	// 1 one
}

func ExampleOrderedMap_Filter() {
	m := New[int, string]()
	m.PushBack(1, "one")
	m.PushBack(2, "two")
	m.PushBack(3, "three")
	m.PushBack(4, "four")

	isKeyEven := func(key int, val string) bool {
		return key%2 == 0
	}

	filteredMap := m.Filter(isKeyEven)

	for e, ok := filteredMap.Front(); ok; e, ok = filteredMap.Next(e.Key) {
		fmt.Println(e.Key, e.Value)
	}
	// Output:
	// 2 two
	// 4 four
}
