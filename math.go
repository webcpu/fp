package fp

import (
	"fmt"
	"math"
	"math/cmplx"
	"reflect"
)

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

func Abs( x interface{}) interface{} {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Int:
		return int(math.Abs(float64(x.(int))))
	case reflect.Int8:
		return int8(math.Abs(float64(x.(int8))))
	case reflect.Int16:
		return int16(math.Abs(float64(x.(int16))))
	case reflect.Int32:
		return int32(math.Abs(float64(x.(int32))))
	case reflect.Int64:
		return int64(math.Abs(float64(x.(int64))))
	case reflect.Uint:
		return x.(uint)
	case reflect.Uint8:
		return x.(uint8)
	case reflect.Uint16:
		return x.(uint16)
	case reflect.Uint32:
		return x.(uint32)
	case reflect.Uint64:
		return x.(uint64)
	case reflect.Float32:
		return float64(math.Abs(float64(x.(float32))))
	case reflect.Float64:
		return math.Abs(x.(float64))
	case reflect.Complex64:
		var y complex64 = x.(complex64)
		return math.Hypot(float64(real(y)), float64(imag(y)))
	case reflect.Complex128:
		return cmplx.Abs(x.(complex128))
	default:
		msg := fmt.Sprintf("Abs: Unsupported type of %v", x)
		panic(msg)
	}
}
