set dotenv-load

# Disable CGo globally
export CGO_ENABLED := "0"

@default:
	just --list --unsorted

@deps:
	go mod tidy
	go mod download

@build: deps
	go build ./...

@lint: deps
	go vet ./...

@test: deps
	go test ./...

@run: deps
	go run ./...

alias fmt := format
@format:
	# try to use golines, fall back to go fmt
	(which golines >> /dev/null && golines -w -m 80 --ignore-generated *.go) || go fmt

alias dbg := debug
@debug:
	mkdir -p build/
	dlv test --output build/ ./...
