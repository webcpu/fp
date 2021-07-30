package fp

import (
	"fmt"
	"reflect"
	"time"
)

func Timing(f interface{}) (float64, interface{}) {
	fv := reflect.ValueOf(f)
	mustBe(fv, reflect.Func)
	if fv.Type().NumIn() > 0 {
		msg := fmt.Sprintf("Timing: %v should have zero parameters.", f)
		panic(msg)
	}

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
