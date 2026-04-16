/*
Package lck implements thread-safe locking-types, for the Go programming language.

Package lck provides two generic wrappers that serialize concurrent access to a shared value: [Lockable] for a single value, and [Map] for a keyed collection of values.
Both types are ready to use at their zero value; no constructor is required.

A [Lockable] holds one value of type T, guarded by a read/write mutex.
It supports atomic read, write, exchange, and read-modify-write operations.

A [Map] holds a keyed collection of values, guarded by a mutex.
It supports all the operations of [Lockable] on a per-key basis, plus bulk operations such as iteration, key enumeration, and clearing.

Both types must not be copied after first use; instead pass them by pointer.

Example usage:

	var counter lck.Lockable[int]
	counter.Set(0)

	var inventory lck.Map[string, int]
	inventory.Set("apple", 5)

See also:

	• [Lockable]
	• [Map]
*/
package lck
