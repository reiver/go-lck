package lck

import (
	"sync"
)

// Lockable wraps a single value of type T with a read/write mutex, providing thread-safe read, write, exchange, and read-modify-write operations.
//
// The zero value is ready to use.
// A Lockable must not be copied after first use; instead pass it by pointer.
//
// Example usage:
//
//	var status lck.Lockable[string]
//	status.Set("ready")
//	current := status.Get()
//
// See also:
//
//	• [Lockable.Get]
//	• [Lockable.Let]
//	• [Lockable.Set]
//	• [Lockable.Swap]
type Lockable[T comparable] struct {
	value T
	mutex sync.RWMutex
}

// Get returns the current value.
//
// Get acquires a read lock; multiple concurrent Get calls may proceed in parallel.
// If the receiver is nil, Get returns the zero value of T.
//
// Example usage:
//
//	var flag lck.Lockable[bool]
//	flag.Set(true)
//	value := flag.Get()
func (receiver *Lockable[T]) Get() T {
	if nil == receiver {
		var nada T
		return nada
	}

	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()

	return receiver.value
}

// Let atomically replaces the current value with fn(current) and returns the previous value.
//
// The callback funtion fn runs while the write lock is held; fn MUST NOT call any method on the same Lockable — doing so will deadlock.
// Panics inside fn correctly release the lock.
//
// If the receiver is nil, Let returns the zero value of T and does not call the function fn.
//
// Example usage:
//
//	var counter lck.Lockable[int]
//	prev := counter.Let(func(v int) int {
//		return v + 1
//	})
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

// Set stores value, replacing any previous value.
//
// If the receiver is nil, Set does nothing.
//
// Example usage:
//
//	var name lck.Lockable[string]
//	name.Set("alice")
func (receiver *Lockable[T]) Set(value T) {
	if nil == receiver {
		return
	}

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.value = value
}

// Swap stores value and returns the previous value.
//
// If the receiver is nil, Swap returns the zero value of T and does not store anything.
//
// Example usage:
//
//	var cache lck.Lockable[string]
//	old := cache.Swap("new-value")
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
