package main

import (
	"dim"
	"fmt"
)

type Service2 struct {
	S    *Service1 `dim:"on"`
	asdf string
}

func (s *Service2) Hoi() {
	s.asdf = "1234"
	fmt.Println("asdfsaf")
}

type Service1 struct {
	S *Service2 `dim:"on"`
}

func (s *Service1) Init() error {
	s.S.Hoi()
	return nil
}

func provide1() *Service1 {
	return &Service1{}
}

func provide2() *Service2 {
	return &Service2{}
}

func main() {
	d := dim.New()
	d.Provide(provide1)
	d.Provide(provide2)
	d.Init("")
}
