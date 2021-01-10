package services

import (
	"fmt"

	"github.com/sunho/dim"
)

// Definition of config for PrintService
type PrintServiceConf struct {
	AppName string `yaml:"app_name"`
}

// Implements dim.Service
type PrintService struct {
	appname string
}

// This will be run by dim
func (p *PrintService) Init(conf dim.ServiceConfig) error {
	c := conf.(*PrintServiceConf)
	p.appname = c.AppName
	return nil
}

// Name of config file for this service
func (PrintService) ConfigName() string {
	return "printservice"
}

// Default configuration that will be used to generate initial config file
func (PrintService) DefaultConfig() dim.ServiceConfig {
	return &PrintServiceConf{
		AppName: "dim example",
	}
}

func (p *PrintService) Print(str string) {
	fmt.Println("[" + p.appname + "] " + str + " (this is printed by services.PrintService)")
}
