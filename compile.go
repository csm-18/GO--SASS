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
		for _, error := range lexErrors {
			fmt.Println(error)
		}
		os.Exit(1)
	}

	// for _, token := range tokens {
	// 	fmt.Printf("name: %v,value: %v, index: %v\n", token.name, token.value, token.index)
	// }

	output, parseErrors := parse(tokens)
	if len(parseErrors) != 0 {
		for _, error := range parseErrors {
			fmt.Println(error)
		}
		os.Exit(1)
	}

	fmt.Println(output)

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
