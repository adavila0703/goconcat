package goconcat

import (
	"errors"
	"fmt"
	"io/ioutil"
)

var (
	ErrReadingDirectories = errors.New("error reading the directories")
)

func Goconcat() error {
	fileNames, err := getFileNames(".")
	if err != nil {
		return err
	}

	fmt.Println(fileNames)
	return nil
}

func getFileNames(path string) ([]string, error) {
	var fileNames []string
	var directories []string

	directPath := path

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, ErrReadingDirectories
	}

	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		} else {
			directories = append(directories, directPath+"/"+file.Name())
		}
	}

	for {

		for index, directory := range directories {
			files, err := ioutil.ReadDir(directory)
			if err != nil {
				return nil, ErrReadingDirectories
			}

			fmt.Println(directories, directory)
			directories = append(directories[:index], directories[index+1:]...)
			fmt.Println(directories, directory)

			for _, file := range files {
				if !file.IsDir() {
					fileNames = append(fileNames, file.Name())
				} else {
					directories = append(directories, directPath+"/"+file.Name())
				}
			}
		}

		if len(directories) == 0 {
			break
		}
	}

	return fileNames, nil
}
