module github.com/nanobus/nanobus/example/welcome

go 1.17

require (
	github.com/go-logr/logr v1.2.2
	github.com/go-logr/zapr v1.2.3
	github.com/mattn/go-colorable v0.1.12
	github.com/nanobus/adapter-go v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.21.0
)

require (
	github.com/google/uuid v1.1.2 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jjeffcaii/reactor-go v0.5.2 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/panjf2000/ants/v2 v2.4.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rsocket/rsocket-go v0.8.8 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

replace github.com/nanobus/adapter-go => ../../../adapter-go
