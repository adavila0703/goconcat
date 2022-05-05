package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/adavila0703/goconcat"
	"github.com/adavila0703/goconcat/internal/utils"
	"github.com/adavila0703/goconcat/pkg/concat"
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
)

var (
	app = &cli.App{
		Name: `
		██████╗  ██████╗  ██████╗ ██████╗ ███╗   ██╗ ██████╗ █████╗ ████████╗
		██╔════╝ ██╔═══██╗██╔════╝██╔═══██╗████╗  ██║██╔════╝██╔══██╗╚══██╔══╝
		██║  ███╗██║   ██║██║     ██║   ██║██╔██╗ ██║██║     ███████║   ██║   
		██║   ██║██║   ██║██║     ██║   ██║██║╚██╗██║██║     ██╔══██║   ██║   
		╚██████╔╝╚██████╔╝╚██████╗╚██████╔╝██║ ╚████║╚██████╗██║  ██║   ██║   
		 ╚═════╝  ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝╚═╝  ╚═╝   ╚═╝   
																																					
		`,
		Usage:     "Concatenate your Go files!",
		UsageText: "goconcat can be used to concatenate multiple go files.",
		Flags:     []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name: "simple",
				Usage: `A simple concatenation of multiple Go files.
					args: [root-path][ignored-directories][file-prefix][destination][delete-old-files]
				`,
				UsageText: `
				args details:
				root-path: the directory you would like to start
					example: "./home" - will only walk the home directory.
				ignored-directories: directories you would like to ignore
					example: "home,pkg,internal" - make sure you use a comma as a common delimite.
				file-prefix: if you want to only point to files with a prefix
					example: "mocks_" - will only look for files with a prefix of 'mocks_'.
				destination: directory you would like your files to be moved to
				delete-old-files: if you would like to delete your old files
					exmple: true or false
				`,
				Action:    simpleConcat,
				Aliases:   []string{"s"},
				ArgsUsage: "[root-path][ignored-directories][file-prefix][destination][delete-old-files]",
			},
			{
				Name: "json",
				Usage: `Json allows you to write options in a json file. Use -h for help on the json structure.
					args: [path-to-json]`,
				UsageText: `
				Json seetings example
				{
					"rootPath": ".",
					"ignoredDirectories": ["file1", "file2"],
					"filePrefix": ["mocks_", "mock"],
					"destination": "./destination",
					"deleteOldFiles": false,
					"concatPkg": false
				}
				`,
				Action:    jsonConcat,
				Aliases:   []string{"j"},
				ArgsUsage: "[path-to-json]",
			},
		},
	}
)

func simpleConcat(c *cli.Context) error {
	rootPath := c.Args().Get(0)
	ignoredDirectories := utils.StringToType[utils.Directory](
		strings.Split(c.Args().Get(1), ","),
	)

	filePrefix := utils.StringToType[utils.PrefixType](
		strings.Split(c.Args().Get(2), ","),
	)

	destination := c.Args().Get(3)
	deleteOldFiles, err := strconv.ParseBool(c.Args().Get(4))
	if err != nil {
		log.Fatal(concat.ErrBoolCouldNotBeParsed)
	}

	fileTypes := []utils.FileType{
		utils.FileGo,
	}

	options := goconcat.NewOptions(
		rootPath,
		ignoredDirectories,
		filePrefix,
		destination,
		deleteOldFiles,
		false,
		false,
		fileTypes,
	)

	if err := goconcat.GoConcat(options); err != nil {
		log.Fatal(err)
	}

	return nil
}

func jsonConcat(c *cli.Context) error {
	jsonFilePath := c.Args().Get(0)
	if jsonFilePath == "" {
		log.Fatal(concat.ErrNoFilePathForJson)
	}

	file, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var options goconcat.Options

	err = jsoniter.UnmarshalFromString(string(file), &options)
	if err != nil {
		log.Fatal(err)
	}

	options.FileType = append(options.FileType, utils.FileGo)

	if err := goconcat.GoConcat(&options); err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
