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

	tokens, errors := lexer(srcPath, src)
	if len(errors) != 0 {
		for _, error := range errors {
			fmt.Println(error)
		}
	}
	for _, token := range tokens {
		fmt.Printf("name: %v,value: %v, index: %v\n", token.name, token.value, token.index)
	}

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
