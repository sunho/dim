package main

import (
	"dim"
	"fmt"

	"github.com/labstack/echo"
)

type PrintService struct {
}

func (p *PrintService) Print(str string) {
	fmt.Println(str)
}

type LogService struct {
	PrintService *PrintService `dim:"on"`
}

func (l *LogService) Log(str string) {
	l.PrintService.Print(str)
}

func providePrintService() *PrintService {
	return &PrintService{}
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
	err := d.Init("")
	if err != nil {
		panic(err)
	}
	d.Register(func(g *dim.Group) {
		g.Route("/log", &LogRoute{})
	})
	d.Start(":8080")
}
