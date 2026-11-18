package main

import (
	"fmt"
	"os"
)

func compile(srcPath, outPath string) (errors []string) {
	src, readError := readFile(srcPath)
	if readError != "" {
		errors = append(errors, readError)
		return errors
	}

	fmt.Println(src)

	return errors
}

func readFile(path string) (content string, readError string) {
	data, err := os.ReadFile(path)
	if err != nil {
		readError = fmt.Sprintf("error: Cannot read file: %s", path)
	}
	content = string(data)
	return content, readError
}
