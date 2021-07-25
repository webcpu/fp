package fp

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
)

func isMap(v reflect.Value) bool {
	return typeQ(v, reflect.Map)
}

func isSlice(v reflect.Value) bool {
	return typeQ(v, reflect.Slice)
}

func typeQ(v reflect.Value, t reflect.Kind) bool {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return v.Kind() == t
	}
	return false
}

func mustBeMap(v reflect.Value) {
	mustBe(v, reflect.Map)
}

func mustBeSlice(v reflect.Value) {
	mustBe(v, reflect.Slice)
}

func mustBe(v reflect.Value, kind reflect.Kind) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		if v.Kind() != kind {
			panic(&reflect.ValueError{details.Name(), v.Kind()})
		}
	}
}

func panicTypeError(v reflect.Value) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		panic(&reflect.ValueError{details.Name(), v.Kind()})
	} else {
		panic(&reflect.ValueError{"type error", v.Kind()})
	}
}

func Map(f func(interface{}) interface{}, slice interface{}) []interface{} {
	sv := reflect.ValueOf(slice)
	mustBeSlice(sv)

	var ys = make([]interface{}, sv.Len())
	//if sv.Len() > 0 {
	//	x := sv.Index(0).Interface()
	//	y := f(x)
	//	//kind := reflect.ValueOf(y).Kind()
	//	t := reflect.TypeOf(y)
	//	zs := reflect.MakeSlice(reflect.SliceOf(t), 0, sv.Len())//.([]interface{})
	//	zs = reflect.Append(zs, reflect.ValueOf(y))
	//	fmt.Println(reflect.TypeOf(zs))
	//	fmt.Println(zs)
	//}
	for i := 0; i < sv.Len(); i++ {
		x := sv.Index(i).Interface()
		//y := f(x)
		ys[i] = f(x)
	}
	return ys
}

var Filter = Select

func Select(f func(interface{}) bool, slice interface{}) []interface{} {
	sv := reflect.ValueOf(slice)
	mustBeSlice(sv)

	var xs = []interface{}{}
	for i := 0; i < sv.Len(); i++ {
		x := sv.Index(i).Interface()
		if f(x) {
			xs = append(xs, x)
		}
	}
	return xs
}

var Reduce = Fold
func Fold(f func(r interface{}, element interface{}) interface{}, initial interface{}, slice interface{}) interface{} {
	sv := reflect.ValueOf(slice)
	mustBeSlice(sv)

	var result = initial
	for i := 0; i < sv.Len(); i++ {
		x := sv.Index(i).Interface()
		result = f(result, x)
	}
	return result
}

func MapIndexed(f func(interface{}, interface{}) interface{}, slice interface{}) []interface{}{
	var xs = Map(Identity, slice)
	var rs []interface{} = make([]interface{}, len(xs))
	for index, x := range(xs) {
		rs[index] = f(x, index)
	}
	return rs
}

func Identity(x interface{}) interface{} {
	return x
}

func Range(nums ...int) []int {
	switch len(nums) {
	case 1:
		return _Range(1, nums[0], 1)
	case 2:
		return _Range(nums[0], nums[1], 1)
	case 3:
		return _Range(nums[0], nums[1], nums[2])
	default:
		msg := fmt.Sprintf("Range: Range called with %v arguments; between 1 and 3 arguments are expected.", len(nums))
		panic(msg)
	}
}

func _Range(imin, imax, step int) []int {
	if step == 0 {
		msg := fmt.Sprintf("Range: Range specification in Range[%v,%v,%v] does not have appropriate bounds.", imin, imax, step)
		panic(msg)
	}
	if imin > imax && step > 0 {
		return []int{}
	}

	isInRange := func(min, max, step, i int) bool {
		if step > 0 {
			return imin + i*step <= imax
		} else {
			return imin + i*step >= imax
		}
	}
	rs := make([]int, (imax-imin)/step + 1)
	for i := 0; isInRange(imax, imax, step, i); i++ {
		rs[i] = imin + i*step
	}
	return rs
}

func Length(slice interface{}) int {
	return reflect.ValueOf(slice).Len()
}

func First(slice interface{}) interface{} {
	if Length(slice) == 0 {
		msg := fmt.Sprintf("First: %v has zero length and no first element.", slice)
		panic(msg)
	}
	sv := reflect.ValueOf(slice)
	mustBeSlice(sv)
	return sv.Index(0).Interface()
}

func Last(slice interface{}) interface{} {
	if Length(slice) == 0 {
		msg := fmt.Sprintf("Last: %v has zero length and no first element.", slice)
		panic(msg)
	}
	sv := reflect.ValueOf(slice)
	mustBeSlice(sv)
	return sv.Index(sv.Len() - 1).Interface()
}

func Take(slice interface{}, n int) interface{} {
	if Length(slice) == 0 || n == 0 {
		msg := fmt.Sprintf("Take: %v has zero length and no first element.", slice)
		panic(msg)
	}
	sv := reflect.ValueOf(slice)
	mustBeSlice(sv)

	if sv.Len() < int(math.Abs(float64(n))) {
		var msg string
		if n > 0 {
			msg = fmt.Sprintf("Take: Cannot take positions 0 through %v", n-1)
		} else {
			msg = fmt.Sprintf("Take: Cannot take positions %v through -1", n)
		}
		panic(msg)
	}

	var ys = []interface{}{}
	start, end := func() (int, int) {
		if n > 0 {
			return 0, n
		} else {
			return sv.Len() + n, sv.Len()
		}
	}()

	for i := start; i < end; i++ {
		x := sv.Index(i).Interface()
		ys = append(ys, x)
	}
	return ys
}

func Drop(slice interface{}, n int) interface{} {
	if Length(slice) < int(math.Abs(float64(n))) {
		msg := fmt.Sprintf("Drop: %v has zero length and no first element.", slice)
		panic(msg)
	}
	sv := reflect.ValueOf(slice)
	mustBeSlice(sv)

	if sv.Len() < int(math.Abs(float64(n))) {
		var msg string
		if n > 0 {
			msg = fmt.Sprintf("Drop: Cannot drop positions 1 through %v", n-1)
		} else {
			msg = fmt.Sprintf("Drop: Cannot drop positions %v through -1", n)
		}
		panic(msg)
	}

	var ys = []interface{}{}
	start, end := func() (int, int) {
		if n > 0 {
			return n, sv.Len()
		} else {
			return 0, sv.Len()+n
		}
	}()

	for i := start; i < end; i++ {
		x := sv.Index(i).Interface()
		ys = append(ys, x)
	}
	return ys
}


func Position(expr interface{}, pattern interface{}) [][]interface{} {
	v := reflect.ValueOf(expr)
	switch v.Kind() {
	case reflect.Slice:
		return positionInSlice(v, pattern)
	case reflect.Map:
		return positionInMap(v, pattern)
	default:
		panicTypeError(v)
	}
	return [][]interface{}{}
}

func positionInSlice(sv reflect.Value, pattern interface{}) [][]interface{} {
	results := [][]interface{}{}
	for i := 0; i < sv.Len(); i++ {
		x := sv.Index(i).Interface()
		if reflect.DeepEqual(x, pattern) {
			results = append(results, []interface{}{i})
		}
	}
	return results
}

func positionInMap(sv reflect.Value, pattern interface{}) [][]interface{} {
	results := [][]interface{}{}
	keys := sv.MapKeys()
	for i := 0; i < sv.Len(); i++ {
		key := keys[i]
		if reflect.DeepEqual(sv.MapIndex(key).Interface(), pattern) {
			results = append(results, []interface{}{key.Interface()})
		}
	}
	return results
}

func Count(expr interface{}, pattern interface{}) int {
	v := reflect.ValueOf(expr)
	switch v.Kind() {
	case reflect.Slice:
		return countInSlice(v, pattern)
	case reflect.Map:
		return countInMap(v, pattern)
	default:
		panicTypeError(v)
	}
	return 0
}

func countInSlice(sv reflect.Value, pattern interface{}) int {
	count := 0
	for i := 0; i < sv.Len(); i++ {
		x := sv.Index(i).Interface()
		if reflect.DeepEqual(x, pattern) {
			count ++
		}
	}
	return count
}

func countInMap(sv reflect.Value, pattern interface{}) int {
	count := 0
	keys := sv.MapKeys()
	for i := 0; i < sv.Len(); i++ {
		key := keys[i]
		if reflect.DeepEqual(sv.MapIndex(key).Interface(), pattern) {
			count++
		}
	}
	return count
}



