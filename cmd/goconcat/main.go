package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	app = &cli.App{
		Name:      "goconcat",
		UsageText: "goconcat can be used to concatenate multiple go files.",
		Flags:     []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:    "mockery",
				Usage:   `Using the mockery flag will automatically look for all files which have a "mock_" prefix and sort them to a mockery package.`,
				Action:  mockery,
				Aliases: []string{"-m"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "tes",
						Usage:    "test",
						Value:    "hi",
						Required: true,
					},
				},
				ArgsUsage: "[test] [anothertest]",
			},
		},
	}
)

func mockery(c *cli.Context) error {
	// hello := c.Args().Get(0)

	test := &flag.FlagSet{}
	for _, value := range c.Command.Flags {
		fmt.Println("name", value.Names())
		value.Apply(test)
		fmt.Println(test)
	}

	fmt.Println()

	// if err := goconcat.GoConcat(
	// 	".",
	// 	[]utils.Directory{
	// 		utils.DirectoryGit,
	// 	},
	// 	[]utils.FileType{
	// 		utils.FileGo,
	// 	},
	// 	[]utils.PrefixType{
	// 		utils.PrefixMock,
	// 	},
	// ); err != nil {
	// 	log.Fatal(err)
	// }

	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
