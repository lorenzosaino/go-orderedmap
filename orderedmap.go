// Package ordermap implements an ordered map.
//
// An ordered map is a map whose values are all connected
// with a doubly-linked list. It provides O(1) lookup,
// removal, insertion and move to front/back.
//
// To iterate over a map (where m is a *OrderedMap):
//	for e, ok := m.Front(); ok; e, ok = m.Next(e) {
//		// do something with e
//	}
//
package orderedmap

import (
	"errors"

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

// Remove removes an item from a map and returns the value removed.
//
// If the item to be deleted was already missing from the map, ok is set to false.
func (m *OrderedMap[K, V]) Remove(key K) (value V, ok bool) {
	el, ok := m.m[key]
	if !ok {
		return value, false
	}
	val := m.l.Remove(el)
	delete(m.m, key)
	return val.Value, true
}

// Len returns the number of items stored in the ordered map
func (m *OrderedMap[K, V]) Len() int {
	return len(m.m)
}

// Clear empties the ordered map
func (m *OrderedMap[K, V]) Clear() {
	m.m = make(map[K]*list.Element[Item[K, V]])
	m.l.Init()
}

// Map returns a map of all items stored in the OrderedMap.
func (m *OrderedMap[K, V]) Map() map[K]V {
	out := make(map[K]V, m.l.Len())
	for k, v := range m.m {
		out[k] = v.Value.Value
	}
	return out
}

// Item returns the a ordered slice of keys of the content of the map
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

// Item returns the a ordered slice of items of the content of the map
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
func (m *OrderedMap[K, V]) Next(item Item[K, V]) (next Item[K, V], ok bool) {
	e, ok := m.m[item.Key]
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
func (m *OrderedMap[K, V]) Prev(item Item[K, V]) (prev Item[K, V], ok bool) {
	e, ok := m.m[item.Key]
	if !ok {
		return prev, false
	}
	e = e.Prev()
	if e == nil {
		return prev, false
	}
	return e.Value, true
}
