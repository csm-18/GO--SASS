package main

import "fmt"

func compile(srcPath, outPath string) (errors []string) {
	fmt.Println("Compiling!")
	fmt.Printf("srcPath: %s, outPath: %s\n", srcPath, outPath)

	return errors
}
