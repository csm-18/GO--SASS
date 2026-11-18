package main

import (
	"fmt"
	"os"
	"strings"
)

type cssBlock struct {
	name       string
	selector   string
	properties []string
}

type Options struct {
	version  bool
	help     bool
	watchDir string
	outDir   string
	srcFiles []string
	errors   []string
}

const VERSION = "minisass 0.2"
const HELP = `A small SCSS-to-CSS compiler written in Go.
Usage: minisass <input.scss>

  Options:
    -v              -> Print version.
    -h              -> Print help.
    -w folder-name  -> Watch folder for changes and re-generate .css files.
    -o folder-name  -> Specify where the generated css files should go.`

func main() {
	args := os.Args[1:]
	options := parseCommandLineArgs(args)

	if len(options.errors) != 0 {
		for _, error := range options.errors {
			fmt.Println(error)
		}
		os.Exit(1)
	}

	if options.watchDir != "" && (options.version || options.help || len(options.srcFiles) != 0) {
		fmt.Println("error: -w option cannot be used with other options except -o option.")
		os.Exit(1)
	}

	if options.version && options.help {
		fmt.Println(VERSION)
		fmt.Println(HELP)
		os.Exit(0)
	} else if options.version {
		fmt.Println(VERSION)
		os.Exit(0)
	} else if options.help {
		fmt.Println(VERSION)
		fmt.Println(HELP)
		os.Exit(0)
	} else if options.watchDir != "" {
		// Watch option
		fmt.Println("This option is under construction 😅")
		os.Exit(0)
	} else if len(options.srcFiles) != 0 {
		// Compile SCSS files to CSS files
		fmt.Println("Compiling!")
	} else {
		fmt.Println("error: No .scss source files provided!")
		os.Exit(1)
	}

}

func parseCommandLineArgs(args []string) (options Options) {
	if len(args) == 0 {
		options.version = true
		options.help = true
		return options
	}

	x := 0
	for x < len(args) {
		if args[x] == "-v" {
			options.version = true
		} else if args[x] == "-h" {
			options.help = true
		} else if strings.HasSuffix(args[x], ".scss") {
			options.srcFiles = append(options.srcFiles, args[x])
		} else if args[x] == "-o" {
			x += 1
			if x < len(args) {
				if isDirectory(args[x]) {
					options.outDir = args[x]
				} else {
					options.errors = append(options.errors, "error: Invalid path provided to -o option: "+args[x])
				}
			} else {
				error := "error: No path provided to -o option."
				options.errors = append(options.errors, error)
			}
		} else if args[x] == "-w" {
			x += 1
			if x < len(args) {
				if isDirectory(args[x]) {
					options.watchDir = args[x]
				} else {
					options.errors = append(options.errors, "error: Invalid path provided to -w option: "+args[x])
				}
			} else {
				error := "error: No path provided to -w option."
				options.errors = append(options.errors, error)
			}

		} else {
			error := "error: Unknown option: " + args[x]
			options.errors = append(options.errors, error)
		}
		x += 1
	}

	return options
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	return true
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
