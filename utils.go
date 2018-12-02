package dim

import (
	"errors"
	"reflect"
)

func parseFactory(factory interface{}) (reflect.Type, reflect.Type) {
	typ := reflect.TypeOf(factory)
	if typ.NumIn() > 1 || typ.NumOut() > 2 || typ.NumOut() == 0 {
		panic("Invalid factory function")
	}
	if typ.NumOut() == 2 && typ.Out(1) != reflect.TypeOf(errors.New("")) {
		panic("Invalid factory function")
	}
	if typ.NumIn() == 0 {
		return typ.Out(0), nil
	}
	return typ.Out(0), typ.In(0)
}

func getConfName(serv reflect.Type) string {
	vals := reflect.New(serv).MethodByName("ConfName").Call(nil)
	if len(vals) != 1 {
		panic("Invalid ConfName")
	}
	val := vals[0]
	if val.Kind() != reflect.String {
		panic("Invalid ConfName")
	}
	return val.String()
}

func callInit(serv interface{}) error {
	fn := reflect.ValueOf(serv).MethodByName("Init")
	if fn.IsValid() {
		if fn.Type().NumIn() != 0 || fn.Type().NumOut() != 1 {
			panic("Invalid Init")
		}
		err := fn.Call(nil)[0].Interface()
		if err != nil {
			return err.(error)
		}
	}
	return nil
}

func callFactory(factory, conf interface{}) (interface{}, error) {
	args := make([]reflect.Value, 0, 1)
	if conf != nil {
		args = append(args, reflect.ValueOf(conf))
	}

	vals := reflect.ValueOf(factory).Call(args)
	if len(vals) == 2 {
		err := vals[1].Interface()
		if err != nil {
			return nil, err.(error)
		}
	}

	return vals[0].Interface(), nil
}
