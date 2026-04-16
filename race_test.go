package lck_test

import (
	"sync"
	"testing"

	"github.com/reiver/go-lck"
)

// TO CHECK THESE TESTS, YOU MUST RUN:
//
//	go test -race
//
// Runnng `go test` without the `-race` won' test thread safety.

const raceGoroutines = 64
const raceIterations = 200

func TestRace_Map_ConcurrentSetGetUnset(t *testing.T) {
	var m lck.Map[int, int]

	var wg sync.WaitGroup
	for g := 0; g < raceGoroutines; g++ {
		wg.Add(3)

		go func(base int) {
			defer wg.Done()
			for i := 0; i < raceIterations; i++ {
				m.Set(base+i, i)
			}
		}(g * raceIterations)

		go func(base int) {
			defer wg.Done()
			for i := 0; i < raceIterations; i++ {
				_, _ = m.Get(base + i)
			}
		}(g * raceIterations)

		go func(base int) {
			defer wg.Done()
			for i := 0; i < raceIterations; i++ {
				m.Unset(base + i)
			}
		}(g * raceIterations)
	}
	wg.Wait()
}

func TestRace_Map_ForDuringWrites(t *testing.T) {
	var m lck.Map[int, int]

	// Pre-seed.
	for i := 0; i < 100; i++ {
		m.Set(i, i)
	}

	var wg sync.WaitGroup

	// Writers.
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for i := 0; i < raceIterations; i++ {
				m.Set(base+i, i)
				if i%3 == 0 {
					m.Unset(base + i)
				}
			}
		}(g * raceIterations)
	}

	// Iterators running concurrently with writers — proves the snapshot
	// isolates iteration from concurrent mutation.
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				m.For(func(k int, v int) {
					_ = k
					_ = v
				})
			}
		}()
	}

	wg.Wait()
}

func TestRace_Lockable_ConcurrentSetGet(t *testing.T) {
	var l lck.Lockable[int]

	var wg sync.WaitGroup
	for g := 0; g < raceGoroutines; g++ {
		wg.Add(2)

		go func(v int) {
			defer wg.Done()
			for i := 0; i < raceIterations; i++ {
				l.Set(v + i)
			}
		}(g)

		go func() {
			defer wg.Done()
			for i := 0; i < raceIterations; i++ {
				_ = l.Get()
			}
		}()
	}
	wg.Wait()
}

func TestRace_Lockable_LetMutationVsGet(t *testing.T) {
	var l lck.Lockable[int]

	var wg sync.WaitGroup
	for g := 0; g < raceGoroutines; g++ {
		wg.Add(2)

		go func() {
			defer wg.Done()
			for i := 0; i < raceIterations; i++ {
				l.Let(func(v int) int {
					return v + 1
				})
			}
		}()

		go func() {
			defer wg.Done()
			for i := 0; i < raceIterations; i++ {
				_ = l.Get()
			}
		}()
	}
	wg.Wait()
}
