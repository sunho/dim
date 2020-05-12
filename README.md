# Introduction

Dim wraps echo to provide the dependecy injection for go web development.

It has been used to devlop the server for [Minda](https://github.com/sdbx/minda), a game published in Steam.

# Features

## Easily configurable service

The service instances are created by a function you implement. The function can take a yaml deserializable struct to configure your services. When you call Dim.Init(path), Dim will read yaml files from the path, unmarshal them as the struct parameter of your function and use them to call your function.


# Examples

## Service Configuration

print.go
```go
package main

import (
	"github.com/sunho/dim"
	"fmt"
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
// Dim will read "print.yaml" and pass it by conf here.
func ProvidePrintService(conf PrintServiceConf) *PrintService {
	return &PrintService{
		test: conf.Test,
	}
}
```

log.go
```go
package main

import (
	"dim"
	"fmt"

	"github.com/labstack/echo"
)

type LogService struct {
    PrintService *PrintService `dim:"on"` // you can specify another service as dependency
    // dim:"on" will trigger Dim to inject
}

func (l *LogService) Log(str string) {
    	// use the injected service
	l.PrintService.Print(str)
}

// creator function
func ProvideLogService() *LogService {
	return &LogService{}
}

```

routes.go
```go
package main

// you will primarily use injected services in the API routes
type LogRoute struct {
	LogService *LogService `dim:"on"` // specify dependency
}

// when you add this route by calling Dim.Route
// it will call this function to register sub-routes
func (l *LogRoute) Register(g *dim.Group) {
	g.GET("/", l.get) // root of this route
}

// handler
func (l *LogRoute) get(e echo.Context) error {
	l.LogService.Log("Hello Dim!") // use injected service
	return e.String(200, "asdf")
}
```

```go

func main() {
    d := dim.New()

    // register service creator functions
    d.Provide(ProvideLogService)
    d.Provide(ProvidePrintService)

    // create service instances
    // read yaml files from config folder 
    err := d.Init("config")
    if err != nil {
    	panic(err)
    }

    // register routes
    d.Register(func(g *dim.Group) {
    	g.Route("/log", &LogRoute{}) // this is the route that Dim will inject dependencies into
    })

    // start http server
    d.Start(":8080")
}
```
