package lck

import (
	"sync"
)

type Lockable[T any] struct {
	value T
	mutex sync.RWMutex
}

func (receiver *Lockable[T]) Get() T {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()

	return receiver.value
}

func (receiver *Lockable[T]) Let(fn func(*T)) {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()

	fn(&receiver.value)
}
func (receiver *Lockable[T]) Set(value T) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.value = value
}
