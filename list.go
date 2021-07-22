package fp

import (
	"fmt"
	"reflect"
	"runtime"
)

func mustBeSlice(v reflect.Value) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		if v.Kind() != reflect.Slice {
			panic(&reflect.ValueError{details.Name(), v.Kind()})
		}
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

func Filter(f func(interface{}) bool, slice interface{}) []interface{} {
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

func Reduce(initial interface{}, f func(r interface{}, element interface{}) interface{}, slice interface{}) interface{} {
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

func Length(slice interface{}) interface{} {
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

