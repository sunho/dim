package main

import (
	"github.com/sunho/dim"
	"github.com/sunho/dim/example/services"

	"github.com/labstack/echo"
)

type LogMiddleware struct {
	LogService *services.LogService `dim:"on"` // will be injected by dim
}

func (l *LogMiddleware) Act(c echo.Context) error {
	l.LogService.Log("Log from middleware")
	return nil
}

type HelloWorldRoute struct {
	LogService *services.LogService `dim:"on"` // will be injected by dim
}

func (l *HelloWorldRoute) Register(g *dim.Group) {
	g.GET("/", l.get)
}

func (l *HelloWorldRoute) get(e echo.Context) error {
	l.LogService.Log("Log from hello world route")
	return e.String(200, "Hello World")
}

func main() {
	d := dim.New()

	// Register service instances
	d.Provide(&services.PrintService{}, &services.LogService{})

	// Load yaml files from config folder
	// and call Init of each service with config struct
	// Default config files will be generated if they don't exist (based on DefaultConfig of each service)
	err := d.Init("config", false)
	// if you want a single config file instead of configuration directory, try d.Init("config.yaml", true)
	if err != nil {
		panic(err)
	}

	// Register routes
	d.Register(func(g *dim.Group) {
		g.Use(&LogMiddleware{})
		g.Route("", &HelloWorldRoute{})
	})

	// Start http server
	d.Start(":8080")
}
