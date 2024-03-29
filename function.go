package fp

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
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

func Bind(fs ...interface{}) func(...interface{}) (interface{}, error) {
	checkBindArguments(fs)
	return func(args ...interface{}) (interface{}, error) {
		var values []reflect.Value = make([]reflect.Value, len(args))
		for i := 0; i < len(args); i++ {
			values[i] = reflect.ValueOf(args[i])
		}
		for _, f := range fs {
			fv := reflect.ValueOf(f)
			outs := fv.Call(values)
			values = outs[:(len(outs) - 1)]
			err, ok := outs[len(outs)-1].Interface().(error)
			if ok && err != nil {
				return nil, err
			}
		}
		outputs := make([]interface{}, len(values)-1)
		for i := 0; i < len(outputs); i++ {
			outputs[i] = values[i].Interface()
		}
		return values[0].Interface(), nil
	}
}

func checkBindArguments(fs []interface{}) {
	for i := 0; i < len(fs); i++ {
		fv := reflect.ValueOf(fs[i])
		mustBe(fv, reflect.Func)
		if i > 0 {
			checkFunctionArgumentNumber(fs, i, fv)
		}

		if i == len(fs)-1 {
			checkFunctionOutputArguments(fv)
		}
	}
}

func checkFunctionOutputArguments(fv reflect.Value) {
	if fv.Type().NumOut() != 2 {
		panic("the output arguments of the last function should be (interface{}, error)")
	}
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	t := fv.Type().Out(1)
	if !(t.Kind() == reflect.Interface && t.Implements(errorInterface)) {
		panic("the output arguments of the last function should be (interface{}, error)")
	}
}

func checkFunctionArgumentNumber(fs []interface{}, i int, fv reflect.Value) {
	prev := reflect.ValueOf(fs[i-1])
	if prev.Type().NumOut()-1 != fv.Type().NumIn() {
		msg := fmt.Sprintf("argument #%v and argument %v doesn't match.", signature(fs[i-1]), signature(fs[i]))
		panic(msg)
	}
}

func signature(f interface{}) string {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return "<not a function>"
	}

	buf := strings.Builder{}
	buf.WriteString("func (")
	for i := 0; i < t.NumIn(); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(t.In(i).String())
	}
	buf.WriteString(")")
	if numOut := t.NumOut(); numOut > 0 {
		if numOut > 1 {
			buf.WriteString(" (")
		} else {
			buf.WriteString(" ")
		}
		for i := 0; i < t.NumOut(); i++ {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(t.Out(i).String())
		}
		if numOut > 1 {
			buf.WriteString(")")
		}
	}

	return buf.String()
}
