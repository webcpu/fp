package fp

import (
	"fmt"
	"reflect"
)

func Sort(args ...interface{}) []interface{} {
	switch len(args) {
	case 0:
		msg := fmt.Sprintf("Sort: no enough arguments.")
		panic(msg)
	case 1:
		v := reflect.ValueOf(args[0])
		mustBeArraySlice(v)
		if v.Len() > 0 {
			if !isPrimitiveComparable(v.Index(0).Interface()) {
				msg := fmt.Sprintf("compare function is missing, you must use Sort(xs, less) but not Sort(xs) to sort, function less is a ordering function.")
				panic(msg)
			}
		}
		return _Sort(args[0], Less)
	case 2:
		var less func(interface{}, interface{}) bool = args[1].(func(interface{}, interface{})bool)
		return _Sort(args[0], less)
	default:
		msg := fmt.Sprintf("Sort: too many arguments.")
		panic(msg)
	}
}

func _Sort(expr interface{}, less func(interface{},interface{}) bool) []interface{} {
	v := reflect.ValueOf(expr)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return sortArraySlice(expr, less)
	default:
		panicTypeError(v)
	}
	return []interface{}{}
}

// An implementation of Interface can be sorted by the routines in this package.
// The methods refer to elements of the underlying collection by integer index.
type SortInterface interface {
	// Len is the number of elements in the collection.
	Len() int

	// Less reports whether the element with index i
	// must sort before the element with index j.
	//
	// If both Less(i, j) and Less(j, i) are false,
	// then the elements at index i and j are considered equal.
	// Sort may place equal elements in any order in the final result,
	// while Stable preserves the original input order of equal elements.
	//
	// Less must describe a transitive ordering:
	//  - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
	//  - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
	//
	// Note that floating-point comparison (the < operator on float32 or float64 values)
	// is not a transitive ordering when not-a-number (NaN) values are involved.
	// See Float64Slice.Less for a correct implementation for floating-point values.

	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}

type XSlice []interface{}
func (x XSlice) Len() int           { return len(x) }
func (x XSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func sortArraySlice(expr interface{}, _less func(interface{}, interface{}) bool) []interface{} {
	var slice XSlice = Map(Identity, expr).(XSlice)
	less := func(i,j int) bool { return _less(slice[i], slice[j]) }
	_sortSlice(slice, less)
	return []interface{}{}
	//return Map(Identity, slice)
}

// Sort sorts data.
// It makes one call to data.Len to determine n and O(n*log(n)) calls to
// data.Less and data.Swap. The sort is not guaranteed to be stable.
func _sortSlice(data SortInterface, less func(int,int)bool) {
	n := data.Len()
	quickSort(data, 0, n, maxDepth(n), less)
}

// maxDepth returns a threshold at which quicksort should switch
// to heapsort. It returns 2*ceil(lg(n+1)).
func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

func quickSort(data SortInterface, a, b, maxDepth int, less func(int,int)bool) {
	for b-a > 12 { // Use ShellSort for slices <= 12 elements
		if maxDepth == 0 {
			heapSort(data, a, b, less)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(data, a, b, less)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort(data, a, mlo, maxDepth, less)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSort(data, mhi, b, maxDepth, less)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		// Do ShellSort pass with gap 6
		// It could be written in this simplified form cause b-a <= 12
		for i := a + 6; i < b; i++ {
			if less(i, i-6) {
				data.Swap(i, i-6)
			}
		}
		insertionSort(data, a, b, less)
	}
}

// insertionSort sorts data[a:b] using insertion sort.
func insertionSort(data SortInterface, a, b int, less func(int, int) bool) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}

// siftDown implements the heap property on data[lo:hi].
// first is an offset into the array where the root of the heap lies.
func siftDown(data SortInterface, lo, hi, first int, less func(int,int)bool) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && less(first+child, first+child+1) {
			child++
		}
		if !less(first+root, first+child) {
			return
		}
		data.Swap(first+root, first+child)
		root = child
	}
}

func heapSort(data SortInterface, a, b int, less func(int,int)bool) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first,less)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)
		siftDown(data, lo, i, first,less)
	}
}

// Quicksort, loosely following Bentley and McIlroy,
// ``Engineering a Sort Function,'' SP&E November 1993.

// medianOfThree moves the median of the three values data[m0], data[m1], data[m2] into data[m1].
func medianOfThree(data SortInterface, m1, m0, m2 int, less func(int,int)bool) {
	// sort 3 elements
	if less(m1, m0) {
		data.Swap(m1, m0)
	}
	// data[m0] <= data[m1]
	if less(m2, m1) {
		data.Swap(m2, m1)
		// data[m0] <= data[m2] && data[m1] < data[m2]
		if less(m1, m0) {
			data.Swap(m1, m0)
		}
	}
	// now data[m0] <= data[m1] <= data[m2]
}

func swapRange(data SortInterface, a, b, n int) {
	for i := 0; i < n; i++ {
		data.Swap(a+i, b+i)
	}
}

func doPivot(data SortInterface, lo, hi int, less func(int,int)bool) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1) // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree(data, lo, lo+s, lo+2*s, less)
		medianOfThree(data, m, m-s, m+s, less)
		medianOfThree(data, hi-1, hi-1-s, hi-1-2*s, less)
	}
	medianOfThree(data, lo, m, hi-1, less)

	// Invariants are:
	//	data[lo] = pivot (set up by ChoosePivot)
	//	data[lo < i < a] < pivot
	//	data[a <= i < b] <= pivot
	//	data[b <= i < c] unexamined
	//	data[c <= i < hi-1] > pivot
	//	data[hi-1] >= pivot
	pivot := lo
	a, c := lo+1, hi-1

	for ; a < c && less(a, pivot); a++ {
	}
	b := a
	for {
		for ; b < c && !less(pivot, b); b++ { // data[b] <= pivot
		}
		for ; b < c && less(pivot, c-1); c-- { // data[c-1] > pivot
		}
		if b >= c {
			break
		}
		// data[b] > pivot; data[c-1] <= pivot
		data.Swap(b, c-1)
		b++
		c--
	}
	// If hi-c<3 then there are duplicates (by property of median of nine).
	// Let's be a bit more conservative, and set border to 5.
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		// Lets test some points for equality to pivot
		dups := 0
		if !less(pivot, hi-1) { // data[hi-1] = pivot
			data.Swap(c, hi-1)
			c++
			dups++
		}
		if !less(b-1, pivot) { // data[b-1] = pivot
			b--
			dups++
		}
		// m-lo = (hi-lo)/2 > 6
		// b-lo > (hi-lo)*3/4-1 > 8
		// ==> m < b ==> data[m] <= pivot
		if !less(m, pivot) { // data[m] = pivot
			data.Swap(m, b-1)
			b--
			dups++
		}
		// if at least 2 points are equal to pivot, assume skewed distribution
		protect = dups > 1
	}
	if protect {
		// Protect against a lot of duplicates
		// Add invariant:
		//	data[a <= i < b] unexamined
		//	data[b <= i < c] = pivot
		for {
			for ; a < b && !less(b-1, pivot); b-- { // data[b] == pivot
			}
			for ; a < b && less(a, pivot); a++ { // data[a] < pivot
			}
			if a >= b {
				break
			}
			// data[a] == pivot; data[b-1] < pivot
			data.Swap(a, b-1)
			a++
			b--
		}
	}
	// Swap pivot into middle
	data.Swap(pivot, b-1)
	return b - 1, c
}