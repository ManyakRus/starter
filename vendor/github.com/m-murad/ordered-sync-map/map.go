package ordered_sync_map

import (
	"container/list"
	"sync"
)

type mapElement[K comparable, V any] struct {
	key   K
	value V
}

// Map is a thread safe and ordered implementation of standard map.
// K is the type of key and V is the type of value.
type Map[K comparable, V any] struct {
	mp  map[K]*list.Element
	mu  sync.RWMutex
	dll *list.List
}

// New returns an initialized Map[K, V].
func New[K comparable, V any]() *Map[K, V] {
	m := new(Map[K, V])
	m.mp = make(map[K]*list.Element)
	m.dll = list.New()
	return m
}

// Get returns the value stored in the map for a key.
// If the key is not found in the Map it return the zero value of type V.
// The bool indicates whether value was found in the map.
func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.mp[key]
	if !ok {
		var value V
		return value, ok
	}

	me := v.Value.(mapElement[K, V])
	return me.value, ok
}

// Put sets the value for the given key.
// It will replace the value if the key already exists in the map
// even if the values are same.
func (m *Map[K, V]) Put(key K, val V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if e, ok := m.mp[key]; !ok {
		m.mp[key] = m.dll.PushFront(mapElement[K, V]{key: key, value: val})
	} else {
		e.Value = mapElement[K, V]{key: key, value: val}
	}
}

// Delete deletes the value for a key.
// It returns a boolean indicating weather the key existed and it was deleted.
func (m *Map[K, V]) Delete(key K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	e, ok := m.mp[key]
	if !ok {
		return false
	}

	m.dll.Remove(e)
	delete(m.mp, key)
	return true
}

// UnorderedRange will range over the map in an unordered sequence.
// This is same as ranging over a map using the "for range" syntax.
// Parameter func f should not call any method of the Map, eg Get, Put, Delete, UnorderedRange, OrderedRange etc
// It will cause a deadlock.
func (m *Map[K, V]) UnorderedRange(f func(key K, value V)) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.mp {
		f(k, v.Value.(mapElement[K, V]).value)
	}
}

// OrderedRange will range over the map in ab ordered sequence.
// Parameter func f should not call any method of the Map, eg Get, Put, Delete, UnorderedRange, OrderedRange etc
// It will cause a deadlock.
func (m *Map[K, V]) OrderedRange(f func(key K, value V)) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cur := m.dll.Back()
	for cur != nil {
		me := cur.Value.(mapElement[K, V])
		f(me.key, me.value)
		cur = cur.Prev()
	}
}

// Length will return the length of Map.
func (m *Map[k, V]) Length() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.dll.Len()
}

// GetOrPut will return the existing value if the key exists in the Map.
// If the key did not exist previously it will be added to the Map.
// updated will be true if the key existed previously
// otherwise it will be false if the key did not exist and was added to the Map.
func (m *Map[K, V]) GetOrPut(key K, value V) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if e, exists := m.mp[key]; exists {
		me := e.Value.(mapElement[K, V])
		return me.value, true
	} else {
		m.mp[key] = m.dll.PushFront(mapElement[K, V]{key: key, value: value})
		return value, false
	}
}

// GetAndDelete will get the value saved against the given key.
// deleted will be true if the key existed previously
// otherwise it will be false.
func (m *Map[K, V]) GetAndDelete(key K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if e, exists := m.mp[key]; exists {
		m.dll.Remove(e)
		delete(m.mp, key)
		me := e.Value.(mapElement[K, V])
		return me.value, true
	} else {
		var value V
		return value, false
	}
}
