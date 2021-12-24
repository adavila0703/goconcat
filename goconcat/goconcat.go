package goconcat

import (
	"io/ioutil"
	"mockconcat/flags"
)

func Goconcat() {
	ioutil.ReadDir(*flags.Path)
}
