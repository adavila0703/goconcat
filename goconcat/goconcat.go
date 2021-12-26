package goconcat

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	ErrReadingDirectories = errors.New("error reading the directories")
)

func Goconcat() error {
	filePaths, err := getFilePaths(".", []string{".git"}, ".go", "mock_")
	if err != nil {
		return err
	}

	test, err := ioutil.ReadFile(filePaths[0])
	if err != nil {
		return err
	}

	fmt.Println("file", test)
	fmt.Println("file", test)

	fmt.Println(filePaths)
	return nil
}

func getFilePaths(path string, ignoredDirectories []string, fileType string, prefix string) ([]string, error) {
	var filePaths []string
	var directories []string
	var currentPath string

	directories = append(directories, path)

	for {
		for _, directory := range directories {
			if checkDirectoryIgnore(directory, ignoredDirectories) {
				continue
			}
			currentPath += directory
			files, err := ioutil.ReadDir(directory)
			if err != nil {
				return nil, ErrReadingDirectories
			}

			directories = removePathFromDirectories(directories, directory)

			for _, file := range files {
				path := currentPath + "/" + file.Name()
				if file.IsDir() {
					directories = append(directories, path)
				} else {
					if strings.HasSuffix(file.Name(), fileType) && strings.HasPrefix(file.Name(), prefix) {
						filePaths = append(filePaths, path)
					}
				}
			}
			currentPath = ""
		}
		if len(directories) == 0 {
			break
		}
	}

	return filePaths, nil
}

func removePathFromDirectories(directories []string, path string) []string {
	var newDirectories []string
	for _, directory := range directories {
		if directory == path {
			continue
		}

		newDirectories = append(newDirectories, directory)
	}
	return newDirectories
}

func checkDirectoryIgnore(directory string, ignoredDirectories []string) bool {
	for _, d := range ignoredDirectories {
		if d == directory {
			return true
		}
	}
	return false
}
