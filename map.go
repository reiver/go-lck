package lck

import (
	"cmp"
	"slices"
	"sync"
)

type Map[K cmp.Ordered, V any] struct {
	mutex sync.Mutex
	data map[K]V
}

func (receiver *Map[K,V]) For(fn func(K,V)) {
	if nil == receiver {
		return
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	var keys []K
	for key,_ := range receiver.data {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	for _, key := range keys {
		value := receiver.data[key]
		fn(key,value)
	}
}

func (receiver *Map[K,V]) Get(key K) (V, bool) {
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

func (receiver *Map[K,V]) Len() int {
	if nil == receiver {
		return 0
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return len(receiver.data)
}

func (receiver *Map[K,V]) Set(key K, value V) {
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

func (receiver *Map[K,V]) Unset(key K) {
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
