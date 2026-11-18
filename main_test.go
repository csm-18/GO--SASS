package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompileNestedSCSS(t *testing.T) {
	input := `
$primary: red;
nav {
  a {
    color: $primary;
  }
}
`

	got := compileSCSS(input, false)
	if !strings.Contains(got, "nav a") {
		t.Fatalf("expected nested selector in output, got: %s", got)
	}
	if !strings.Contains(got, "color: red;") {
		t.Fatalf("expected variable replacement in output, got: %s", got)
	}
}

func TestCompileSCSSComments(t *testing.T) {
	input := `
// main comment
nav {
  color: red; // inline comment
}
`

	got := compileSCSS(input, false)
	if !strings.Contains(got, "/* main comment */") {
		t.Fatalf("expected SCSS line comment converted to CSS comment, got: %s", got)
	}
	if !strings.Contains(got, "/* inline comment */") {
		t.Fatalf("expected inline SCSS comment converted to CSS comment, got: %s", got)
	}
}

func TestCompileMinify(t *testing.T) {
	input := `
$color: blue;
button {
  color: $color;
}
`

	got := compileSCSS(input, true)
	if !strings.Contains(got, "button{color: blue;}") {
		t.Fatalf("expected minified CSS output, got: %q", got)
	}
}

func TestPrintUsage(t *testing.T) {
	var buf bytes.Buffer
	printUsage(&buf)
	output := buf.String()

	expectedSubstrings := []string{
		"Usage:",
		"Flags:",
		"-i string",
		"-o string",
		"-minify",
		"-h, --help",
		"Arguments:",
		"Examples:",
	}

	for _, expected := range expectedSubstrings {
		if !strings.Contains(output, expected) {
			t.Errorf("expected usage output to contain %q, but got:\n%s", expected, output)
		}
	}
}

func TestErrorTypes(t *testing.T) {
	unclosedErr := ErrorUnclosedBlock(4, "nav")
	expectedUnclosed := "line 4: block 'nav' was never closed"
	if unclosedErr.Error() != expectedUnclosed {
		t.Errorf("expected %q, got %q", expectedUnclosed, unclosedErr.Error())
	}

	closingBraceErr := ErrorUnexpectedClosingBrace(7)
	expectedBrace := "line 7: unexpected closing brace '}'"
	if closingBraceErr.Error() != expectedBrace {
		t.Errorf("expected %q, got %q", expectedBrace, closingBraceErr.Error())
	}
}

