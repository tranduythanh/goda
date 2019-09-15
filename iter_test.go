package goda

import (
	"reflect"
	"testing"
)

// Test iterators for element equality. Allow it1 to be longer than it2
func testIter(t *testing.T, it1, it2 Iter) {
	t.Log("Start")
	for el1 := range it1 {
		if el2, ok := <-it2; !ok {
			t.Error("it2 shorter than it1!", el1)
			return
		} else if !reflect.DeepEqual(el1, el2) {
			t.Error("Elements are not equal", el1, el2)
		} else {
			t.Log(el1, el2)
		}
	}
	t.Log("Stop")
}

// Test iterators for element equality. Don't allow it1 to be longer than it2
func testIterEq(t *testing.T, it1, it2 Iter) {
	t.Log("Start")
	for el1 := range it1 {
		if el2, ok := <-it2; !ok {
			t.Error("it2 shorter than it1!", el1)
			return
		} else if !reflect.DeepEqual(el1, el2) {
			t.Error("Elements are not equal", el1, el2)
		} else {
			t.Log(el1, el2)
		}
	}
	if el2, ok := <-it2; ok {
		t.Error("it1 shorter than it2!", el2)
	}
	t.Log("Stop")
}

func TestList(t *testing.T) {
	list := New(1, 2, 3).List()
	if !reflect.DeepEqual(list, []interface{}{1, 2, 3}) {
		t.Error("List didn't make a list", list)
	}
}

func TestCount(t *testing.T) {
	testIter(t, New(1, 2, 3, 4, 5, 6, 7, 8, 9), Count(1))
}

func TestCycle(t *testing.T) {
	testIter(t, New("a", "b", "ccc", "a", "b", "ccc", "a"), New("a", "b", "ccc").Cycle())
}

func TestRepeat(t *testing.T) {
	testIterEq(t, Uint64(100, 100, 100, 100), Repeat(uint64(100), 4))
	testIter(t, Uint64(100, 100, 100, 100), Repeat(uint64(100)))
}

func TestChain(t *testing.T) {
	testIterEq(t, Int32(1, 2, 3, 4, 5, 5, 4, 3, 2, 1, 100), Chain(Int32(1, 2, 3, 4, 5), Int32(5, 4, 3, 2, 1), Int32(100)))
}

func TestDropWhile(t *testing.T) {
	pred := func(i interface{}) bool {
		return i.(int) < 10
	}
	testIter(t, New(10, 11, 12, 13, 14, 15), Count(0).DropWhile(pred))
}

func TestTakeWhile(t *testing.T) {
	pred := func(i interface{}) bool {
		return i.(string)[:3] == "abc"
	}
	testIterEq(t, New("abcdef", "abcdaj"), New("abcdef", "abcdaj", "ajcde").Cycle().TakeWhile(pred))
}

func TestFilter(t *testing.T) {
	pred := func(i interface{}) bool {
		return i.(uint64)%2 == 1
	}
	testIterEq(t, Uint64(1, 3, 5, 7, 9), Uint64(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Filter(pred))
	testIterEq(t, Uint64(2, 4, 6, 8, 10), Uint64(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).FilterFalse(pred))
}

func TestSlice(t *testing.T) {
	testIter(t, New(5, 6, 7, 8, 9, 10), Count(0).Slice(5))
	testIterEq(t, New(2, 3, 4, 5, 6, 7, 8), Count(0).Slice(2, 9))
	testIterEq(t, New(3, 6, 9), Count(0).Slice(3, 11, 3))
}

func TestMap(t *testing.T) {
	mapper := func(i interface{}) interface{} {
		return len(i.(string))
	}
	testIterEq(t, New(1, 2, 3, 4), New("a", "ab", "abc", "abcd").Map(mapper))
}

func TestMultiMap(t *testing.T) {
	multiMapper := func(is ...interface{}) interface{} {
		var s float64
		for _, i := range is {
			s += i.(float64)
		}
		return s
	}
	testIterEq(t, Float64(10.4, 3.2), MultiMap(multiMapper, Float64(5.2, 1.6, 2.2), Float64(5.2, 1.0), Float64(0, 0.6, 0)))
}

func TestZip(t *testing.T) {
	a, b, c := []interface{}{1, "a"}, []interface{}{2, nil}, []interface{}{3, nil}
	test1, test2 := New(a), New(a, b, c)

	testIterEq(t, test1, Zip(Count(1), New("a")))
	testIterEq(t, test2, ZipLongest(Count(1).Slice(0, 3), New("a")))
}

func TestStarmap(t *testing.T) {
	multiMapper := func(is ...interface{}) interface{} {
		var s = 1
		for _, i := range is {
			s *= i.(int)
		}
		return s
	}
	testIterEq(t, New(10, 20, 30), Zip(New(1, 2, 3), Repeat(10, 3)).Starmap(multiMapper))
}

func TestReduce(t *testing.T) {
	summer := func(memo interface{}, el interface{}) interface{} {
		return memo.(float64) + el.(float64)
	}
	if float64(.82)-Float64(.1, .2, .3, .22).Reduce(summer, float64(0)).(float64) > .000001 {
		t.Error("Sum Reduce failed")
	}
}

func TestTee2(t *testing.T) {
	it1, it2 := New(5, 4, 3, 2, 1).Tee2()
	for i := range it1 {
		j := <-it2
		if i != j {
			t.Error("Tees are not coming off equal")
		}
	}

	it1, it2 = New(1, 2, 3, 4, 5, 6).Tee2()
	testIterEq(t, New(1, 2, 3, 4, 5, 6), it1)
	testIterEq(t, New(1, 2, 3, 4, 5, 6), it2)
}

func TestTee(t *testing.T) {
	its := New(3, 4, 5).Tee(3)
	if len(its) != 3 {
		t.Error("its length wrong")
	}
	for _, it := range its {
		testIter(t, New(3, 4, 5), it)
	}
}
