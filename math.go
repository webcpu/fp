package fp

import (
	"errors"
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
var Power = Pow
func Pow( x interface{}, y interface{}) interface{} {
	if isIntegerFloat(x) && isIntegerFloat(y) {
		return powForNumber(x, y)
	} else if isComplex(x) {
		return powerForComplex(x, y)
	} else {
		msg := fmt.Sprintf("Pow: Unsupported type of %v", x)
		panic(msg)
	}
}

func isComplex(x interface{}) bool {
	v := reflect.ValueOf(x)
	return 	v.Kind() == reflect.Complex64 || v.Kind() == reflect.Complex128
}

func normalize(x interface{}, y interface{}) (interface{}, interface{})  {
	if !(isIntegerFloat(y) || isComplex(y)) {
		msg := fmt.Sprintf("Pow: %v and %v should be both complex64 or complex128", x, y)
		panic(msg)
	}
	vx := reflect.ValueOf(x)
	if vx.Kind() == reflect.Complex64 {
		return normalizeComplex64(x, y)
	}
	if vx.Kind() == reflect.Complex128 {
		return normalizeComplex128(x, y)
	}
	return nil, nil
}

func normalizeComplex128(x interface{}, y interface{}) (interface{}, interface{}) {
	vy := reflect.ValueOf(y)
	if isIntegerFloat(y) {
		r, _ := toFloat64(y)
		y1 := complex(r, 0)
		return x, y1
	} else if isComplex(y) {
		if vy.Kind() == reflect.Complex64 {
			y1 := complex(float32(real(y.(complex64))), float32(imag(y.(complex64))))
			return x, y1
		} else {
			return x, y
		}
	}
	return nil, nil
}

func normalizeComplex64(x interface{}, y interface{}) (interface{}, interface{}) {
	vy := reflect.ValueOf(y)
	if isIntegerFloat(y) {
		r, _ := toFloat32(y)
		y1 := complex(r, 0)
		return x.(complex64), y1
	} else if isComplex(y) {
		if vy.Kind() == reflect.Complex64 {
			return x, y
		} else {
			y1 := complex(float32(real(y.(complex128))), float32(imag(y.(complex128))))
			return x, y1
		}
	}
	return nil, nil
}

func powerForComplex(x interface{}, y interface{}) interface{} {
	x, y = normalize(x, y)
	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Complex64 {
		var x1 complex128 = complex128(complex(real(x.(complex64)), imag(x.(complex64))))
		var y1 complex128 = complex128(complex(real(y.(complex64)), imag(y.(complex64))))
		var c complex128 = cmplx.Pow(x1, y1)
		return complex(float32(real(c)), float32(imag(c)))
	}
	if v.Kind() == reflect.Complex128 {
		return cmplx.Pow(x.(complex128), y.(complex128))
	}

	msg := fmt.Sprintf("Pow: Unsupported type of %v", x)
	panic(msg)
}

func powForNumber(x interface{}, y interface{}) interface{} {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Int:
		return int(_Pow(x, y))
	case reflect.Int8:
		return int8(_Pow(x, y))
	case reflect.Int16:
		return int16(_Pow(x, y))
	case reflect.Int32:
		return int32(_Pow(x, y))
	case reflect.Int64:
		return int64(_Pow(x, y))
	case reflect.Uint:
		return uint(_Pow(x, y))
	case reflect.Uint8:
		return uint8(_Pow(x, y))
	case reflect.Uint16:
		return uint16(_Pow(x, y))
	case reflect.Uint32:
		return uint32(_Pow(x, y))
	case reflect.Uint64:
		return uint64(_Pow(x, y))
	case reflect.Float32:
		return float32(_Pow(x, y))
	case reflect.Float64:
		return _Pow(x, y)
	default:
		panic("Should not happend")
	}
}

func _Pow( x interface{}, y interface{}) float64 {
	x1, _ := toFloat64(x)
	y1, _ := toFloat64(y)
	return math.Pow(x1, y1)
}

func toFloat64(x interface{}) (float64, error) {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Int:
		return float64(x.(int)), nil
	case reflect.Int8:
		return float64(x.(int8)), nil
	case reflect.Int16:
		return float64(x.(int16)), nil
	case reflect.Int32:
		return float64(x.(int32)), nil
	case reflect.Int64:
		return float64(x.(int64)), nil
	case reflect.Uint:
		return float64(x.(uint)), nil
	case reflect.Uint8:
		return float64(x.(uint8)), nil
	case reflect.Uint16:
		return float64(x.(uint16)), nil
	case reflect.Uint32:
		return float64(x.(uint32)), nil
	case reflect.Uint64:
		return float64(x.(uint64)), nil
	case reflect.Float32:
		return float64(x.(float32)), nil
	case reflect.Float64:
		return x.(float64), nil
	default:
		return float64(0), errors.New("wrong type")
	}
}

func toFloat32(x interface{}) (float32, error) {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Int:
		return float32(x.(int)), nil
	case reflect.Int8:
		return float32(x.(int8)), nil
	case reflect.Int16:
		return float32(x.(int16)), nil
	case reflect.Int32:
		return float32(x.(int32)), nil
	case reflect.Int64:
		return float32(x.(int64)), nil
	case reflect.Uint:
		return float32(x.(uint)), nil
	case reflect.Uint8:
		return float32(x.(uint8)), nil
	case reflect.Uint16:
		return float32(x.(uint16)), nil
	case reflect.Uint32:
		return float32(x.(uint32)), nil
	case reflect.Uint64:
		return float32(x.(uint64)), nil
	case reflect.Float32:
		return float32(x.(float32)), nil
	case reflect.Float64:
		return x.(float32), nil
	default:
		return float32(0), errors.New("wrong type")
	}
}

func isIntegerFloat(x interface{}) bool {
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
	default:
		return false
	}
}
