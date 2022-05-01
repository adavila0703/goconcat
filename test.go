package test_concat

import (
	"fmt"
	"go/token"
)

var (
	TestVar = 3

	_ IDog = (*Dog)(nil)
)

const (
	TestConst = "Test"
)

type IPerson interface {
	SayHello()
}

var (
	_ IPerson = (*Person)(nil)
	_ IDog    = (*Dog)(nil)
)

type Person struct {
	Name string
	Age  int
}

func (p *Person) SayHello() {
	fmt.Println("hello")
}

func test2() {
	fmt.Println(token.ADD)
}
