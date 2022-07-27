package orderedmap_test

import (
	"fmt"

	"github.com/lorenzosaino/go-orderedmap"
)

func ExampleOrderedMap_iteration() {
	m := orderedmap.New[int, string]()
	m.PushBack(1, "one")
	m.PushBack(2, "two")
	m.PushBack(3, "three")

	f := func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	}
	m.Range(f)
	// Output:
	// 1 one
	// 2 two
	// 3 three
}

func ExampleOrderedMap_reverseIteration() {
	m := orderedmap.New[int, string]()
	m.PushBack(1, "one")
	m.PushBack(2, "two")
	m.PushBack(3, "three")

	f := func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	}
	m.RangeReverse(f)
	// Output:
	// 3 three
	// 2 two
	// 1 one
}

func ExampleOrderedMap_Filter() {
	m := orderedmap.New[int, string]()
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

func ExampleOrderedMap_Reverse() {
	m := orderedmap.New[int, string]()
	m.PushBack(1, "one")
	m.PushBack(2, "two")
	m.PushBack(3, "three")

	fmt.Println("original map:")
	for e, ok := m.Front(); ok; e, ok = m.Next(e.Key) {
		fmt.Println(e.Key, e.Value)
	}

	m = m.Reverse()

	fmt.Println("reversed map:")
	for e, ok := m.Front(); ok; e, ok = m.Next(e.Key) {
		fmt.Println(e.Key, e.Value)
	}
	// Output:
	// original map:
	// 1 one
	// 2 two
	// 3 three
	// reversed map:
	// 3 three
	// 2 two
	// 1 one
}
