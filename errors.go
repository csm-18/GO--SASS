package main

import (
	"fmt"
)

type CompileError struct {
	Line    int
	Message string
}

func (e CompileError) Error() string {
	return fmt.Sprintf(
		"line %d: %s",
		e.Line,
		e.Message,
	)
}

func ErrorUnexpectedClosingBrace(line int) error {
	return CompileError{
		Line:    line,
		Message: "unexpected closing brace '}'",
	}
}

func ErrorUnclosedBlock(line int, selector string) error {
	return CompileError{
		Line:    line,
		Message: fmt.Sprintf("block '%s' was never closed", selector),
	}
}
