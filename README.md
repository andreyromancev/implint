# Implint
A simple tool for linting package imports.

## Usage
Execute `go run cmd/implint/implint.go -dir /path/to/lint/`

Provided path should contain file named `.implint.yml` with defined hierarchy.
See `example` dir for more details (`go run cmd/implint/implint.go -dir example`).

**Warning:** only supports `GOPATH` imports.
