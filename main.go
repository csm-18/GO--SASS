package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type cssBlock struct {
	name       string
	selector   string
	properties []string
	Line       int
}

func printUsage(w io.Writer) {
	prog := "GO-SASS.exe"
	if len(os.Args) > 0 {
		base := filepath.Base(os.Args[0])
		if !strings.Contains(base, "go-build") && base != "main" && base != "main.exe" && !strings.HasPrefix(base, "minisass") {
			prog = base
		}
	}

	fmt.Fprintf(w, "GO-SASS - A small SCSS-to-CSS compiler written in Go\n\n")
	fmt.Fprintf(w, "Usage:\n")
	fmt.Fprintf(w, "  %s [flags] [input file]\n", prog)
	fmt.Fprintf(w, "  %s -i <input.scss> [-o <output.css>] [-minify]\n\n", prog)
	fmt.Fprintf(w, "Flags:\n")
	fmt.Fprintf(w, "  -i string\n")
	fmt.Fprintf(w, "    \tInput SCSS file path (default: \"test.scss\")\n")
	fmt.Fprintf(w, "  -o string\n")
	fmt.Fprintf(w, "    \tOutput CSS file path (default: \"output.css\")\n")
	fmt.Fprintf(w, "  -minify\n")
	fmt.Fprintf(w, "    \tMinify CSS output\n")
	fmt.Fprintf(w, "  -h, --help\n")
	fmt.Fprintf(w, "    \tDisplay this help and usage information\n\n")
	fmt.Fprintf(w, "Arguments:\n")
	fmt.Fprintf(w, "  [input file]\n")
	fmt.Fprintf(w, "    \tOptional positional argument specifying the input SCSS file\n\n")
	fmt.Fprintf(w, "Examples:\n")
	fmt.Fprintf(w, "  %s -i input.scss -o output.css\n", prog)
	fmt.Fprintf(w, "  %s -i input.scss -minify\n", prog)
	fmt.Fprintf(w, "  %s input.scss\n", prog)
	fmt.Fprintf(w, "  %s\n", prog)
}

func main() {
	flag.Usage = func() {
		printUsage(flag.CommandLine.Output())
	}

	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" || arg == "-help" || arg == "help" {
			flag.CommandLine.SetOutput(os.Stdout)
			if arg == "help" {
				printUsage(os.Stdout)
				return
			}
			break
		}
	}

	inputPath := flag.String("i", "", "Input SCSS file")
	outputPath := flag.String("o", "", "Output CSS file Path")
	minify := flag.Bool("minify", false, "Minify CSS output")
	flag.Parse()

	args := flag.Args()
	if len(args) >= 4 && args[0] == "i" && args[2] == "o" {
		*inputPath = args[1]
		*outputPath = args[3]
	} else if *inputPath == "" && len(args) > 0 {
		*inputPath = args[0]
	}

	if *inputPath == "" {
		*inputPath = "test.scss"
	}

	if *outputPath == "" {
		*outputPath = "output.css"
	}

	fmt.Println("Reading SCSS file:", *inputPath)

	data, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	fmt.Println("Compiling CSS...")
	output := compileSCSS(string(data), *minify)

	fmt.Println("Writing output to :", *outputPath)
	if err := os.WriteFile(*outputPath, []byte(output), 0644); err != nil {
		fmt.Println("Error writing output file:", err)
		os.Exit(1)
	}

	fmt.Println("Done.")
}

func compileSCSS(source string, minify bool) string {
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

	for lineNum, rawLine := range lines {
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
				parent := stack[len(stack)-1].selector

				if strings.Contains(selector, "&") {
					fullSelector = strings.ReplaceAll(
						selector,
						"&",
						parent,
					)
				} else {
					fullSelector = strings.TrimSpace(
						parent + " " + selector,
					)
				}
			}

			stackNames := []string{}
			for _, block := range stack {
				stackNames = append(stackNames, block.name)
			}
			stackNames = append(stackNames, selector)
			fmt.Printf("PUSH %s\nSTACK = [%s]\n", selector, strings.Join(stackNames, " "))
			stack = append(stack, cssBlock{name: selector, selector: fullSelector, Line: lineNum + 1})
			continue
		}

		if strings.Contains(line, "}") {
			if len(stack) == 0 {
				err := ErrorUnexpectedClosingBrace(lineNum + 1)
				fmt.Println(err)
				os.Exit(1)
			}
			block := stack[len(stack)-1]
			if len(block.properties) > 0 {
				output = append(output, formatBlock(block.selector, block.properties, minify))
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

	if len(stack) > 0 {
		block := stack[len(stack)-1]

		err := ErrorUnclosedBlock(
			block.Line,
			block.selector,
		)

		fmt.Println(err)
		os.Exit(1)
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

func formatBlock(selector string, properties []string, minify bool) string {
	fmt.Println("Minify:", minify)

	if len(properties) == 0 {
		return ""
	}

	if minify {
		result := selector + "{"
		for _, property := range properties {
			property = strings.ReplaceAll(property, ":", ":")
			result += property
		}
		result += "}"
		return result
	}

	result := []string{selector + " {"}
	for _, property := range properties {
		result = append(result, "  "+property)
	}
	result = append(result, "}")
	return strings.Join(result, "\n")
}
