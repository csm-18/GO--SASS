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

	tokens, lexErrors := lexer(srcPath, src)
	if len(lexErrors) != 0 {
		errors = append(errors, lexErrors...)
		return errors
	}

	// for _, token := range tokens {
	// 	fmt.Printf("name: %v,value: %v, index: %v\n", token.name, token.value, token.index)
	// }

	output, parseErrors := parse(tokens)
	if len(parseErrors) != 0 {
		errors = append(errors, parseErrors...)
		return errors
	}

	writeError := createFile(outPath, output)
	if writeError != "" {
		errors = append(errors, writeError)
		return errors
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

func createFile(path string, content string) (error string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		error = "error: Unable to generate file: " + path
	}
	return error
}
