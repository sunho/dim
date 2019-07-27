package dim

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
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

func (d *Dim) readConfig(path, name string, typ reflect.Type) (interface{}, error) {
	return d.readConfigMultipleFiles(path, name, typ)
}

func (d *Dim) readConfigMultipleFiles(path, name string, typ reflect.Type) (interface{}, error) {
	path2 := filepath.Join(path, name+".yaml")
	if _, err := os.Stat(path2); os.IsNotExist(err) {
		c, support := getDefaultConf(typ)
		if !support {
			return nil, errors.New(path2 + " doesn't exist")
		}
		buf, err := yaml.Marshal(c)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(path2, buf, 0644)
		if err != nil {
			return nil, err
		}
		log.Println(path2 + " made with the default configuration")
		return c, nil
	}

	buf, err := ioutil.ReadFile(path2)
	if err != nil {
		return nil, err
	}

	c := reflect.New(typ).Interface()
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (d *Dim) Init(path string) error {
	servs := map[reflect.Type]interface{}{}
	for _, factory := range d.factories {
		serv, conf := parseFactory(factory)
		name, ok := getConfName(serv)
		if ok {
			c, err := d.readConfig(path, name, conf)
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

	servidx := make(map[int]reflect.Type)
	idxserv := make(map[reflect.Type]int)
	t := 0
	for serv := range servs {
		servidx[t] = serv
		idxserv[serv] = t
		t++
	}
	g := newGraph(t)
	for serv := range servs {
		d.connect(serv, idxserv, g)
	}
	idx, err := g.TopologicalSort()
	if err != nil {
		return err
	}
	log.Println(servidx)
	log.Println(idx)

	for _, i := range idx {
		err := callInit(servs[servidx[i]])
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

// func(conf Struct) *Serv
// func(conf Struct) (*Serv, error)
// func (s *Serv) ConfigName() string
// func (s *Serv) Init() error
func (d *Dim) Provide(factories ...interface{}) {
	for _, factory := range factories {
		serv, _ := parseFactory(factory)
		if _, ok := d.factories[serv]; ok {
			panic("Duplicate factory func")
		}
		d.factories[serv] = factory
	}
}

func (d *Dim) Start(addr string) error {
	return d.e.Start(addr)
}

func (d *Dim) connect(serv reflect.Type, servidx map[reflect.Type]int, g *graph) {
	for i := 0; i < serv.NumField(); i++ {
		field := serv.Field(i)
		serv2 := indirectType(field.Type)
		if field.Tag.Get("dim") == "on" {
			g.Connect(servidx[serv], servidx[serv2])
		}
	}
}

func (d *Dim) inject(bean interface{}) error {
	typ := reflect.TypeOf(bean).Elem()

	toinject := make(map[int]reflect.Type)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		serv := indirectType(field.Type)
		if field.Tag.Get("dim") == "on" {
			if field.Type.Kind() != reflect.Ptr {
				return errors.New(typ.Name() + "." + field.Name + " should be a pointer")
			}
			toinject[i] = serv
		} else if _, ok := d.servs[serv]; ok && field.Tag.Get("dim") != "-" {
			log.Println(typ.Name() + "." + field.Name + " can be injected but doesn't have dim tag. If you don't want to see this message set the dim tag of it to '-'")
		}
	}

	val := reflect.ValueOf(bean).Elem()
	for i, serv := range toinject {
		serv2, ok := d.servs[serv]
		if !ok {
			return errors.New(serv.Name() + " not resolved")
		}
		val.Field(i).Set(reflect.ValueOf(serv2))
	}
	return nil
}
