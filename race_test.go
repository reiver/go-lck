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
// Running `go test` without the `-race` won't test thread safety.

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

// TestRace_Map_ForSnapshotConsistency verifies that a single For() call
// observes an internally-consistent snapshot even while other goroutines
// are mutating the map: no duplicate keys in one iteration, and every
// observed (k, v) pair respects the writer invariant v == k. A broken
// snapshot (e.g. iterating the live map without copying) would eventually
// surface duplicates under contention.
func TestRace_Map_ForSnapshotConsistency(t *testing.T) {
	var m lck.Map[int, int]

	const keySpace = 500

	stop := make(chan struct{})

	var writerWG sync.WaitGroup
	for g := 0; g < 8; g++ {
		writerWG.Add(1)
		go func(base int) {
			defer writerWG.Done()
			k := base
			for {
				select {
				case <-stop:
					return
				default:
				}
				k = (k + 1) % keySpace
				if k%3 == 0 {
					m.Unset(k)
				} else {
					m.Set(k, k) // writer invariant: value == key
				}
			}
		}(g)
	}

	var readerWG sync.WaitGroup
	for g := 0; g < 8; g++ {
		readerWG.Add(1)
		go func() {
			defer readerWG.Done()
			for i := 0; i < 50; i++ {
				seen := map[int]struct{}{}
				m.For(func(k int, v int) {
					if v != k {
						t.Errorf("snapshot inconsistency: key=%d value=%d (want value==key)", k, v)
					}
					if _, dup := seen[k]; dup {
						t.Errorf("duplicate key in snapshot: %d", k)
					}
					seen[k] = struct{}{}
				})
			}
		}()
	}

	readerWG.Wait()
	close(stop)
	writerWG.Wait()
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

// TestRace_Lockable_LetMutationVsGet exercises Let concurrently with Get
// and — critically — asserts that Let performs an atomic read-modify-write.
// N goroutines each increment the value M times via Let(v+1); the final
// value MUST equal N*M. A non-atomic Let would lose updates and this test
// would fail deterministically (the race detector alone would not catch it).
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

	got := l.Get()
	want := raceGoroutines * raceIterations
	if got != want {
		t.Errorf("Let is not atomic: got %d, want %d (lost %d updates)", got, want, want-got)
	}
}

// TestRace_Lockable_WideStructNoTearing stores a struct wide enough that a
// naive (unlocked) read would not be atomic at the hardware level. The
// invariant is that all fields are always equal; writers only ever publish
// "all-equal" values. If Lockable's locking were broken, a reader would
// eventually observe a struct with mismatched fields — a torn read.
//
// Using Lockable[int] would not catch this: on common 64-bit architectures
// aligned int stores are atomic, so the race test would pass even with the
// mutex removed. A wide struct forces the locking to actually matter.
func TestRace_Lockable_WideStructNoTearing(t *testing.T) {
	type wide struct {
		a, b, c, d, e, f, g, h int64
	}

	var l lck.Lockable[wide]

	stop := make(chan struct{})

	var writerWG sync.WaitGroup
	for g := 0; g < raceGoroutines; g++ {
		writerWG.Add(1)
		go func(seed int64) {
			defer writerWG.Done()
			v := seed
			for {
				select {
				case <-stop:
					return
				default:
				}
				l.Set(wide{v, v, v, v, v, v, v, v})
				v++
			}
		}(int64(g) * 1_000_000)
	}

	var readerWG sync.WaitGroup
	for g := 0; g < raceGoroutines; g++ {
		readerWG.Add(1)
		go func() {
			defer readerWG.Done()
			for i := 0; i < raceIterations*10; i++ {
				w := l.Get()
				if w.a != w.b || w.b != w.c || w.c != w.d ||
					w.d != w.e || w.e != w.f || w.f != w.g || w.g != w.h {
					t.Errorf("torn read (locking broken): %+v", w)
					return
				}
			}
		}()
	}

	readerWG.Wait()
	close(stop)
	writerWG.Wait()
}
