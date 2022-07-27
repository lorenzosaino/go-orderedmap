// Package orderedmap implements an ordered map using generics.
//
// An ordered map is a map whose values are ordered and all connected
// with a doubly-linked list. It provides O(1) lookup,
// update, removal, insertion to front/back,
// insertion before/after a specific key,
// move to front/back, move before/after a specific key.
//
// This implementation is not safe for concurrent usage. You
// may want to use a sync.RWLock to synchronize access to it
// if you intend to use it concurrently.
package orderedmap

import (
	"errors"
	"fmt"

	"github.com/lorenzosaino/go-orderedmap/internal/list"
)

var (
	// ErrKeyMissing indicates that the key specified is not present in the ordered map
	ErrKeyMissing = errors.New("key missing")

	// ErrMarkKeyMissing indicates that the mark key specified is not present in the ordered map
	ErrMarkKeyMissing = errors.New("mark key missing")

	// ErrKeyAlreadyPresent indicates that key to be inserted is already present in the ordered map
	ErrKeyAlreadyPresent = errors.New("key already present")
)

// Item is a key-value item stored in the ordered map
type Item[K comparable, V any] struct {
	Key   K
	Value V
}

// OrderedMap is an implementation of an ordered map.
//
// K and V are respectively the types of keys and values.
type OrderedMap[K comparable, V any] struct {
	m map[K]*list.Element[Item[K, V]]
	l *list.List[Item[K, V]]
}

// New returns a new ordered map instance.
func New[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		m: make(map[K]*list.Element[Item[K, V]]),
		l: list.New[Item[K, V]](),
	}
}

// Get returns the value associated to a key in the map.
//
// If the key is not present in the map, it returns the zero value of V
// and ok is set to false.
func (m *OrderedMap[K, V]) Get(key K) (value V, ok bool) {
	if el, ok := m.m[key]; ok {
		return el.Value.Value, true
	}
	return value, false
}

// Update updates the value associated to an existing key and returns the old value.
//
// If the key is not present, then ErrKeyMissing is returned.
func (m *OrderedMap[K, V]) Update(key K, value V) (oldValue V, err error) {
	el, ok := m.m[key]
	if !ok {
		return oldValue, ErrKeyMissing
	}
	oldValue = el.Value.Value
	el.Value.Value = value
	return oldValue, nil
}

// Front returns the item at the front of the map.
//
// If the map is empty, it returns the zero value of Item[K, V]
// and ok is set to false.
func (m *OrderedMap[K, V]) Front() (item Item[K, V], ok bool) {
	if front := m.l.Front(); front != nil {
		return front.Value, true
	}
	return item, false
}

// Front returns the item at the back of the map.
//
// If the map is empty, it returns the zero value of Item[K, V]
// and ok is set to false.
func (m *OrderedMap[K, V]) Back() (item Item[K, V], ok bool) {
	if back := m.l.Back(); back != nil {
		return back.Value, true
	}
	return item, false
}

// PushFront insert a new key and value at the front of the map.
//
// It returns ErrKeyAlreadyPresent if the key to be inserted is already present.
func (m *OrderedMap[K, V]) PushFront(key K, value V) error {
	if _, ok := m.m[key]; ok {
		return ErrKeyAlreadyPresent
	}
	newVal := Item[K, V]{key, value}
	m.m[key] = m.l.PushFront(newVal)
	return nil
}

// PushBack insert a new key and value at the back of the map.
//
// It returns ErrKeyAlreadyPresent if the key to be inserted is already present.
func (m *OrderedMap[K, V]) PushBack(key K, value V) error {
	if _, ok := m.m[key]; ok {
		return ErrKeyAlreadyPresent
	}
	newVal := Item[K, V]{key, value}
	m.m[key] = m.l.PushBack(newVal)
	return nil
}

// InsertAfter insert a new key and value immediately after a mark key.
//
// It returns ErrKeyAlreadyPresent if the key to be inserted is already present
// and ErrMarkKeyMissing if the mark key is missing.
func (m *OrderedMap[K, V]) InsertAfter(key K, value V, mark K) error {
	if _, ok := m.m[key]; ok {
		return ErrKeyAlreadyPresent
	}
	markEl, ok := m.m[mark]
	if !ok {
		return ErrMarkKeyMissing
	}
	newVal := Item[K, V]{key, value}
	newEl := m.l.InsertAfter(newVal, markEl)
	m.m[key] = newEl
	return nil
}

// InsertBefore insert a new key and value immediately before a mark key.
//
// It returns ErrKeyAlreadyPresent if the key to be inserted is already present
// and ErrMarkKeyMissing if the mark key is missing.
func (m *OrderedMap[K, V]) InsertBefore(key K, value V, mark K) error {
	if _, ok := m.m[key]; ok {
		return ErrKeyAlreadyPresent
	}
	markEl, ok := m.m[mark]
	if !ok {
		return ErrMarkKeyMissing
	}
	newVal := Item[K, V]{key, value}
	newEl := m.l.InsertBefore(newVal, markEl)
	m.m[key] = newEl
	return nil
}

// MoveToFront moves an existing key to the front of the map.
//
// It returns ErrKeyMissing if the key to be moved is not in the map.
func (m *OrderedMap[K, V]) MoveToFront(key K) error {
	e, ok := m.m[key]
	if !ok {
		return ErrKeyMissing
	}
	m.l.MoveToFront(e)
	return nil
}

// MoveToBack moves an existing key to the back of the map.
//
// It returns ErrKeyMissing if the key to be moved is not in the map.
func (m *OrderedMap[K, V]) MoveToBack(key K) error {
	e, ok := m.m[key]
	if !ok {
		return ErrKeyMissing
	}
	m.l.MoveToBack(e)
	return nil
}

// MoveAfter moves an existing key immediately after a mark key.
//
// It returns ErrKeyMissing if the key to be moved is missing
// and ErrMarkKeyMissing if the mark key is missing.
func (m *OrderedMap[K, V]) MoveAfter(key K, mark K) error {
	if key == mark {
		return nil
	}
	el, ok := m.m[key]
	if !ok {
		return ErrKeyMissing
	}
	markEl, ok := m.m[mark]
	if !ok {
		return ErrKeyMissing
	}
	m.l.MoveAfter(el, markEl)
	return nil
}

// MoveAfter moves an existing key immediately before a mark key.
//
// It returns ErrKeyMissing if the key to be moved is missing
// and ErrMarkKeyMissing if the mark key is missing.
func (m *OrderedMap[K, V]) MoveBefore(key K, mark K) error {
	if key == mark {
		return nil
	}
	el, ok := m.m[key]
	if !ok {
		return ErrKeyMissing
	}
	markEl, ok := m.m[mark]
	if !ok {
		return ErrKeyMissing
	}
	m.l.MoveBefore(el, markEl)
	return nil
}

// Delete deletes an item from a map and returns the value deleted.
//
// If the item to be deleted was already missing from the map, ok is set to false.
func (m *OrderedMap[K, V]) Delete(key K) (value V, ok bool) {
	el, ok := m.m[key]
	if !ok {
		return value, false
	}
	val := m.l.Remove(el)
	delete(m.m, key)
	return val.Value, true
}

// PopFront pops the element at the front of the map and returns its value.
//
// If the map is empty, it returns the zero value of Item[K, V]
// and ok is set to false.
func (m *OrderedMap[K, V]) PopFront() (item Item[K, V], ok bool) {
	el := m.l.Front()
	if el == nil {
		return item, false
	}

	delete(m.m, el.Value.Key)
	item = m.l.Remove(el)

	return item, true
}

// PopBack pops the element at the back of the map and returns its value.
//
// If the map is empty, it returns the zero value of Item[K, V]
// and ok is set to false.
func (m *OrderedMap[K, V]) PopBack() (item Item[K, V], ok bool) {
	el := m.l.Back()
	if el == nil {
		return item, false
	}

	delete(m.m, el.Value.Key)
	item = m.l.Remove(el)

	return item, true
}

// Len returns the number of items stored in the ordered map.
func (m *OrderedMap[K, V]) Len() int {
	return len(m.m)
}

// Clear empties the ordered map.
func (m *OrderedMap[K, V]) Clear() {
	m.m = make(map[K]*list.Element[Item[K, V]])
	m.l.Init()
}

// Reverse returns a copy of the ordered map with reversed ordering.
func (m *OrderedMap[K, V]) Reverse() *OrderedMap[K, V] {
	out := New[K, V]()
	for item, ok := m.Front(); ok; item, ok = m.Next(item.Key) {
		if err := out.PushFront(item.Key, item.Value); err != nil {
			// while generally we should not panic from within a library, this
			// error should never happen because all keys of the ordered map
			// should be unique. If this error occurs, it is because of a bug
			// in this library that needs to be fixed.
			panic(fmt.Sprintf("error trying to insert key %v: %v", item.Key, err))
		}
	}
	return out
}

// Filter returns a filtered copy of the ordered map.
//
// The returned map only includes the (key, value) items such that
// f(key, value) == true.
func (m *OrderedMap[K, V]) Filter(f func(key K, value V) bool) *OrderedMap[K, V] {
	out := New[K, V]()
	for item, ok := m.Front(); ok; item, ok = m.Next(item.Key) {
		if f != nil && !f(item.Key, item.Value) {
			continue
		}
		if err := out.PushBack(item.Key, item.Value); err != nil {
			// while generally we should not panic from within a library, this
			// error should never happen because all keys of the ordered map
			// should be unique. If this error occurs, it is because of a bug
			// in this library that needs to be fixed.
			panic(fmt.Sprintf("error trying to insert key %v: %v", item.Key, err))
		}
	}
	return out
}

// Range calls f sequentially for each key and value present in the ordered map
// starting from the front element. If f returns false, Range stops the iteration.
func (m *OrderedMap[K, V]) Range(f func(key K, value V) bool) {
	for e := m.l.Front(); e != nil; e = e.Next() {
		if !f(e.Value.Key, e.Value.Value) {
			return
		}
	}
}

// Range calls f sequentially for each key and value present in the ordered map
// starting from the back element. If f returns false, RangeReverse stops the iteration.
func (m *OrderedMap[K, V]) RangeReverse(f func(key K, value V) bool) {
	for e := m.l.Back(); e != nil; e = e.Prev() {
		if !f(e.Value.Key, e.Value.Value) {
			return
		}
	}
}

// Map returns a map of all items stored in the OrderedMap.
func (m *OrderedMap[K, V]) Map() map[K]V {
	out := make(map[K]V, m.l.Len())
	for k, v := range m.m {
		out[k] = v.Value.Value
	}
	return out
}

// Item returns the a ordered slice of keys of the content of the map.
//
// Note that while this function could be used to iterate over the items
// stored in the ordered map, it allocates a new slice and copy all items
// in the map. For better performance, you may want to iterate using
// Prev() and Next() instead.
func (m *OrderedMap[K, V]) Keys() []K {
	out := make([]K, 0, m.l.Len())
	for e := m.l.Front(); e != nil; e = e.Next() {
		out = append(out, e.Value.Key)
	}
	return out
}

// Item returns the a ordered slice of items of the content of the map.
//
// Note that while this function could be used to iterate over the items
// stored in the ordered map, it allocates a new slice and copy all items
// in the map. For better performance, you may want to iterate using
// Prev() and Next() instead.
func (m *OrderedMap[K, V]) Items() []Item[K, V] {
	out := make([]Item[K, V], 0, m.l.Len())
	for e := m.l.Front(); e != nil; e = e.Next() {
		out = append(out, e.Value)
	}
	return out
}

// Next returns the item succeeding a given item in the map.
//
// If the specified item is missing or it is at the back of the map, ok is set to false.
func (m *OrderedMap[K, V]) Next(key K) (next Item[K, V], ok bool) {
	e, ok := m.m[key]
	if !ok {
		return next, false
	}
	e = e.Next()
	if e == nil {
		return next, false
	}
	return e.Value, true
}

// Prev returns the item preceding a given item in the map.
//
// If the specified item is missing or it is at the front of the map, ok is set to false.
func (m *OrderedMap[K, V]) Prev(key K) (prev Item[K, V], ok bool) {
	e, ok := m.m[key]
	if !ok {
		return prev, false
	}
	e = e.Prev()
	if e == nil {
		return prev, false
	}
	return e.Value, true
}
