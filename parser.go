package main

func parse(tokens []Token) (output string, errors []string) {
	x := 0
	for x < len(tokens) {
		if tokens[x].kind == TokenMultiLineComment {
			output += tokens[x].value
		} else if tokens[x].kind == TokenNewline {
			output += tokens[x].value
		} else if tokens[x].kind == TokenSpace {
			output += tokens[x].value
		}
		x += 1
	}
	return output, errors
}
