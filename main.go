package main

import (
	"mockconcat/goconcat"
)

func main() {
	files := []string{"test.go", "test2.go"}
	goconcat.Goconcat(files)
}
