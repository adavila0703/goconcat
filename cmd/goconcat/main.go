package main

import (
	"log"
	"os"

	"github.com/adavila0703/goconcat/internal/utils"
	"github.com/adavila0703/goconcat/pkg/concat"

	"github.com/pkg/errors"
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
				Aliases: []string{"m"},
			},
		},
	}
)

func mockery(c *cli.Context) error {
	options := utils.NewOptions(true, true)

	filesToDelete, err := concat.GetFilePaths(
		".",
		[]utils.Directory{
			utils.DirectoryGit,
		},
		[]utils.FileType{
			utils.FileGo,
		},
		[]utils.PrefixType{
			utils.PrefixGoconcat,
		},
	)
	if err != nil {
		return errors.WithStack(err)
	}

	concat.DeleteFiles(filesToDelete)

	filePaths, err := concat.GetFilePaths(
		".",
		[]utils.Directory{
			utils.DirectoryGit,
		},
		[]utils.FileType{
			utils.FileGo,
		},
		[]utils.PrefixType{
			utils.PrefixMockery,
		},
	)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := concat.GoConcat(
		".",
		[]utils.Directory{
			utils.DirectoryGit,
		},
		[]utils.FileType{
			utils.FileGo,
		},
		[]utils.PrefixType{
			utils.PrefixMockery,
		},
		utils.DestinationMockery,
		options,
	); err != nil {
		log.Fatal(err)
	}

	if len(filePaths) > 0 {
		concat.DeleteFiles(filePaths)
	}

	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
