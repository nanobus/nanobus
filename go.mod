module github.com/nanobus/nanobus

go 1.16

require (
	github.com/actgardner/gogen-avro/v9 v9.0.0
	github.com/agrea/ptr v0.0.0-20180711073057-77a518d99b7b
	github.com/antonmedv/expr v1.9.0
	github.com/cenkalti/backoff/v4 v4.1.1
	github.com/dapr/components-contrib v1.4.0-rc1
	github.com/dapr/dapr v1.3.0
	github.com/dapr/kit v0.0.2-0.20210614175626-b9074b64d233
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da
	github.com/google/cel-go v0.7.3
	github.com/gorilla/mux v1.8.0
	github.com/hamba/avro v1.5.6
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/anko v0.1.8
	github.com/mitchellh/mapstructure v1.4.1
	github.com/nanobus/go-functions v0.0.0-20210811180606-4acc97e47cad
	github.com/oklog/run v1.1.0
	github.com/sony/gobreaker v0.4.1
	github.com/spf13/cast v1.4.0
	github.com/stretchr/testify v1.7.0
	github.com/valyala/fasthttp v1.28.0
	github.com/vmihailenco/msgpack/v5 v5.3.4
	github.com/wapc/wapc-go v0.3.0
	github.com/wapc/widl-go v0.0.0-20210922050642-a089e96973c3
	google.golang.org/genproto v0.0.0-20210813162853-db860fec028c
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/nanobus/go-functions => ../go-functions

replace github.com/dapr/dapr => ../../dapr/dapr
