package goda

import "sync"

// Iter ...
type Iter chan interface{}

// Predicate ...
type Predicate func(interface{}) bool

// Mapper ...
type Mapper func(interface{}) interface{}

// MultiMapper ...
type MultiMapper func(...interface{}) interface{}

// Reducer ...
type Reducer func(memo interface{}, element interface{}) interface{}

// New ...
func New(els ...interface{}) Iter {
	c := make(Iter)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

// Int64 ...
func Int64(els ...int64) Iter {
	c := make(Iter)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

// Int32 ...
func Int32(els ...int32) Iter {
	c := make(Iter)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

// Float64 ...
func Float64(els ...float64) Iter {
	c := make(Iter)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

// Float32 ...
func Float32(els ...float32) Iter {
	c := make(Iter)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

// Uint ...
func Uint(els ...uint) Iter {
	c := make(Iter)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

// Uint64 ...
func Uint64(els ...uint64) Iter {
	c := make(Iter)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

// Uint32 ...
func Uint32(els ...uint32) Iter {
	c := make(Iter)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

// String ...
func String(els ...string) Iter {
	c := make(Iter)
	go func() {
		for _, el := range els {
			c <- el
		}
		close(c)
	}()
	return c
}

// List ...
func (it Iter) List() List {
	arr := make([]interface{}, 0, 1)
	for el := range it {
		arr = append(arr, el)
	}
	return List(arr)
}

// Count from i to infinity
func Count(i int) Iter {
	c := make(Iter)
	go func() {
		for ; true; i++ {
			c <- i
		}
	}()
	return c
}

// Cycle through an iterator infinitely (requires memory)
func (it Iter) Cycle() Iter {
	c, a := make(Iter), make([]interface{}, 0, 1)
	go func() {
		for el := range it {
			a = append(a, el)
			c <- el
		}
		for {
			for _, el := range a {
				c <- el
			}
		}
	}()
	return c
}

// Repeat an element n times or infinitely
func Repeat(el interface{}, n ...int) Iter {
	c := make(Iter)
	go func() {
		for i := 0; len(n) == 0 || i < n[0]; i++ {
			c <- el
		}
		close(c)
	}()
	return c
}

// Chain together multiple iterators
func Chain(its ...Iter) Iter {
	c := make(Iter)
	go func() {
		for _, it := range its {
			for el := range it {
				c <- el
			}
		}
		close(c)
	}()
	return c
}

// DropWhile : Elements after pred(el) == true
func (it Iter) DropWhile(pred Predicate) Iter {
	c := make(Iter)
	go func() {
		for el := range it {
			if drop := pred(el); !drop {
				c <- el
				break
			}
		}
		for el := range it {
			c <- el
		}
		close(c)
	}()
	return c
}

// TakeWhile : Elements before pred(el) == false
func (it Iter) TakeWhile(pred Predicate) Iter {
	c := make(Iter)
	go func() {
		for el := range it {
			if take := pred(el); take {
				c <- el
			} else {
				break
			}
		}
		close(c)
	}()
	return c
}

// Filter out any elements where pred(el) == false
func (it Iter) Filter(pred Predicate) Iter {
	c := make(Iter)
	go func() {
		for el := range it {
			if keep := pred(el); keep {
				c <- el
			}
		}
		close(c)
	}()
	return c
}

// FilterFalse filters out any elements where pred(el) == true
func (it Iter) FilterFalse(pred Predicate) Iter {
	c := make(Iter)
	go func() {
		for el := range it {
			if drop := pred(el); !drop {
				c <- el
			}
		}
		close(c)
	}()
	return c
}

// Slice : Sub-iterator from start (inclusive) to [stop (exclusive) every [step (default 1)]]
func (it Iter) Slice(startstopstep ...int) Iter {
	start, stop, step := 0, 0, 1
	if len(startstopstep) == 1 {
		start = startstopstep[0]
	} else if len(startstopstep) == 2 {
		start, stop = startstopstep[0], startstopstep[1]
	} else if len(startstopstep) >= 3 {
		start, stop, step = startstopstep[0], startstopstep[1], startstopstep[2]
	}

	c := make(Iter)
	go func() {
		i := 0
		// Start
		for el := range it {
			if i >= start {
				c <- el // inclusive
				break
			}
			i++
		}

		// Stop
		i, j := i+1, 1
		for el := range it {
			if stop > 0 && i >= stop {
				break
			} else if j%step == 0 {
				c <- el
			}

			i, j = i+1, j+1
		}

		close(c)
	}()
	return c
}

// Map an iterator to fn(el) for el in it
func (it Iter) Map(fn Mapper) Iter {
	c := make(Iter)
	go func() {
		for el := range it {
			c <- fn(el)
		}
		close(c)
	}()
	return c
}

// MultiMap :
// Map p, q, ... to fn(pEl, qEl, ...)
// Breaks on first closed channel
func MultiMap(fn MultiMapper, its ...Iter) Iter {
	c := make(Iter)
	go func() {
	Outer:
		for {
			els := make([]interface{}, len(its))
			for i, it := range its {
				if el, ok := <-it; ok {
					els[i] = el
				} else {
					break Outer
				}
			}
			c <- fn(els...)
		}
		close(c)
	}()
	return c
}

// MultiMapLongest :
// Map p, q, ... to fn(pEl, qEl, ...)
// Breaks on last closed channel
func MultiMapLongest(fn MultiMapper, its ...Iter) Iter {
	c := make(Iter)
	go func() {
		for {
			els := make([]interface{}, len(its))
			n := 0
			for i, it := range its {
				if el, ok := <-it; ok {
					els[i] = el
				} else {
					n++
				}
			}
			if n < len(its) {
				c <- fn(els...)
			} else {
				break
			}
		}
		close(c)
	}()
	return c
}

// Starmap :
// Map an iterator if arrays to a fn(els...)
// Iter must be an iterator of []interface{} (possibly created by Zip)
// If not, Starmap will act like MultiMap with a single iterator
func (it Iter) Starmap(fn MultiMapper) Iter {
	c := make(Iter)
	go func() {
		for els := range it {
			if elements, ok := els.([]interface{}); ok {
				c <- fn(elements...)
			} else {
				c <- fn(els)
			}
		}
		close(c)
	}()
	return c
}

// Zip up multiple interators into one
// Close on shortest iterator
func Zip(its ...Iter) Iter {
	c := make(Iter)
	go func() {
		defer close(c)
		for {
			els := make([]interface{}, len(its))
			for i, it := range its {
				if el, ok := <-it; ok {
					els[i] = el
				} else {
					return
				}
			}
			c <- els
		}
	}()
	return c
}

// ZipLongest :
// Zip up multiple iterators into one
// Close on longest iterator
func ZipLongest(its ...Iter) Iter {
	c := make(Iter)
	go func() {
		for {
			els := make([]interface{}, len(its))
			n := 0
			for i, it := range its {
				if el, ok := <-it; ok {
					els[i] = el
				} else {
					n++
				}
			}
			if n < len(its) {
				c <- els
			} else {
				break
			}
		}
		close(c)
	}()
	return c
}

// Reduce the iterator (aka fold) from the left
func (it Iter) Reduce(red Reducer, memo interface{}) interface{} {
	for el := range it {
		memo = red(memo, el)
	}
	return memo
}

// Tee splits an iterator into n multiple iterators
// Requires memory to keep values for n iterators
func (it Iter) Tee(n int) []Iter {
	deques := make([][]interface{}, n)
	iters := make([]Iter, n)
	for i := 0; i < n; i++ {
		iters[i] = make(Iter)
	}

	mutex := new(sync.Mutex)

	gen := func(myiter Iter, i int) {
		for {
			if len(deques[i]) == 0 {
				mutex.Lock()
				if len(deques[i]) == 0 {
					if newval, ok := <-it; ok {
						for i, d := range deques {
							deques[i] = append(d, newval)
						}
					} else {
						mutex.Unlock()
						close(myiter)
						break
					}
				}
				mutex.Unlock()
			}
			var popped interface{}
			popped, deques[i] = deques[i][0], deques[i][1:]
			myiter <- popped
		}
	}
	for i, iter := range iters {
		go gen(iter, i)
	}
	return iters
}

// Tee2 helper to tee just into two iterators
func (it Iter) Tee2() (Iter, Iter) {
	iters := it.Tee(2)
	return iters[0], iters[1]
}
