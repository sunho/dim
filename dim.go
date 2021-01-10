package dim

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/labstack/echo"
	"gopkg.in/yaml.v2"
)

// Dim instance. It wraps an echo instance.
type Dim struct {
	E           *echo.Echo
	servs       map[reflect.Type]Service
	confs       map[string]interface{}
	visitedConf map[string]bool
	multiple    bool
}

// ServiceConfig is an YAML configuration struct for each service
type ServiceConfig interface{}

// Service is an interface for a service that dim will inject to other services, route groups, or middlewares.
type Service interface {
	// Initialization function that will be called inside Dim.Init.
	Init(conf ServiceConfig) error
	// Default yaml config struct that will be used to generate an initial config file. Return nil to disable configuration for this service.
	DefaultConfig() ServiceConfig
	// Name of config file for this service.
	ConfigName() string
}

// New returns a new instance of dim.
func New() *Dim {
	return &Dim{
		E:           echo.New(),
		servs:       map[reflect.Type]Service{},
		visitedConf: map[string]bool{},
	}
}

// Init will do two things.
// 1. Load configuration files from specified path. If they don't exist, dim will generate initial config files based on Service.DefaultConfig()
// 2. Call Init of service instances provided by dim.Provide. When Init of some service is called, its service fields are guranteed to be already injected with initialized service instances.
// When the single argument is true, a single yaml file from specified path (e.g. config.yaml) will be used for service configuration.
// Otherwise, multiple yaml files inside the directory specified by path (e.g. config) will be used for service configuration.
func (d *Dim) Init(path string, single bool) error {
	d.multiple = !single
	err := d.initialConfig(path)
	if err != nil {
		return err
	}

	for _, serv := range d.servs {
		d.Inject(serv)
	}

	servidx := make(map[int]reflect.Type)
	idxserv := make(map[reflect.Type]int)
	t := 0
	for serv := range d.servs {
		servidx[t] = serv
		idxserv[serv] = t
		t++
	}
	g := newGraph(t)
	for serv := range d.servs {
		connectEdge(serv, idxserv, g)
	}
	idx, err := g.TopologicalSort()
	if err != nil {
		return err
	}

	for _, i := range idx {
		serv := d.servs[servidx[i]]
		conf := serv.DefaultConfig()
		if conf != nil {
			err := d.readConfig(path, serv.ConfigName(), conf)
			if err != nil {
				return err
			}
		}
		err := d.servs[servidx[i]].Init(conf)
		if err != nil {
			return err
		}
	}
	return nil
}

// Register registers the root route handlers by calling specified register function
func (d *Dim) Register(register RegisterFunc) {
	t := newGroup(d, d.E.Group(""))
	register(t)
}

// Provide registers serivce instances that are to be injected
func (d *Dim) Provide(services ...Service) {
	for _, serv := range services {
		typ := indirectType(reflect.TypeOf(serv))
		if _, ok := d.servs[typ]; ok {
			panic("Duplicate factory func")
		}
		d.servs[typ] = serv
	}
}

// Start starts a web server listening to specifed addr
func (d *Dim) Start(addr string) error {
	return d.E.Start(addr)
}

// Inject can be used to inject service instances to any custom struct containing service fields with dim:"on"
func (d *Dim) Inject(bean interface{}) error {
	typ := reflect.TypeOf(bean).Elem()

	toinject := make(map[int]reflect.Type)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		serv := indirectType(field.Type)
		tag := parseTag(field)
		if tag.on {
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

func (d *Dim) readConfig(path, name string, conf ServiceConfig) error {
	if _, ok := d.visitedConf[name]; ok {
		return errors.New("Duplicated config name: " + name)
	}
	d.visitedConf[name] = true
	if d.multiple {
		return d.readConfigMultipleFiles(path, name, conf)
	}
	return d.readConfigSingleFile(path, name, conf)
}

func (d *Dim) readConfigSingleFile(path, name string, conf ServiceConfig) error {
	if d.confs == nil {
		d.confs = make(map[string]interface{})
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(buf, &d.confs)
		if err != nil {
			return err
		}
	}

	conf, ok := d.confs[name]
	if !ok {
		return errors.New(name + " not configured")
	}
	buf, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dim) readConfigMultipleFiles(path, name string, conf ServiceConfig) error {
	path2 := filepath.Join(path, name+".yaml")
	buf, err := ioutil.ReadFile(path2)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dim) initialConfig(path string) error {
	if _, err := os.Stat(filepath.Join(path)); !os.IsNotExist(err) && !d.multiple {
		return nil
	}

	if d.multiple {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.MkdirAll(path, 0755)
			if err != nil {
				return err
			}
		}
	}

	confs := make(map[string]interface{})
	for _, serv := range d.servs {
		c := serv.DefaultConfig()
		if c == nil {
			continue
		}
		if d.multiple {
			path2 := filepath.Join(path, serv.ConfigName()+".yaml")
			if _, err := os.Stat(path2); !os.IsNotExist(err) {
				continue
			}
			buf, err := yaml.Marshal(c)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(path2, buf, 0755)
			if err != nil {
				return err
			}
			log.Println(path2 + " made with the default configuration")
		} else {
			confs[serv.ConfigName()] = c
		}
	}
	if !d.multiple {
		buf, err := yaml.Marshal(confs)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, buf, 0755)
		if err != nil {
			return err
		}
		log.Println(path + " made with the default configuration")
	}
	return nil
}

type dimTag struct {
	on     bool
	noinit bool
}

func parseTag(field reflect.StructField) dimTag {
	splitted := strings.Split(field.Tag.Get("dim"), ",")
	out := dimTag{}
	for _, item := range splitted {
		if item == "on" {
			out.on = true
		} else if item == "noinit" {
			out.noinit = true
		}
	}
	return out
}

func connectEdge(serv reflect.Type, servidx map[reflect.Type]int, g *graph) {
	for i := 0; i < serv.NumField(); i++ {
		field := serv.Field(i)
		serv2 := indirectType(field.Type)
		tag := parseTag(field)
		if tag.on && !tag.noinit {
			g.Connect(servidx[serv], servidx[serv2])
		}
	}
}
