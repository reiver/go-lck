package lck

import (
	"sync"
)

type Locking[T any] struct {
	value T
	mutex sync.RWMutex
}

func (receiver *Locking[T]) Get() T {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()

	return receiver.value
}

func (receiver *Locking[T]) Set(value T) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.value = value
}
