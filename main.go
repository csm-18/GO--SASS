package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type cssBlock struct {
	name       string
	selector   string
	properties []string
}

func main() {
	inputPath := flag.String("i", "", "Input SCSS file")
	flag.Parse()

	args := flag.Args()
	if len(args) >= 2 && args[0] == "i" {
		*inputPath = args[1]
	} else if *inputPath == "" && len(args) > 0 {
		*inputPath = args[0]
	}
	if *inputPath == "" {
		*inputPath = "test.scss"
	}

	fmt.Println("Reading SCSS file:", *inputPath)

	data, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		panic(err)
	}

	fmt.Println("Compiling CSS...")
	output := compileSCSS(string(data))

	fmt.Println("Writing output to output.css")
	if err := os.WriteFile("output.css", []byte(output), 0644); err != nil {
		fmt.Println("Error writing output.css:", err)
		panic(err)
	}

	fmt.Println("Done.")
}

func compileSCSS(source string) string {
	source = convertSCSSCommentsToCSS(source)
	variables := parseVariables(source)
	lines := strings.Split(source, "\n")
	stack := []cssBlock{}
	output := []string{}

	fmt.Println("Searching for SCSS variables...")
	for name, value := range variables {
		fmt.Println("Found variable:", name, "=", value)
	}
	if len(variables) == 0 {
		fmt.Println("No variables found.")
	}

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "/*") && strings.HasSuffix(line, "*/") {
			output = append(output, line)
			continue
		}
		if strings.HasPrefix(line, "$") {
			continue
		}

		if openIndex := strings.Index(line, "{"); openIndex >= 0 {
			selector := strings.TrimSpace(line[:openIndex])
			if selector == "" {
				continue
			}
			fullSelector := selector
			if len(stack) > 0 {
				fullSelector = strings.TrimSpace(stack[len(stack)-1].selector + " " + selector)
			}

			stackNames := []string{}
			for _, block := range stack {
				stackNames = append(stackNames, block.name)
			}
			stackNames = append(stackNames, selector)
			fmt.Printf("PUSH %s\nSTACK = [%s]\n", selector, strings.Join(stackNames, " "))
			stack = append(stack, cssBlock{name: selector, selector: fullSelector})
			continue
		}

		if strings.Contains(line, "}") {
			if len(stack) == 0 {
				continue
			}
			block := stack[len(stack)-1]
			if len(block.properties) > 0 {
				output = append(output, formatBlock(block.selector, block.properties))
			}
			stackNames := []string{}
			for _, parent := range stack[:len(stack)-1] {
				stackNames = append(stackNames, parent.name)
			}
			fmt.Printf("POP\nSTACK = [%s]\n", strings.Join(stackNames, " "))
			stack = stack[:len(stack)-1]
			continue
		}

		if strings.Contains(line, ":") {
			propertyText, commentText := splitCodeAndTrailingComment(line)
			property := replaceVariables(strings.TrimSuffix(strings.TrimSpace(propertyText), ";"), variables)
			if property == "" || len(stack) == 0 {
				continue
			}
			fullProperty := property + ";"
			if commentText != "" {
				fullProperty += " " + commentText
			}
			stack[len(stack)-1].properties = append(stack[len(stack)-1].properties, fullProperty)
			fmt.Println("PROPERTY", property)
		}
	}

	for len(stack) > 0 {
		block := stack[len(stack)-1]
		if len(block.properties) > 0 {
			output = append(output, formatBlock(block.selector, block.properties))
		}
		stack = stack[:len(stack)-1]
	}

	return strings.Join(output, "\n") + "\n"
}

func convertSCSSCommentsToCSS(source string) string {
	lines := strings.Split(source, "\n")
	converted := make([]string, 0, len(lines))

	for _, rawLine := range lines {
		line := rawLine
		for {
			commentIndex := strings.Index(line, "//")
			if commentIndex == -1 {
				break
			}

			before := strings.TrimRight(line[:commentIndex], " \t")
			after := strings.TrimSpace(line[commentIndex+2:])

			if strings.TrimSpace(before) == "" {
				if after == "" {
					line = "/* */"
				} else {
					line = "/* " + after + " */"
				}
				break
			}

			if after == "" {
				line = before
			} else {
				line = before + " /* " + after + " */"
			}
			break
		}
		converted = append(converted, line)
	}

	return strings.Join(converted, "\n")
}

func splitCodeAndTrailingComment(line string) (string, string) {
	trimmed := strings.TrimSpace(line)
	idx := strings.Index(trimmed, "/*")
	if idx < 0 {
		return trimmed, ""
	}
	return strings.TrimSpace(trimmed[:idx]), strings.TrimSpace(trimmed[idx:])
}

func parseVariables(source string) map[string]string {
	variables := make(map[string]string)
	for _, rawLine := range strings.Split(source, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "/*") || !strings.HasPrefix(line, "$") {
			continue
		}

		line, _ = splitCodeAndTrailingComment(line)
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(strings.TrimPrefix(parts[0], "$"))
		value := strings.TrimSpace(parts[1])
		value = strings.TrimSuffix(value, ";")
		variables[name] = value
	}
	return variables
}

func replaceVariables(line string, variables map[string]string) string {
	for name, value := range variables {
		line = strings.ReplaceAll(line, "$"+name, value)
	}
	return strings.TrimSpace(line)
}

func formatBlock(selector string, properties []string) string {
	if len(properties) == 0 {
		return ""
	}

	result := []string{selector + " {"}
	for _, property := range properties {
		result = append(result, "  "+property)
	}
	result = append(result, "}")
	return strings.Join(result, "\n")
}
