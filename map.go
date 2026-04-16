package lck

import (
	"cmp"
	"slices"
	"sync"
)

// Map is a thread-safe generic map keyed by K and holding values of type V.
// All operations are guarded by an internal mutex.
//
// The zero value is ready to use.
// A Map must not be copied after first use; instead pass it by pointer.
//
// Example usage:
//
//	var inventory lck.Map[string, int]
//	inventory.Set("apple", 5)
//	count, found := inventory.Get("apple")
//
// See also:
//
//	• [Map.Clear]
//	• [Map.For]
//	• [Map.Get]
//	• [Map.Has]
//	• [Map.Keys]
//	• [Map.Let]
//	• [Map.Len]
//	• [Map.Set]
//	• [Map.Swap]
//	• [Map.Unset]
type Map[K cmp.Ordered, V any] struct {
	mutex sync.Mutex
	data  map[K]V
}

// Clear removes every entry from the map.
//
// The underlying hash table remains allocated, so subsequent inserts do not pay reallocation cost.
// If the receiver is nil, or the map has never held entries, Clear does nothing.
//
// Example usage:
//
//	var cache lck.Map[string, int]
//	cache.Set("a", 1)
//	cache.Clear()
func (receiver *Map[K, V]) Clear() {
	if nil == receiver {
		return
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if nil == receiver.data {
		return
	}

	clear(receiver.data)
}

// For iterates over the map, invoking the function fn once per entry in sorted key order.
//
// For takes a snapshot of the entries while the lock is held, then  releases the lock and invokes fn on each snapshot entry.
// Mutations made by the function fn (or by other goroutines) after the snapshot is taken are not reflected in the current iteration.
//
// If the receiver is nil, For does nothing.
//
// Example usage:
//
//	inventory.For(func(key string, value int) {
//		fmt.Println(key, value)
//	})
func (receiver *Map[K, V]) For(fn func(K, V)) {
	if nil == receiver {
		return
	}

	type keyvalue struct {
		key   K
		value V
	}

	var keyvalues []keyvalue

	receiver.mutex.Lock()
	for key, value := range receiver.data {
		keyvalues = append(keyvalues, keyvalue{key, value})
	}
	receiver.mutex.Unlock()

	slices.SortFunc(keyvalues, func(a keyvalue, b keyvalue) int {
		return cmp.Compare(a.key, b.key)
	})

	for _, kv := range keyvalues {
		fn(kv.key, kv.value)
	}
}

// Get returns the value stored at key, and a flag indicating whether the key was present.
//
// If the key is absent, or the receiver is nil, Get returns the zero value of V and false.
//
// Example usage:
//
//	value, found := inventory.Get("apple")
//	if !found {
//		// key is absent
//	}
func (receiver *Map[K, V]) Get(key K) (V, bool) {
	var nada V

	if nil == receiver {
		return nada, false
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if len(receiver.data) <= 0 {
		return nada, false
	}

	value, found := receiver.data[key]
	if !found {
		return nada, false
	}

	return value, true
}

// Has reports whether key is present in the map.
//
// Has is equivalent to discarding the value returned by [Map.Get], but it avoids copying the value, which matters when V is large.
//
// If the receiver is nil, Has returns false.
//
// Example usage:
//
//	if inventory.Has("apple") {
//		// key is present
//	}
func (receiver *Map[K, V]) Has(key K) bool {
	if nil == receiver {
		return false
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if len(receiver.data) <= 0 {
		return false
	}

	_, found := receiver.data[key]
	return found
}

// Keys returns a snapshot of every key currently in the map.
//
// The returned slice is independent of the map — subsequent mutations do not affect it.
//
// If the receiver is nil, Keys returns nil.
//
// Example usage:
//
//	for _, key := range inventory.Keys() {
//		// ...
//	}
func (receiver *Map[K, V]) Keys() []K {
	if nil == receiver {
		return nil
	}

	receiver.mutex.Lock()
	keys := make([]K, 0, len(receiver.data))
	for key := range receiver.data {
		keys = append(keys, key)
	}
	receiver.mutex.Unlock()

	slices.Sort(keys)

	return keys
}

// Len returns the number of entries currently in the map.
//
// If the receiver is nil, Len returns 0.
//
// Example usage:
//
//	size := inventory.Len()
func (receiver *Map[K, V]) Len() int {
	if nil == receiver {
		return 0
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return len(receiver.data)
}

// Let atomically replaces the value at key with fn(current, found) and returns the previous value and presence flag.
//
// If key was absent, fn is called with the zero value of V and found=false.
// Let always writes fn's return value, so after Let returns the key is present with the value fn returned.
//
// The callback fn runs while the write lock is held; fn MUST NOT call any method on the same Map — doing so will deadlock.
//
// If the receiver is nil, Let returns the zero value of V and false, and does not call fn.
//
// Example usage:
//
//	prev, had := counters.Let("hits", func(v int, found bool) int {
//		return v + 1
//	})
func (receiver *Map[K, V]) Let(key K, fn func(V, bool) V) (V, bool) {
	if nil == receiver {
		var nada V
		return nada, false
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if nil == receiver.data {
		receiver.data = map[K]V{}
	}
	if nil == receiver.data {
		var nada V
		return nada, false
	}

	value, found := receiver.data[key]
	receiver.data[key] = fn(value, found)
	return value, found
}

// Set stores value at key, replacing any previous entry.
//
// If the receiver is nil, Set does nothing.
//
// Example usage:
//
//	inventory.Set("apple", 5)
func (receiver *Map[K, V]) Set(key K, value V) {
	if nil == receiver {
		return
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if nil == receiver.data {
		receiver.data = map[K]V{}
	}
	if nil == receiver.data {
		return
	}

	receiver.data[key] = value
}

// Swap stores value at key and returns the previous value and presence flag.
//
// If the key was absent, Swap returns the zero value of V and false, then inserts the new entry. After Swap returns the key is always present with the supplied value.
//
// If the receiver is nil, Swap returns the zero value of V and false, and does not store anything.
//
// Example usage:
//
//	old, found := inventory.Swap("apple", 10)
func (receiver *Map[K, V]) Swap(key K, value V) (V, bool) {
	if nil == receiver {
		var nada V
		return nada, false
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if nil == receiver.data {
		receiver.data = map[K]V{}
	}
	if nil == receiver.data {
		var nada V
		return nada, false
	}

	prev, found := receiver.data[key]
	receiver.data[key] = value
	return prev, found
}

// Unset removes the entry at key and returns its previous value and presence flag.
//
// If the key was absent, Unset returns the zero value of V and false.
// If the receiver is nil, Unset returns the zero value of V and false, and does not remove anything.
//
// Example usage:
//
//	inventory.Unset("apple")
//
// Another example usage:
//
//	old, found := inventory.Unset("apple")
func (receiver *Map[K, V]) Unset(key K) (V, bool) {
	if nil == receiver {
		var nada V
		return nada, false
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if len(receiver.data) <= 0 {
		var nada V
		return nada, false
	}

	prev, found := receiver.data[key]
	delete(receiver.data, key)
	return prev, found
}
