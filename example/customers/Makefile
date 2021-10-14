APP_ID                 ?= customers
APP_PORT               ?= 9000
HTTP_LISTEN_ADDRESS    ?= :18081
REST_LISTEN_ADDRESS    ?= :8091
BUS_LISTEN_ADDRESS     ?= 127.0.0.1:32321
METRICS_PORT           ?= 19091

PLACEMENT_HOST_ADDRESS ?= localhost:50005

.PHONY: all deps codegen watch build run run-nanobus run-app clean doc test

all: deps codegen tidy build

deps:

codegen:
	nanogen generate

watch:
	nanogen watch

tidy:
	go mod tidy

build:
	mkdir -p build
	go build -o build/$(APP_ID) cmd/main.go

run:
	@nanobus --components-path ../components --app-id $(APP_ID) --metrics-port $(METRICS_PORT) --placement-host-address $(PLACEMENT_HOST_ADDRESS) --http-listen-addr $(HTTP_LISTEN_ADDRESS) --rest-listen-addr $(REST_LISTEN_ADDRESS) --bus-listen-addr $(BUS_LISTEN_ADDRESS) --app-port $(APP_PORT) -- go run cmd/main.go

run-nanobus:
	@nanobus --components-path ../components --app-id $(APP_ID) --metrics-port $(METRICS_PORT) --placement-host-address $(PLACEMENT_HOST_ADDRESS) --http-listen-addr $(HTTP_LISTEN_ADDRESS) --rest-listen-addr $(REST_LISTEN_ADDRESS) --bus-listen-addr $(BUS_LISTEN_ADDRESS) --app-port $(APP_PORT)

run-app:
	@BUS_URI=http://$(BUS_LISTEN_ADDRESS) PORT=$(APP_PORT) go run cmd/main.go

clean:
	rm -Rf build

doc:

test:
	go test --count=1 ./pkg/...
