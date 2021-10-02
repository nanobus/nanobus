.PHONY: all deps codegen watch build run clean doc test

all: deps codegen build

deps:
	npm install

codegen:
	nanogen generate

watch:
	nanogen watch

build:
	npm run build

run:
	npm run ts-node-server

clean:
	rm -Rf build

doc:

test: build
	npm run test
