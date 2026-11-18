package main

import "fmt"

type TokenKind int

const (
	TokenMultiLineComment = iota // 0
	TokenSingleQuoteString
	TokenDoubleQuoteString
	TokenNewline
	TokenSpace
)

type Token struct {
	name  string
	kind  TokenKind
	value string
	index int
}

func lexer(srcPath, src string) (tokens []Token, errors []string) {
	x := 0
	for x < len(src) {
		if src[x] == '\n' {
			token := Token{"Newline", TokenNewline, "\n", x}
			tokens = append(tokens, token)
		} else if src[x] == ' ' {
			token := Token{"Space", TokenSpace, " ", x}
			tokens = append(tokens, token)
		} else if src[x] <= ' ' {
			// Skip ASCII control characters
		} else if src[x] == '/' && x+1 < len(src) {
			if src[x+1] == '/' {
				// Skip single-line comments
				foundNewline := false
				y := x + 2
				for y < len(src) {
					if src[y] == '\n' {
						foundNewline = true
						break
					}
					y += 1
				}
				if foundNewline {
					x = y
					continue
				} else {
					// EOF
					break
				}
			} else if src[x+1] == '*' {
				// Muti-line comments parsing
				foundCommentClosing := false
				y := x + 2
				for y < len(src) {
					if src[y] == '*' && y+1 < len(src) && src[y+1] == '/' {
						y += 1
						foundCommentClosing = true
						break
					}
					y += 1
				}

				if foundCommentClosing {
					token := Token{"MultiLineComment", TokenMultiLineComment, src[x : y+1], x}
					tokens = append(tokens, token)
					x = y + 1
					continue
				} else {
					line, col := getLineCol(src, x)
					error := fmt.Sprintf("error: Unclosed Multi-line comment in file '%v' on line %v at %v", srcPath, line, col)
					errors = append(errors, error)
					return tokens, errors
				}
			}
		} else if src[x] == '"' {
			// Double-quote string parsing

			endQuote := false
			y := x + 1
			for y < len(src) {
				if src[y] == '\\' {
					y += 2
					continue
				} else if src[y] == '"' {
					endQuote = true
					break
				}
				y += 1
			}

			if endQuote {
				token := Token{"DoubleQuoteString", TokenDoubleQuoteString, src[x : y+1], x}
				tokens = append(tokens, token)
				x = y + 1
				continue
			} else {
				line, col := getLineCol(src, x)
				error := fmt.Sprintf("error: Unclosed string-literal in file '%v' on line %v at %v", srcPath, line, col)
				errors = append(errors, error)
				return tokens, errors
			}
		} else if src[x] == '\'' {
			// Single-quote string parsing

			endQuote := false
			y := x + 1
			for y < len(src) {
				if src[y] == '\\' {
					y += 2
					continue
				} else if src[y] == '\'' {
					endQuote = true
					break
				}
				y += 1
			}

			if endQuote {
				token := Token{"DoubleQuoteString", TokenDoubleQuoteString, src[x : y+1], x}
				tokens = append(tokens, token)
				x = y + 1
				continue
			} else {
				line, col := getLineCol(src, x)
				error := fmt.Sprintf("error: Unclosed string-literal in file '%v' on line %v at %v", srcPath, line, col)
				errors = append(errors, error)
				return tokens, errors
			}
		}
		x += 1
	}

	return tokens, errors
}

func getLineCol(src string, index int) (line, col int) {
	line, col = 1, 1
	for i := 0; i < index; i++ {
		if src[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return line, col
}
