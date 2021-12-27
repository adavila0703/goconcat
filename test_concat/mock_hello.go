package test_concat

import "fmt"

var (
	TestVar = 3
)

const (
	TestConst = "Test"
)

type IPerson interface {
	SayHello()
}

var _ IPerson = (*Person)(nil)

type Person struct {
	Name string
	Age  int
}

func (p *Person) SayHello() {
	fmt.Println("hello")
}

func test2() {

}
