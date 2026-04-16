package lck

import (
	"sync"
)

type Lockable[T comparable] struct {
	value T
	mutex sync.RWMutex
}

func (receiver *Lockable[T]) Get() T {
	if nil == receiver {
		var nada T
		return nada
	}

	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()

	return receiver.value
}

func (receiver *Lockable[T]) Let(fn func(T) T) T {
	if nil == receiver {
		var nada T
		return nada
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	prev := receiver.value
	receiver.value = fn(receiver.value)
	return prev
}

func (receiver *Lockable[T]) Set(value T) {
	if nil == receiver {
		return
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.value = value
}

func (receiver *Lockable[T]) Swap(value T) T {
	if nil == receiver {
		var nada T
		return nada
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	prev := receiver.value
	receiver.value = value
	return prev
}
