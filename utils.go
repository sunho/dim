package dim

import (
	"reflect"

	"github.com/labstack/echo"
)

func indirectType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func parseFactory(factory interface{}) (reflect.Type, reflect.Type) {
	typ := reflect.TypeOf(factory)
	if typ.NumIn() > 1 || typ.NumOut() > 2 || typ.NumOut() == 0 {
		panic("Invalid factory function")
	}
	if typ.NumOut() == 2 && typ.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		panic("Invalid factory function")
	}
	if typ.NumIn() == 0 {
		return indirectType(typ.Out(0)), nil
	}
	return indirectType(typ.Out(0)), indirectType(typ.In(0))
}

func getDefaultConf(conf reflect.Type) (interface{}, bool) {
	fn := reflect.New(conf).Elem().MethodByName("Default")
	if fn.IsValid() {
		if fn.Type().NumIn() != 0 || fn.Type().NumOut() != 1 {
			panic(conf.Name() + "'s Default should be a void function that returns a struct")
		}
		vals := fn.Call(nil)
		val := vals[0]
		if indirectType(val.Type()) != conf {
			panic(conf.Name() + "'s Default should return a " + conf.Name())
		}
		return ptr(reflect.Indirect(val)).Interface(), true
	}
	return nil, false
}

func ptr(v reflect.Value) reflect.Value {
	pt := reflect.PtrTo(v.Type())
	pv := reflect.New(pt.Elem())
	pv.Elem().Set(v)
	return pv
}

func getConfName(serv reflect.Type) (string, bool) {
	fn := reflect.New(serv).Elem().MethodByName("ConfigName")
	if fn.IsValid() {
		if fn.Type().NumIn() != 0 || fn.Type().NumOut() != 1 {
			panic(serv.Name() + "'s ConfigName should be a void function that returns a string")
		}
		vals := fn.Call(nil)
		val := vals[0]
		if val.Kind() != reflect.String {
			panic(serv.Name() + "'s ConfigName should return a string value")
		}
		return val.String(), true
	}
	return "", false
}

func callValidate(conf interface{}) bool {
	fn := reflect.ValueOf(conf).MethodByName("Validate")
	if fn.IsValid() {
		if fn.Type().NumIn() != 0 || fn.Type().NumOut() != 1 {
			panic(reflect.TypeOf(conf).Name() + "'s Validate should be a void function that returns a bool")
		}
		vals := fn.Call(nil)
		val := vals[0]
		if val.Kind() != reflect.Bool {
			panic(reflect.TypeOf(conf).Name() + "'s Validate should return a bool value")
		}
		return val.Bool()
	}
	return true
}

func callInit(serv interface{}) error {
	fn := reflect.ValueOf(serv).MethodByName("Init")
	if fn.IsValid() {
		if fn.Type().NumIn() != 0 || fn.Type().NumOut() != 1 {
			panic(reflect.TypeOf(serv).Name() + "'s Init should be a void function that returns a error")
		}
		err := fn.Call(nil)[0].Interface()
		if err == nil {
			return nil
		}
		out, ok := err.(error)
		if !ok {
			panic(reflect.TypeOf(serv).Name() + "'s Init should should return a error value")
		}
		return out
	}
	return nil
}

func callFactory(factory, conf interface{}) (interface{}, error) {
	fn := reflect.ValueOf(factory)
	args := make([]reflect.Value, 0, 1)
	if conf != nil {
		val := reflect.ValueOf(conf)
		if fn.Type().In(0).Kind() == reflect.Ptr {
			args = append(args, val)
		} else {
			args = append(args, val.Elem())
		}
	}

	vals := fn.Call(args)
	if len(vals) == 2 {
		err := vals[1].Interface()
		if err != nil {
			return nil, err.(error)
		}
	}

	return ptr(reflect.Indirect(vals[0])).Interface(), nil
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

func callHandler(handler, context interface{}) error {
	fn := reflect.ValueOf(handler)
	t := fn.Type()

	c := reflect.ValueOf(context)
	t2 := c.Type()
	if t.In(0) != t2 {
		panic("Invalid context type: " + t2.Name() + " Expected: " + t.In(0).Name())
	}
	out := fn.Call([]reflect.Value{c})
	return out[0].Interface().(error)
}

func convertHandler(handler interface{}) echo.HandlerFunc {
	if !isHandler(handler) {
		panic("Invalid handler")
	}
	return func(c echo.Context) error {
		return callHandler(handler, c)
	}
}
