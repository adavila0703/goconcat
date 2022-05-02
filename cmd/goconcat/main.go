package main

import (
	"fmt"
	"log"
	"mockconcat/goconcat"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	app = &cli.App{
		Name:  "goconcat",
		Usage: "Concat go files.",
		Flags: []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:    "mockery",
				Usage:   `Using the mockery flag will automatically look for all files which have a "mock_" prefix and sort them to a mockery package.`,
				Action:  woop,
				Aliases: []string{"-m"},
			},
		},
	}
)

func woop(c *cli.Context) error {
	fmt.Println("hello")
	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	err := goconcat.GoConcat()
	if err != nil {
		log.Fatal(err)
	}
}
