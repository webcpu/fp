package fp

import (
	"fmt"
	"reflect"
)

func Apply(f interface{}, expr interface{}) interface{} {
	sv := reflect.ValueOf(expr)
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	if !(sv.Kind() == reflect.Array || sv.Kind() == reflect.Slice) {
		return expr
	}
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
	return fv.Call(values)[0].Interface()
}
