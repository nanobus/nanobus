APP_ID                 ?= customers
HTTP_LISTEN_ADDRESS    ?= :18081
REST_LISTEN_ADDRESS    ?= :8091
BUS_LISTEN_ADDRESS     ?= 127.0.0.1:32321
METRICS_PORT           ?= 19091

PLACEMENT_HOST_ADDRESS ?= localhost:50005

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
	@nanobus --components-path ../components --app-id $(APP_ID) --metrics-port $(METRICS_PORT) --placement-host-address $(PLACEMENT_HOST_ADDRESS) --http-listen-addr $(HTTP_LISTEN_ADDRESS) --rest-listen-addr $(REST_LISTEN_ADDRESS) --bus-listen-addr $(BUS_LISTEN_ADDRESS)

clean:
	rm -Rf build

doc:

test: build
	npm run test
