# GO--SASS
A small SCSS-to-CSS compiler written in Go.

## Usage
Build a Windows executable:
```bash
go build -o GO-SASS.exe
```

Run the compiler with an input file:
```powershell
.\GO-SASS.exe -i input.scss
```

You can also set the output file explicitly:
```powershell
.\GO-SASS.exe -i input.scss -o output.css
```

A positional argument also works:
```powershell
.\GO-SASS.exe input.scss
```

If no input file is provided, the program defaults to `test.scss`. If no output file is provided, it writes to `output.css`.

### Example
```scss
$primary-color: red;
$secondary-color: green;
$font: monospace;

// Nested example
nav {
    color: $primary-color;
    font: $font;

    a {
        color: $secondary-color;
    }
}
```

Output:
```css
nav {
  color: red;
  font: monospace;
}

nav a {
  color: green;
}
```

## Supported syntax features

| Feature | Supported |
| --- | --- |
| Nested selectors | Yes |
| Parent selector (`&`) | Yes |
| Variable replacement | Yes |
| SCSS `//` comments converted to CSS `/* ... */` comments | Yes |

## Testing
Run the project tests with:
```bash
go test ./...
```

Compile a sample file:
```bash
go run . -i tests/test.scss -o test.css
```

## Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines and contribution workflow.