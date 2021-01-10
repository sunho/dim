package services

import "github.com/sunho/dim"

// Implements dim.Service
type LogService struct {
	PrintService *PrintService `dim:"on"` // will be injected by dim
}

// You can return nil to say this service doesn't use any config file
func (LogService) DefaultConfig() dim.ServiceConfig {
	return nil
}

func (LogService) ConfigName() string {
	return ""
}

func (LogService) Init(conf dim.ServiceConfig) error {
	return nil
}

func (l *LogService) Log(str string) {
	// Use the injected service
	l.PrintService.Print(str)
}
