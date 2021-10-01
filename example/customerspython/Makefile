.PHONY: all deps codegen watch run clean doc test

all: deps codegen

deps:
	pip3 install -r requirements.txt

codegen:
	nanogen generate

watch:
	nanogen watch

run:
	python3 service.py

doc:

test:
