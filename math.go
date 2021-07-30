package fp

import (
	"fmt"
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
