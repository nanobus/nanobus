.PHONY: all deps codegen build clean doc test

all: deps codegen tidy build

deps:

codegen:
	wapc generate codegen.yaml

tidy:
	go mod tidy

build:
	mkdir -p build
	go build -o build/customers cmd/main.go

run:
	@go run cmd/main.go

clean:
	rm -Rf build

doc:

test:
	go test --count=1 ./pkg/...
