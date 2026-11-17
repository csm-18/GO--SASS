# Contributing to GO--SASS

Thanks for helping improve GO--SASS.

## How to contribute

1. Fork the repository.
2. Create a feature branch for your work.
3. Make your changes and keep them focused.
4. Run the test suite before submitting.
5. Open a pull request with a clear summary.

## Development setup

Clone the project:

```bash
git clone https://github.com/your-username/GO--SASS.git
cd GO--SASS
```

Install dependencies (if any are added later):

```bash
go mod tidy
```

Run tests:

```bash
go test ./...
```

Build the compiler:

```bash
go build -o GO-SASS.exe
```

## Project conventions

- Keep code simple and easy to read.
- Follow standard Go formatting and naming conventions.
- Prefer small, focused changes over broad rewrites.
- Add tests for bug fixes and new compiler features when possible.
- Keep behavior consistent with the current SCSS-like syntax supported by the project.

## Reporting issues

When reporting a bug, include:

- the input SCSS file
- the expected CSS output
- the actual output
- the command used to run the compiler
- any relevant error messages

## Pull request guidelines

Please ensure your PR:

- explains the purpose of the change
- mentions any affected syntax/features
- includes test results when relevant
- remains focused on one topic or fix

Thank you for contributing to GO--SASS.
