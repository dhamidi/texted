# Technical specifications for `texted`

- use only the Go stdlib for the core implementation
- you must use I/O-less implementations through appropriate interfaces (often io.Reader, io.Writer will suffice)
- the `texted` package is defined in the repository root directory
- the command line interface lives in `./cmd/texted` and is implemented using cobra
- the actual cobra commands live in `commands/`
