package dim

import (
	"io/ioutil"
	"path/filepath"
	"reflect"

	"github.com/labstack/echo"
	yaml "gopkg.in/yaml.v2"
)

type Dim struct {
	e         *echo.Echo
	servs     map[reflect.Type]interface{}
	factories map[reflect.Type]interface{}
}

func New() *Dim {
	return &Dim{
		e:         echo.New(),
		servs:     map[reflect.Type]interface{}{},
		factories: map[reflect.Type]interface{}{},
	}
}

func (d *Dim) Init(path string) error {
	servs := map[reflect.Type]interface{}{}
	for _, factory := range d.factories {
		serv, conf := parseFactory(factory)
		if conf != nil {
			name := getConfName(serv)
			buf, err := ioutil.ReadFile(filepath.Join(path, name+".yaml"))
			if err != nil {
				return err
			}

			c := reflect.New(conf).Interface()
			err = yaml.Unmarshal(buf, c)
			if err != nil {
				return err
			}

			s, err := callFactory(factory, c)
			if err != nil {
				return err
			}
			servs[serv] = s
		} else {
			s, err := callFactory(factory, nil)
			if err != nil {
				return err
			}
			servs[serv] = s
		}
	}
	d.servs = servs
	for _, serv := range servs {
		d.inject(serv)
	}

	for _, serv := range servs {
		err := callInit(serv)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Dim) Register(register RegisterFunc) {
	t := newGroup(d, d.e.Group(""))
	register(t)
}

func (d *Dim) Provide(factory interface{}) {
	serv, _ := parseFactory(factory)
	if _, ok := d.factories[serv]; ok {
		panic("Duplicate factory func")
	}
	d.factories[serv] = factory
}

func (d *Dim) Start(addr string) error {
	return d.e.Start(addr)
}

func (d *Dim) inject(bean interface{}) {
	typ := reflect.TypeOf(bean).Elem()

	toinject := map[int]reflect.Type{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Tag.Get("dim") == "on" {
			toinject[i] = field.Type
		}
	}

	val := reflect.ValueOf(bean).Elem()
	for val2, serv := range toinject {
		serv2, ok := d.servs[serv]
		if !ok {
			panic("Dependency unresolved")
		}
		val.Field(val2).Set(reflect.ValueOf(serv2))
	}
}
