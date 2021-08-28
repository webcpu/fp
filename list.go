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
		if !(v.Kind() == reflect.Array || v.Kind() == reflect.Slice) {
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
	if (fv.Type().NumIn() != len(types)-numOut) || fv.Type().NumOut() != numOut {
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
		go worker(i)
	}
	wg.Wait()
	return ys.Interface()
}

func MapThread(f interface{}, slices ...interface{}) interface{} {
	checkMapThreadArguments(f, slices)
	return _MapThread(f, slices)
}

func _MapThread(f interface{}, slices []interface{}) interface{} {
	fv := reflect.ValueOf(f)
	minLength := reflect.ValueOf(slices[0]).Len()
	for i := 1; i < len(slices); i++ {
		length := reflect.ValueOf(slices[i]).Len()
		if length < minLength {
			minLength = length
		}
	}

	results := reflect.MakeSlice(reflect.SliceOf(fv.Type().Out(0)), minLength, minLength)
	for i := 0; i < minLength; i++ {
		numIn := len(slices)
		ins := make([]reflect.Value, numIn, numIn)
		for j := 0; j < len(slices); j++ {
			ins[j] = reflect.ValueOf(slices[j]).Index(i)
		}
		results.Index(i).Set(fv.Call(ins[:])[0])
	}

	return results.Interface()
}

func checkMapThreadArguments(f interface{}, slices []interface{}) {
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	if len(slices) == 0 {
		msg := fmt.Sprintf("MapThread: MapThread called with %v elements array/slice; 2 elements array/slice is expected.", len(slices))
		panic(msg)
	}
	elementTypes := []reflect.Type{}
	for i := 0; i < len(slices); i++ {
		sv := reflect.ValueOf(slices[i])
		mustBeArraySlice(sv)
		elementType := sv.Type().Elem()
		fmt.Println(sv.Type().Elem().String())
		elementTypes = append(elementTypes, elementType)
	}

	types := append(elementTypes, nil)
	argumentsTypes := ""
	for i := 0; i < len(elementTypes); i++ {
		if i < len(elementTypes)-1 {
			argumentsTypes = argumentsTypes + elementTypes[i].String() + " , "
		} else {
			argumentsTypes = argumentsTypes + elementTypes[i].String()
		}
	}
	if !verifyFuncSignature(fv, 1, types...) {
		msg := "MapThread : function signature must be func(" + argumentsTypes + ") OutputElementType"
		panic(msg)
	}
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
		go worker(i)
	}
	wg.Wait()
}

func Filter(f interface{}, slice interface{}) interface{} {
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

func MapIndexed(f interface{}, slice interface{}) interface{} {
	sv := reflect.ValueOf(slice)
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	mustBeArraySlice(sv)

	elementType := sv.Type().Elem()
	mustBeFuncSignature(sv, fv, 1, elementType, reflect.ValueOf(0).Type(), nil)

	var ins [2]reflect.Value
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
			return imin+i*step <= imax
		} else {
			return imin+i*step >= imax
		}
	}
	rs := make([]int, (imax-imin)/step+1)
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
		ys.Index(i - start).Set(sv.Index(i))
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
			return 0, sv.Len() + n
		}
	}()

	elementType := sv.Type().Elem()
	var ys = reflect.MakeSlice(reflect.SliceOf(elementType), end-start, end-start)

	for i := start; i < end; i++ {
		ys.Index(i - start).Set(sv.Index(i))
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

func Union(lists ...interface{}) interface{} {
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

	keys := makeKeys(lists, elementType)
	for i := 0; i < len(lists); i++ {
		sv := reflect.ValueOf(lists[i])
		keys = appendIfNotInMap(keys, m, sv)
	}
	return keys.Interface()
}

func appendIfNotInMap(keys reflect.Value, m reflect.Value, sv reflect.Value) reflect.Value {
	existed := reflect.ValueOf(true)
	zeroValue := reflect.Value{}
	for j := 0; j < sv.Len(); j++ {
		key := sv.Index(j)
		if m.MapIndex(key) == zeroValue {
			keys = reflect.Append(keys, key)
		}
		m.SetMapIndex(key, existed)
	}
	return keys
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

func DeleteDuplicates(args ...interface{}) interface{} {
	checkDeleteDuplicatesArguments(args)

	if len(args) == 1 {
		return Union(args[0])
	} else {
		return _DeleteDuplicates(args[0], args[1])
	}
}

func checkDeleteDuplicatesArguments(args []interface{}) {
	switch len(args) {
	case 1:
		mustBeArraySlice(reflect.ValueOf(args[0]))
	case 2:
		sv := reflect.ValueOf(args[0])
		elementType := sv.Type().Elem()
		fv := reflect.ValueOf(args[1])
		mustBeFuncSignature(sv, fv, 1, elementType, elementType, reflect.TypeOf(true))
	default:
		msg := fmt.Sprintf("DeleteDuplicates: DeleteDuplicates called with %v arguments; between 1 and 2 arguments are expected.", len(args))
		panic(msg)
	}
}

func _DeleteDuplicates(list interface{}, f interface{}) interface{} {
	sv := reflect.ValueOf(list)
	elementType := sv.Type().Elem()
	keys := reflect.MakeSlice(reflect.SliceOf(elementType), 0, sv.Len())
	fv := reflect.ValueOf(f)
	for i := 0; i < sv.Len(); i++ {
		keys = appendIfNotDuplicate(sv.Index(i), keys, fv)
	}

	return keys.Interface()
}

func appendIfNotDuplicate(elem reflect.Value, keys reflect.Value, fv reflect.Value) reflect.Value {
	key, found := isDuplicate(elem, keys, fv)
	if !found {
		keys = reflect.Append(keys, key)
	}
	return keys
}

func isDuplicate(elem reflect.Value, slice reflect.Value, fv reflect.Value) (reflect.Value, bool) {
	result := elem
	var found = false
	for j := 0; j < slice.Len(); j++ {
		ins := [2]reflect.Value{slice.Index(j), elem}
		r := fv.Call(ins[:])[0].Interface()
		if r == true {
			result = elem
			found = true
			break
		}
	}
	return result, found
}

func Intersection(args ...interface{}) interface{} {
	if len(args) == 0 {
		return []interface{}{}
	}
	checkIntersectionArguments(args)

	fv := reflect.ValueOf(args[len(args)-1])
	if fv.Kind() == reflect.Func {
		return _IntersectionBy(args[:(len(args)-1)], args[len(args)-1])
	} else {
		return _Intersection(args)
	}
}

func _Intersection(lists []interface{}) interface{} {
	var elementType = reflect.ValueOf(lists[0]).Type().Elem()
	mapType := reflect.MapOf(elementType, reflect.TypeOf(true))

	mapv := reflect.MakeMap(mapType)
	keys := reflect.MakeSlice(reflect.SliceOf(elementType), 0, 0)

	keys, mapv = intersectFirstList(lists[0], elementType, mapType, mapv, keys)
	keys, mapv = intersectRestLists(lists[1:len(lists)], elementType, mapType, mapv, keys)
	return keys.Interface()
}

func intersectFirstList(list interface{}, elementType reflect.Type, mapType reflect.Type, mapv reflect.Value, keys reflect.Value) (reflect.Value, reflect.Value) {
	commonKeys := makeKeys([]interface{}{list}, elementType)
	commonMap := reflect.MakeMap(mapType)
	sv := reflect.ValueOf(list)
	for j := 0; j < sv.Len(); j++ {
		key := sv.Index(j)
		value := mapv.MapIndex(key)
		if !value.IsValid() || value.IsZero() {
			commonMap.SetMapIndex(key, reflect.ValueOf(true))
			commonKeys = reflect.Append(commonKeys, key)
		}
	}
	return commonKeys, commonMap
}

func intersectRestLists(lists []interface{}, elementType reflect.Type, mapType reflect.Type, mapv reflect.Value, keys reflect.Value) (reflect.Value, reflect.Value) {
	commonKeys := reflect.MakeSlice(reflect.SliceOf(elementType), 0, 0)
	commonMap := reflect.MakeMap(mapType)
	for i := 0; i < len(lists); i++ {
		commonKeys, commonMap = intersectRestList(lists[i], elementType, mapType, mapv, keys)
	}
	return commonKeys, commonMap
}

func intersectRestList(list interface{}, elementType reflect.Type, mapType reflect.Type, mapv reflect.Value, keys reflect.Value) (reflect.Value, reflect.Value) {
	sv := reflect.ValueOf(list)
	commonKeys := makeKeys([]interface{}{list}, elementType)
	commonMap := reflect.MakeMap(mapType)
	for j := 0; j < sv.Len(); j++ {
		key := sv.Index(j)
		value := mapv.MapIndex(key)
		if value.IsValid() && value.Interface().(bool) {
			v := commonMap.MapIndex(key)
			if !v.IsValid() {
				commonMap.SetMapIndex(key, reflect.ValueOf(true))
				commonKeys = reflect.Append(commonKeys, key)
			}
		}
	}
	return commonKeys, commonMap
}

func _IntersectionBy(lists []interface{}, f interface{}) interface{} {
	var elementType = reflect.ValueOf(lists[0]).Type().Elem()
	fv := reflect.ValueOf(f)
	results := intersectByFirstList(lists[0], elementType, fv)
	results = intersectByRestLists(lists[:(len(lists)-1)], elementType, fv, results)
	return results.Interface()
}

func intersectByFirstList(list interface{}, elementType reflect.Type, fv reflect.Value) reflect.Value {
	sv := reflect.ValueOf(list)
	keys := reflect.MakeSlice(reflect.SliceOf(elementType), 0, 0)
	for j := 0; j < sv.Len(); j++ {
		key := sv.Index(j)
		index, _, found := isDuplicateKey(key, keys, fv)
		if !found {
			keys = reflect.Append(keys, key)
		} else {
			keys.Index(index).Set(key)
		}
	}
	return keys
}

func intersectByRestLists(lists []interface{}, elementType reflect.Type, fv reflect.Value, results reflect.Value) reflect.Value {
	for i := 0; i < len(lists); i++ {
		results = intersectList(lists, elementType, fv, results, i)
	}
	return results
}

func intersectList(lists []interface{}, elementType reflect.Type, fv reflect.Value, results reflect.Value, i int) reflect.Value {
	keys := reflect.MakeSlice(reflect.SliceOf(elementType), 0, 0)
	sv := reflect.ValueOf(lists[i])
	for j := 0; j < sv.Len(); j++ {
		_, key, found1 := isDuplicateKey(sv.Index(j), results, fv)
		if found1 {
			_, _, found2 := isDuplicateKey(key, keys, fv)
			if !found2 {
				keys = reflect.Append(keys, key)
			}
		}
	}
	return keys
}

func isDuplicateKey(key, keys reflect.Value, fv reflect.Value) (int, reflect.Value, bool) {
	result := key
	var found = false
	var index = -1
	for j := 0; j < keys.Len(); j++ {
		ins := [2]reflect.Value{keys.Index(j), key}
		r := fv.Call(ins[:])[0].Interface()
		if r == true {
			result = keys.Index(j)
			found = true
			index = j
			break
		}
	}
	return index, result, found
}

func checkIntersectionArguments(args []interface{}) {
	switch len(args) {
	case 0:
		msg := fmt.Sprintf("Intersection: Intersection called with %v arguments; at least one argument is expected.", len(args))
		panic(msg)
	case 1:
		mustBeArraySlice(reflect.ValueOf(args[0]))
	default:
		var lists []interface{} = getListsArguments(args)
		checkIntersectionListsArguments(lists)
		checkIntersectionFuncArguments(args)
	}
}

func getListsArguments(args []interface{}) []interface{} {
	fv := reflect.ValueOf(args[len(args)-1])
	if fv.Kind() == reflect.Func {
		return args[:(len(args) - 1)]
	} else {
		return args
	}
}

func checkIntersectionFuncArguments(args []interface{}) {
	fv := reflect.ValueOf(args[len(args)-1])
	if fv.Kind() == reflect.Func {
		sv := reflect.ValueOf(args[0])
		elementType := sv.Type().Elem()
		mustBeFuncSignature(sv, fv, 1, elementType, elementType, reflect.TypeOf(true))
	}
}

func checkIntersectionListsArguments(lists []interface{}) reflect.Type {
	var elementType = reflect.ValueOf(lists[0]).Type().Elem()
	for i := 1; i < len(lists); i++ {
		sv := reflect.ValueOf(lists[i])
		if sv.Type().Elem() != elementType {
			msg := fmt.Sprintf("Union: %v's type should as same as %v's type.", lists[i], lists[0])
			panic(msg)
		}
	}
	return elementType
}

func Complement(list1 interface{}, list2 interface{}) interface{} {
	sv1 := reflect.ValueOf(list1)
	sv2 := reflect.ValueOf(list2)
	mustBeArraySlice(sv1)
	mustBeArraySlice(sv2)
	if sv1.Type().Elem() != sv1.Type().Elem() {
		msg := fmt.Sprintf("Complement: %v's type should as same as %v's type.", list1, list2)
		panic(msg)
	}
	elementType := sv1.Type().Elem()

	mapType := reflect.MapOf(elementType, reflect.TypeOf(true))

	mapv1 := reflect.MakeMap(mapType)
	mapv2 := reflect.MakeMap(mapType)

	for i := 0; i < sv1.Len(); i++ {
		mapv1.SetMapIndex(sv1.Index(i), reflect.ValueOf(true))
	}
	for i := 0; i < sv2.Len(); i++ {
		mapv2.SetMapIndex(sv2.Index(i), reflect.ValueOf(true))
	}

	keys := mapv2.MapKeys()
	result := reflect.MakeSlice(reflect.SliceOf(elementType), 0, 0)
	for _, key := range keys {
		v := mapv1.MapIndex(key) //.Convert(elementType))
		fmt.Printf("%v\n", v)

		if !v.IsValid() {
			result = reflect.Append(result, key)
		}
	}
	return result.Interface()
}

func Transpose(lists []interface{}) interface{} {
	if len(lists) == 0 {
		return lists
	}
	checkTransposeArguments(lists)
	return _Transpose(lists)
}

func _Transpose(lists []interface{}) interface{} {
	length := reflect.ValueOf(lists[0]).Len()
	elementType := reflect.TypeOf([]interface{}{})
	results := reflect.MakeSlice(elementType, length, length)
	for i := 0; i < length; i++ {
		list := make([]interface{}, len(lists), len(lists))
		for j := 0; j < len(lists); j++ {
			sv := reflect.ValueOf(lists[j])
			list[j] = sv.Index(i).Interface()
		}
		results.Index(i).Set(reflect.ValueOf(list))
	}
	return results.Interface()
}

func checkTransposeArguments(lists []interface{}) {
	mustBeArraySlice(reflect.ValueOf(lists[0]))
	length := reflect.ValueOf(lists[0]).Len()
	for i := 1; i < len(lists); i++ {
		sv := reflect.ValueOf(lists[i])
		mustBeArraySlice(sv)
		if sv.Len() != length {
			msg := "Transpose: %v can't be transposed. Each list should have the same length."
			panic(msg)
		}
	}
}
