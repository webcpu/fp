package fp

import (
	"fmt"
	"reflect"
	"runtime"
)

func Apply(f interface{}, expr interface{}) interface{} {
	sv := reflect.ValueOf(expr)
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	//if !(sv.Kind() == reflect.Array || sv.Kind() == reflect.Slice) {
	//	return expr
	//}
	mustBeArraySlice(sv)

	if !fv.Type().IsVariadic() {
		for i := 0; i < fv.Type().NumIn(); i++ {
			left := fv.Type().In(i)
			right := reflect.ValueOf(sv.Index(i).Interface()).Type()
			if left.String() != "interface {}" && left != right {
				msg := fmt.Sprintf("Apply: arguments[%v]'s type should be %v but not %v.", i, left, right)
				panic(msg)
			}
		}
	}

	values := make([]reflect.Value, sv.Len(), sv.Len())//reflect.MakeSlice(elementType, sv.Len(), sv.Len())
	for i := 0; i < len(values); i++ {
		values[i] = reflect.ValueOf(sv.Index(i).Interface())
	}

	return returnResult(fv, values, sv)
}

func Construct(f interface{}, args... interface{}) interface{} {
	sv := reflect.ValueOf(args)
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)

	if !fv.Type().IsVariadic() {
		for i := 0; i < fv.Type().NumIn(); i++ {
			left := fv.Type().In(i)
			right := reflect.ValueOf(sv.Index(i).Interface()).Type()
			if left.String() != "interface {}" && left != right {
				msg := fmt.Sprintf("Construct: arguments[%v]'s type should be %v but not %v.", i, left, right)
				panic(msg)
			}
		}
	}

	values := make([]reflect.Value, sv.Len(), sv.Len())//reflect.MakeSlice(elementType, sv.Len(), sv.Len())
	for i := 0; i < len(values); i++ {
		values[i] = reflect.ValueOf(sv.Index(i).Interface())
	}

	return returnResult(fv, values, sv)
}

func returnResult(fv reflect.Value, values []reflect.Value, sv reflect.Value) interface{} {
	result := fv.Call(values)
	switch len(result) {
	case 0:
		return nil
	case 1:
		return result[0].Interface()
	default:
		output := reflect.MakeSlice(reflect.SliceOf(sv.Type().Elem()), len(result), len(result))
		for i := 0; i < len(result); i++ {
			output.Index(i).Set(result[i])
		}
		return output.Interface()
	}
}

func Composition(fs... interface{}) func(...interface{}) interface{} {
	if len(fs) == 0 {
		return func (args... interface{}) interface{} {
			if len(args) == 1 {
				return Identity(args[0])
			} else {
				msg := "Identity: Identity called with #{len(args)} arguments; 1 argument is expected."
				panic(msg)
			}
		}
	}

	for i := 0; i < len(fs); i++ {
		f := fs[i]
		fv := reflect.ValueOf(f)
		if  fv.Kind() != reflect.Func {
			msg := fmt.Sprintf("Composition: %v is not a function.", getFunctionName(f))
			panic(msg)
		}
	}

	fn := func(args... interface{}) interface{} {
		sv := reflect.ValueOf(args)
		var values []reflect.Value = make([]reflect.Value, len(args), len(args))
		for i := 0; i < len(values); i++ {
			values[i] = reflect.ValueOf(args[i])
		}
		for i := len(fs) - 1; i >= 0; i-- {
			f := fs[i]
			fv := reflect.ValueOf(f)
			if i > 0 {
				values = fv.Call(values)
			}
		}

		return returnResult(reflect.ValueOf(fs[0]), values, sv)
	}
	return fn
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}