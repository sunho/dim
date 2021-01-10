package dim

import (
	"reflect"
)

func indirectType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func ptr(v reflect.Value) reflect.Value {
	pt := reflect.PtrTo(v.Type())
	pv := reflect.New(pt.Elem())
	pv.Elem().Set(v)
	return pv
}

func isHandler(handler interface{}) bool {
	fn := reflect.ValueOf(handler)
	if !fn.IsValid() {
		return false
	}
	t := fn.Type()
	if t.Kind() != reflect.Func {
		return false
	}
	if t.NumIn() != 1 || t.NumOut() != 1 {
		return false
	}
	if t.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return false
	}
	return true
}
