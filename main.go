package main

import (
	"log"
	"mockconcat/goconcat"
)

func main() {
	err := goconcat.Goconcat()
	if err != nil {
		log.Fatal(err)
	}
}
