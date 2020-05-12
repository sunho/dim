# Introduction

Dim wraps echo to provide the dependecy injection for go web development.

It has been used to devlop the server for (Minda)[https://github.com/sdbx/minda], a game published in Steam.

# Features

## Easily configurable service

The service instances are created by a function you implement. The function can take a yaml deserializable struct to configure your services. When you call Dim.Init(path), Dim will read yaml files from the path, unmarshal them as the struct parameter of your function and use them to call your function.


# Examples

## Service Configuration

```go
package main

import (
	"dim"
	"fmt"

	"github.com/labstack/echo"
)

// config struct
type PrintServiceConf struct {
	Test string `yaml:"test"`
}

type PrintService struct {
	test string
}

// provide config file name
func (PrintService) ConfigName() string {
	return "print"
}

func (p *PrintService) Print(str string) {
	fmt.Println(p.test)
}

// creator function
// conf will be provided by Dim
func providePrintService(conf PrintServiceConf) *PrintService {
	return &PrintService{
		test: conf.Test,
	}
}

type LogService struct {
    PrintService *PrintService `dim:"on"` // will be injected by Dim
    // dim:"on" will trigger Dim to inject
}

func (l *LogService) Log(str string) {
    // use the injected service
	l.PrintService.Print(str)
}

// creator function
func provideLogService() *LogService {
	return &LogService{}
}

type LogRoute struct {
	LogService *LogService `dim:"on"`
}

// register routes
func (l *LogRoute) Register(g *dim.Group) {
	g.GET("/", l.get)
}

// handler
func (l *LogRoute) get(e echo.Context) error {
	l.LogService.Log("Hello Dim!")
	return e.String(200, "asdf")
}

func main() {
    d := dim.New()

    // register service creator functions
	d.Provide(provideLogService)
    d.Provide(providePrintService)

    // create service instances
    // unmarshal yaml files from config folder 
    // and provide them to creator functions
    err := d.Init("config")
	if err != nil {
		panic(err)
    }

    // register routes
	d.Register(func(g *dim.Group) {
		g.Route("/log", &LogRoute{})
    })

    // start http server
	d.Start(":8080")
}
```
