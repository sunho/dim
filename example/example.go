package main

import (
	"fmt"

	"github.com/sunho/dim"

	"github.com/labstack/echo"
)

// config struct
type PrintServiceConf struct {
	Test string `yaml:"test"`
}

type PrintService struct {
	LogService *LogService `dim:"on"`
	test       string
}

// provide config file name
func (PrintService) ConfigName() string {
	return "print"
}

func (p *PrintService) Print(str string) {
	fmt.Println(p.test, str)
}

// creator function
// conf will be provided by Dim
func providePrintService(conf PrintServiceConf) (*PrintService, error) {
	return &PrintService{
		test: conf.Test,
	}, nil
}

type LogService struct {
	PrintService *PrintService `dim:"on"` // will be injected by Dim
	// dim:"on" will trigger Dim to inject
}

func (l *LogService) Log(str string) {
	// use the injected service
	l.PrintService.Print(str)
}

type LogMiddleware struct {
	LogService *LogService `dim:"on"`
}

func (l *LogMiddleware) Act(c echo.Context) error {
	l.LogService.Log("Middleware")
	return nil
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
	l.LogService.Log("Route")
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
		g.Use(&LogMiddleware{})
		g.Route("/log", &LogRoute{})
	})

	// start http server
	d.Start(":8080")
}
