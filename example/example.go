package main

import (
	"dim"
	"fmt"

	"github.com/labstack/echo"
)

type PrintServiceConf struct {
	Test string `yaml:"test"`
}

type PrintService struct {
	test string
}

func (PrintService) ConfigName() string {
	return "print"
}

func (p *PrintService) Print(str string) {
	fmt.Println(p.test)
}

type LogService struct {
	PrintService *PrintService `dim:"on"`
}

func (l *LogService) Log(str string) {
	l.PrintService.Print(str)
}

func providePrintService(conf PrintServiceConf) *PrintService {
	return &PrintService{
		test: conf.Test,
	}
}

func provideLogService() *LogService {
	return &LogService{}
}

type LogRoute struct {
	LogService *LogService `dim:"on"`
}

func (l *LogRoute) Register(g *dim.Group) {
	g.GET("/", l.post)
}

func (l *LogRoute) post(e echo.Context) error {
	l.LogService.Log("Hello Dim!")
	return e.String(200, "asdf")
}

func main() {
	d := dim.New()
	d.Provide(provideLogService)
	d.Provide(providePrintService)
	err := d.Init("config")
	if err != nil {
		panic(err)
	}
	d.Register(func(g *dim.Group) {
		g.Route("/log", &LogRoute{})
	})
	d.Start(":8080")
}
