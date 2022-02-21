FROM golang:1.17.5 as build-env

WORKDIR /go/src/app

#RUN apt-get update && apt-get install -y gcc

#COPY . .

#RUN go mod init
#RUN go get -d -v ./...
#RUN go vet -v
#RUN go test -v

#RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/bin/customers cmd/main.go
COPY build/customers-linux-amd64 /go/bin/customers

FROM nanobus/base

COPY --from=build-env /go/bin/customers /app/customers
COPY components /components
COPY common.yaml /app/common.yaml
COPY bus.yaml /app/bus.yaml
COPY spec.apex /app/spec.apex
CMD ["--components-path", "/components", "--app-id", "customers", "--metrics-port", "19091", "--placement-host-address", "host.docker.internal:50005", "--http-listen-addr", ":18081", "--rest-listen-addr", ":8091", "--bus-listen-addr", "127.0.0.1:32321", "--app-port", "9000", "--", "/app/customers"]
