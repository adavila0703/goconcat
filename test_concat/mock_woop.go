package test_concat

import "fmt"

var (
	Number = 3
)

const (
	Another = "Test"
)

type IDog interface {
	Bark()
}

var _ IDog = (*Dog)(nil)

type Dog struct {
	Name string
	Age  int
}

func (p *Dog) Bark() {
	fmt.Println("hello")
}

func test() {

}
