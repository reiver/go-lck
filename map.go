package lck

import (
	"cmp"
	"slices"
	"sync"
)

type Map[K cmp.Ordered, V any] struct {
	mutex sync.Mutex
	data  map[K]V
}

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

	return keys
}

func (receiver *Map[K, V]) Len() int {
	if nil == receiver {
		return 0
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return len(receiver.data)
}

func (receiver *Map[K, V]) Let(key K, fn func(V, bool) V) {
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

	value, found := receiver.data[key]
	receiver.data[key] = fn(value, found)
}

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

func (receiver *Map[K, V]) Unset(key K) {
	if nil == receiver {
		return
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	if len(receiver.data) <= 0 {
		return
	}

	delete(receiver.data, key)
}
