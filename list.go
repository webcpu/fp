package fp

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"sync"
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

func mustBeFuncSignature(sv reflect.Value, fv reflect.Value, numOut int, types ...reflect.Type) {
	if !verifyFuncSignature(fv, numOut, types...) {
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

func verifyFuncSignature(fv reflect.Value, numOut int, types ...reflect.Type) bool {
	mustBe(fv, reflect.Func)
	if (fv.Type().NumIn() != len(types) - numOut) || fv.Type().NumOut() != numOut {
		return false
	}

	for i := 0; i < len(types)-numOut; i++ {
		if fv.Type().In(i) != types[i] {
			return false
		}
	}

	var outType reflect.Type
	if numOut > 0 {
		outType = types[len(types)-numOut]
	}
	return numOut == 0 || outType == nil || fv.Type().Out(0) == outType
}

func Map(f interface{}, slice interface{}) interface{} {
	sv := reflect.ValueOf(slice)
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	mustBeArraySlice(sv)

	elementType := sv.Type().Elem()
	mustBeFuncSignature(sv, fv, 1, elementType, nil)

	ys := reflect.MakeSlice(reflect.SliceOf(fv.Type().Out(0)), sv.Len(), sv.Len())
	for i := 0; i < sv.Len(); i++ {
		x := []reflect.Value{sv.Index(i)}
		ys.Index(i).Set(fv.Call(x)[0])
	}
	return ys.Interface()
}

func ParallelMap(f interface{}, slice interface{}) interface{} {
	sv := reflect.ValueOf(slice)
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	mustBeArraySlice(sv)

	elementType := sv.Type().Elem()
	mustBeFuncSignature(sv, fv, 1, elementType, nil)

	ys := reflect.MakeSlice(reflect.SliceOf(fv.Type().Out(0)), sv.Len(), sv.Len())

	var wg sync.WaitGroup
	wg.Add(sv.Len())

	worker := func(i int) {
		defer wg.Done()
		x := []reflect.Value{sv.Index(i)}
		value := fv.Call(x)[0]
		ys.Index(i).Set(value)
	}
	for i := 0; i < sv.Len(); i++ {
		worker(i)
	}
	wg.Wait()
	return ys.Interface()
}

func Do(f interface{}, slice interface{}) {
	sv := reflect.ValueOf(slice)
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	mustBeArraySlice(sv)

	elementType := sv.Type().Elem()
	mustBeFuncSignature(sv, fv, 0, elementType)

	worker := func(i int) {
		x := []reflect.Value{sv.Index(i)}
		fv.Call(x)
	}
	for i := 0; i < sv.Len(); i++ {
		worker(i)
	}
}

func ParallelDo(f interface{}, slice interface{}) {
	sv := reflect.ValueOf(slice)
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	mustBeArraySlice(sv)

	elementType := sv.Type().Elem()
	mustBeFuncSignature(sv, fv, 0, elementType)

	var wg sync.WaitGroup
	wg.Add(sv.Len())

	worker := func(i int) {
		defer wg.Done()
		x := []reflect.Value{sv.Index(i)}
		fv.Call(x)
	}
	for i := 0; i < sv.Len(); i++ {
		worker(i)
	}
	wg.Wait()
}

var Filter = Select

func Select(f interface{}, slice interface{}) interface{} {
	sv := reflect.ValueOf(slice)
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	mustBeArraySlice(sv)

	elementType := sv.Type().Elem()
	mustBeFuncSignature(sv, fv, 1, elementType, reflect.ValueOf(true).Type())

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
	mustBeFuncSignature(sv, fv, 1, resultType, elementType, resultType)

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
	mustBeFuncSignature(sv, fv, 1, elementType, reflect.ValueOf(0).Type(), nil)

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

func Most(slice interface{}) interface{} {
	if Length(slice) == 0 {
		msg := fmt.Sprintf("Most: Cannot take Most of expression %v with length zero.", slice)
		panic(msg)
	}
	sv := reflect.ValueOf(slice)
	mustBeArraySlice(sv)
	return Take(slice, sv.Len()-1)
}

func Rest(slice interface{}) interface{} {
	if Length(slice) == 0 {
		msg := fmt.Sprintf("Rest: Cannot take Most of expression %v with length zero.", slice)
		panic(msg)
	}
	sv := reflect.ValueOf(slice)
	mustBeArraySlice(sv)
	return Drop(slice, 1)
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

func MemberQ(slice interface{}, x interface{}) bool {
	sv := reflect.ValueOf(slice)
	mustBeArraySlice(sv)

	elementType := sv.Type().Elem()
	if reflect.ValueOf(x).Type() != elementType && elementType.String() != "interface {}" {
		msg := fmt.Sprintf("MemberQ: %v's type should be %v", x, elementType)
		panic(msg)
	}

	for i := 0; i < sv.Len(); i++ {
		if reflect.DeepEqual(sv.Index(i).Interface(), x) {
			return true
		}
	}
	return false
}

func KeyMemberQ(m interface{}, key interface{}) bool {
	sv := reflect.ValueOf(m)
	mustBe(sv, reflect.Map)

	elementType := sv.Type().Key()
	if reflect.ValueOf(key).Type() != elementType && elementType.String() != "interface {}" {
		msg := fmt.Sprintf("MemberQ: %v's type should be %v", key, elementType)
		panic(msg)
	}

	value := sv.MapIndex(reflect.ValueOf(key))
	return value.IsValid() && !value.IsZero()
}

func Keys(m interface{}) interface{} {
	sv := reflect.ValueOf(m)
	mustBe(sv, reflect.Map)

	elementType := sv.Type().Key()

	kvs := sv.MapKeys()
	xs := reflect.MakeSlice(reflect.SliceOf(elementType), len(kvs), len(kvs))
	for i := 0; i < sv.Len(); i++ {
		xs.Index(i).Set(kvs[i])
	}
	return xs.Interface()
}

func Values(m interface{}) interface{} {
	sv := reflect.ValueOf(m)
	mustBe(sv, reflect.Map)

	elementType := sv.Type().Elem()

	kvs := sv.MapKeys()
	xs := reflect.MakeSlice(reflect.SliceOf(elementType), len(kvs), len(kvs))
	for i := 0; i < sv.Len(); i++ {
		xs.Index(i).Set(sv.MapIndex(kvs[i]))
	}
	return xs.Interface()
}

func Union(lists... interface{}) interface{} {
	if len(lists) == 0 {
		return []interface{}{}
	}
	checkUnionArguments(lists)
	return _Union(lists)
}

func checkUnionArguments(lists []interface{}) {
	for i := 0; i < len(lists); i++ {
		sv := reflect.ValueOf(lists[i])
		mustBeArraySlice(sv)
	}

	var elementType = reflect.ValueOf(lists[0]).Type().Elem()
	for i := 1; i < len(lists); i++ {
		sv := reflect.ValueOf(lists[i])
		if sv.Type().Elem() != elementType {
			msg := fmt.Sprintf("Union: %v's type should as same as %v's type.", lists[i], lists[0])
			panic(msg)
		}
	}
}

func _Union(lists []interface{}) interface{} {
	var elementType = reflect.ValueOf(lists[0]).Type().Elem()
	mapType := reflect.MapOf(elementType, reflect.TypeOf(true))
	m := reflect.MakeMap(mapType)

	existed := reflect.ValueOf(true)
	zeroValue := reflect.Value{}
	keys := makeKeys(lists, elementType)
	for i := 0; i < len(lists); i++ {
		sv := reflect.ValueOf(lists[i])
		for j := 0; j < sv.Len(); j++ {
			key := sv.Index(j)
			if m.MapIndex(key) == zeroValue {
				keys = reflect.Append(keys, key)
			}
			m.SetMapIndex(key, existed)
		}
	}
	return keys.Interface()
}

func makeKeys(lists []interface{}, elementType reflect.Type) reflect.Value {
	var capacity int = 0
	for i := 0; i < len(lists); i++ {
		sv := reflect.ValueOf(lists[i])
		capacity += sv.Len()
	}
	keys := reflect.MakeSlice(reflect.SliceOf(elementType), 0, capacity)
	return keys
}

func DeleteDuplicates(args... interface{}) interface{} {
	if len(args) > 2 {
		msg := fmt.Sprintf("DeleteDuplicates: DeleteDuplicates called with %v arguments; between 1 and 3 arguments are expected.", len(args))
		panic(msg)
	}

	sv := reflect.ValueOf(args[0])
	elementType := sv.Type().Elem()
	for i := 0; i < len(args); i++ {
		if i == 0 {
			mustBeArraySlice(reflect.ValueOf(args[0]))
		}
		if i == 1 {
			fv := reflect.ValueOf(args[1])
			mustBeFuncSignature(sv, fv, 1, elementType, elementType, reflect.TypeOf(true))
		}
	}

	if len(args) == 1 {
		return Union(args[0])
	}

	keys := reflect.MakeSlice(reflect.SliceOf(elementType), 0, sv.Len())
	fv := reflect.ValueOf(args[1])
	for i := 0; i < sv.Len(); i++ {
		key := sv.Index(i)
		var ins [2]reflect.Value
		ins[1] = key
		var found = false
		for j := 0; j < keys.Len(); j++ {
			ins[0] = keys.Index(j)
			r := fv.Call(ins[:])[0].Interface()
			if r == true {
				found = true
				break
			}
		}
		if !found {
			keys = reflect.Append(keys, key)
		}
	}

	return keys.Interface()
}