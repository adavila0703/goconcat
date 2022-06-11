# Goconcat

[![](https://godoc.org/github.com/adavila0703/goconcat?status.svg)](http://godoc.org/github.com/adavila0703/goconcat)
[![Code Coverage](https://codecov.io/gh/adavila0703/goconcat/branch/main/graph/badge.svg)](https://app.codecov.io/gh/adavila0703/goconcat)
[![Go Report Card](https://goreportcard.com/badge/github.com/adavila0703/goconcat)](https://goreportcard.com/report/github.com/adavila0703/goconcat)

> Goconcat was originally developed as a tool to consolidate mock files which were generated from mockery. Since mockery's update to this issue, Goconcat has been converted to a general purpose file consolidation tool.

![goconcat](./images/goconcat.png)

## Installation

```shell
    go get github.com/adavila0703/goconcat
```

## CLI tool

```shell
    go install github.com/adavila0703/goconcat@latest
```

## Usage

If you deicide to use a JSON file for your options, follow optionsExample.json

JSON example

```json
{
  "rootPath": ".",
  "ignoredDirectories": ["dir1", "dir2"],
  "filePrefix": ["mocks_", "mock_"],
  "destination": "newdir",
  "deleteOldFiles": true,
  "concatPkg": false
}
```

```go
package main

import "github.com/adavila0703/goconcat"

func main() {
    options := goconcat.NewOptions()
    options.SetJSONOptions("options.json")
    goconcat.GoConcat(options)
}
```

Alternatively, if you decide to not use JSON for options, you can set you options using SetOptions() method.

```go
package main

import "github.com/adavila0703/goconcat"

func main() {
	options := NewOptions()
	options.SetOptions(
		".",
		nil,
		[]PrefixType{"test_"},
		".",
		true,
		false,
		false,
		[]FileType{FileGo},
	)
    goconcat.GoConcat(options)
}
```

You can also go around using options and concatenate files by file paths.

```go
package main

import (
	"log"

	"github.com/adavila0703/goconcat"
)

func main() {
	filePaths := []string{"test/file_one.go", "test/file_two.go"}

	files, fileSet, err := goconcat.ParseASTFiles(filePaths)
	if err != nil {
		log.Fatal(err)
	}

	file, err := goconcat.ConcatFiles(files, fileSet)
	if err != nil {
		log.Fatal(err)
	}

	newFilePath := "test/file_one_two.go"

	goconcat.WriteASTFile(file, fileSet, newFilePath)
}
```
