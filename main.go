package main

import (
	"log"
	"mockconcat/goconcat"
)

func main() {
	err := goconcat.GoConcat()
	if err != nil {
		log.Fatal(err)
	}
}
