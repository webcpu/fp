package fp

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"time"
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

func mustBeArraySlice(v reflect.Value) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		if !(v.Kind() == reflect.Array || v.Kind() == reflect.Slice)  {
			panic(&reflect.ValueError{details.Name(), v.Kind()})
		}
	}
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

func mustBeFuncSignature(sv reflect.Value, fv reflect.Value, types ...reflect.Type) {
	if !verifyFuncSignature(fv, types...) {
		msg := "Map : function signature must be func(" + sv.Type().Elem().String() + ") OutputElementType"
		panic(msg)
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

func verifyFuncSignature(fv reflect.Value, types ...reflect.Type) bool {
	mustBe(fv, reflect.Func)
	if (fv.Type().NumIn() != len(types) - 1) || fv.Type().NumOut() != 1 {
		return false
	}

	for i := 0; i < len(types)-1; i++ {
		if fv.Type().In(i) != types[i] {
			return false
		}
	}

	outType := types[len(types)-1]
	return outType == nil || fv.Type().Out(0) == outType
}

func Map(f interface{}, slice interface{}) interface{} {
	sv := reflect.ValueOf(slice)
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	mustBeArraySlice(sv)

	elementType := sv.Type().Elem()
	mustBeFuncSignature(sv, fv, elementType, nil)

	ys := reflect.MakeSlice(reflect.SliceOf(fv.Type().Out(0)), sv.Len(), sv.Len())
	for i := 0; i < sv.Len(); i++ {
		x := []reflect.Value{sv.Index(i)}
		ys.Index(i).Set(fv.Call(x)[0])
	}
	return ys.Interface()
}

var Filter = Select

func Select(f interface{}, slice interface{}) interface{} {
	sv := reflect.ValueOf(slice)
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	mustBeArraySlice(sv)

	elementType := sv.Type().Elem()
	mustBeFuncSignature(sv, fv, elementType, reflect.ValueOf(true).Type())

	ys := []interface{}{}
	for i := 0; i < sv.Len(); i++ {
		x := sv.Index(i)
		args := []reflect.Value{x}
		if fv.Call(args)[0].Interface().(bool) {
			ys = append(ys, x.Interface())
		}
	}

	zs := reflect.MakeSlice(reflect.SliceOf(elementType), len(ys), len(ys))
	for i := 0; i < zs.Len(); i++ {
		zs.Index(i).Set(reflect.ValueOf(ys[i]))
	}
	return zs.Interface()
}

var Reduce = Fold
func Fold(f interface{}, initial interface{}, slice interface{}) interface{} {
	sv := reflect.ValueOf(slice)
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	mustBeArraySlice(sv)

	elementType := sv.Type().Elem()
	resultType := reflect.ValueOf(initial).Type()
	mustBeFuncSignature(sv, fv, resultType, elementType, resultType)

	var result = reflect.ValueOf(initial)
	var ins [2]reflect.Value
	ins[0] = result
	for i := 0; i < sv.Len(); i++ {
		ins[1] = sv.Index(i)
		result = fv.Call(ins[:])[0]
		ins[0] = result
	}
	return result.Interface()
}

func MapIndexed(f interface{}, slice interface{}) interface{}{
	sv := reflect.ValueOf(slice)
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	mustBeArraySlice(sv)

	elementType := sv.Type().Elem()
	mustBeFuncSignature(sv, fv, elementType, reflect.ValueOf(0).Type(), nil)

	var ins[2]reflect.Value
	ys := reflect.MakeSlice(reflect.SliceOf(fv.Type().Out(0)), sv.Len(), sv.Len())
	for i := 0; i < sv.Len(); i++ {
		ins[0] = sv.Index(i)
		ins[1] = reflect.ValueOf(i)
		ys.Index(i).Set(fv.Call(ins[:])[0])
	}
	return ys.Interface()
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
	mustBeArraySlice(sv)
	return sv.Index(0).Interface()
}

func Last(slice interface{}) interface{} {
	if Length(slice) == 0 {
		msg := fmt.Sprintf("Last: %v has zero length and no first element.", slice)
		panic(msg)
	}
	sv := reflect.ValueOf(slice)
	mustBeArraySlice(sv)
	return sv.Index(sv.Len() - 1).Interface()
}

func Take(slice interface{}, n int) interface{} {
	if Length(slice) == 0 || n == 0 {
		msg := fmt.Sprintf("Take: %v has zero length and no first element.", slice)
		panic(msg)
	}
	sv := reflect.ValueOf(slice)
	mustBeArraySlice(sv)

	if sv.Len() < int(math.Abs(float64(n))) {
		var msg string
		if n > 0 {
			msg = fmt.Sprintf("Take: Cannot take positions 0 through %v", n-1)
		} else {
			msg = fmt.Sprintf("Take: Cannot take positions %v through -1", n)
		}
		panic(msg)
	}

	start, end := func() (int, int) {
		if n > 0 {
			return 0, n
		} else {
			return sv.Len() + n, sv.Len()
		}
	}()

	elementType := sv.Type().Elem()
	var ys = reflect.MakeSlice(reflect.SliceOf(elementType), end-start, end-start)

	for i := start; i < end; i++ {
		ys.Index(i-start).Set(sv.Index(i))
	}
	return ys.Interface()
}

func Drop(slice interface{}, n int) interface{} {
	if Length(slice) < int(math.Abs(float64(n))) {
		msg := fmt.Sprintf("Drop: %v has zero length and no first element.", slice)
		panic(msg)
	}
	sv := reflect.ValueOf(slice)
	mustBeArraySlice(sv)

	if sv.Len() < int(math.Abs(float64(n))) {
		var msg string
		if n > 0 {
			msg = fmt.Sprintf("Drop: Cannot drop positions 1 through %v", n-1)
		} else {
			msg = fmt.Sprintf("Drop: Cannot drop positions %v through -1", n)
		}
		panic(msg)
	}

	start, end := func() (int, int) {
		if n > 0 {
			return n, sv.Len()
		} else {
			return 0, sv.Len()+n
		}
	}()

	elementType := sv.Type().Elem()
	var ys = reflect.MakeSlice(reflect.SliceOf(elementType), end-start, end-start)

	for i := start; i < end; i++ {
		ys.Index(i-start).Set(sv.Index(i))
	}
	return ys.Interface()
}

func Position(expr interface{}, pattern interface{}) [][]interface{} {
	v := reflect.ValueOf(expr)
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
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
	case reflect.Slice, reflect.Array:
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
			count++
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

func Reverse(expr interface{}) []interface{} {
	v := reflect.ValueOf(expr)
	mustBeArraySlice(v)

	var xs = make([]interface{}, v.Len())
	for i, j := 0, len(xs)-1; i <= j; i, j = i+1, j-1 {
		xs[i], xs[j] = v.Index(j).Interface(), v.Index(i).Interface()
	}
	return xs
}

func Max(args ...interface{}) interface{} {
	if len(args) == 0 {
		msg := fmt.Sprintf("Max: no enough arguments.")
		panic(msg)
	}

	if len(args) == 1 {
		v := reflect.ValueOf(args[0])
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			return MaxInSlice(args[0])
		} else {
			return args[0]
		}
	}
	return MaxInSlice(args)
}

func MaxInSlice(expr interface{}) interface{} {
	v := reflect.ValueOf(expr)
	if v.Len() == 0 {
		msg := fmt.Sprintf("Max: %v has zero length and no first element.", expr)
		panic(msg)
	}
	var r interface{} = v.Index(0).Interface()
	for i := 1; i < v.Len(); i++ {
		x := v.Index(i).Interface()
		if Greater(x, r) {
			r = x
		}
	}
	return r
}

func Less(a interface{}, b interface{}) bool {
	return !Greater(a, b)
}

func isPrimitiveComparable(x interface{}) bool {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Int:
		return true
	case reflect.Int8:
		return true
	case reflect.Int16:
		return true
	case reflect.Int32:
		return true
	case reflect.Int64:
		return true
	case reflect.Uint:
		return true
	case reflect.Uint8:
		return true
	case reflect.Uint16:
		return true
	case reflect.Uint32:
		return true
	case reflect.Uint64:
		return true
	case reflect.Float32:
		return true
	case reflect.Float64:
		return true
	case reflect.String:
		return true
	default:
		return false
	}
}

func Greater(a interface{}, b interface{}) bool {
	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)
	if v1.Kind() != v2.Kind() {
		msg := fmt.Sprintf("%v and %v are different kinds of elments", v1, v2)
		panic(msg)
	}
	switch v1.Kind() {
	case reflect.Int:
		return a.(int) > b.(int)
	case reflect.Int8:
		return a.(int8) > b.(int8)
	case reflect.Int16:
		return a.(int16) > b.(int16)
	case reflect.Int32:
		return a.(int32) > b.(int32)
	case reflect.Int64:
		return a.(int64) > b.(int64)
	case reflect.Uint:
		return a.(uint) > b.(uint)
	case reflect.Uint8:
		return a.(uint8) > b.(uint8)
	case reflect.Uint16:
		return a.(uint16) > b.(uint16)
	case reflect.Uint32:
		return a.(uint32) > b.(uint32)
	case reflect.Uint64:
		return a.(uint64) > b.(uint64)
	case reflect.Float32:
		return a.(float32) > b.(float32)
	case reflect.Float64:
		return a.(float64) > b.(float64)
	case reflect.String:
		return a.(string) > b.(string)
	default:
		msg := fmt.Sprintf("compare function is missing, you must use Sort(xs, func(a interface{}, b interface{})bool{}) to sort.")
		panic(msg)
	}
}


func Min(args ...interface{}) interface{} {
	if len(args) == 0 {
		msg := fmt.Sprintf("Max: no enough arguments.")
		panic(msg)
	}

	if len(args) == 1 {
		v := reflect.ValueOf(args[0])
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			return MinInArraySlice(args[0])
		} else {
			return args[0]
		}
	}
	return MinInArraySlice(args)
}

func MinInArraySlice(expr interface{}) interface{} {
	v := reflect.ValueOf(expr)
	if v.Len() == 0 {
		msg := fmt.Sprintf("Min: %v has zero length and no first element.", expr)
		panic(msg)
	}
	var r interface{} = v.Index(0).Interface()
	for i := 1; i < v.Len(); i++ {
		x := v.Index(i).Interface()
		if !Greater(x, r) {
			r = x
		}
	}
	return r
}

//func Timing(f interface{}) (float64, interface{}) {
func Timing(f interface{}) (float64, interface{}) {
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)

	var ins = []reflect.Value{}
	start := time.Now()
	out := fv.Call(ins[:])
	elapsed := time.Since(start)
	var r interface{}
	if len(out) > 0 {
		r = out[0].Interface()
	}
	return float64(elapsed) / float64(1000000000), r
}