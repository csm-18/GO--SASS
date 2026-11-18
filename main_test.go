package main

// import (
// 	"strings"
// 	"testing"
// )

// func TestCompileNestedSCSS(t *testing.T) {
// 	input := `
// $primary: red;
// nav {
//   a {
//     color: $primary;
//   }
// }
// `

// 	got := compileSCSS(input)
// 	if !strings.Contains(got, "nav a") {
// 		t.Fatalf("expected nested selector in output, got: %s", got)
// 	}
// 	if !strings.Contains(got, "color: red;") {
// 		t.Fatalf("expected variable replacement in output, got: %s", got)
// 	}
// }

// func TestCompileSCSSComments(t *testing.T) {
// 	input := `
// // main comment
// nav {
//   color: red; // inline comment
// }
// `

// 	got := compileSCSS(input)
// 	if !strings.Contains(got, "/* main comment */") {
// 		t.Fatalf("expected SCSS line comment converted to CSS comment, got: %s", got)
// 	}
// 	if !strings.Contains(got, "/* inline comment */") {
// 		t.Fatalf("expected inline SCSS comment converted to CSS comment, got: %s", got)
// 	}
// }
