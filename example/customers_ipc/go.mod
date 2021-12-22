module github.com/nanobus/nanobus/example/customers

go 1.17

require (
	github.com/nanobus/go-functions v0.0.0-20210930143304-4e5a7c52d459
	github.com/oklog/run v1.1.0
	go.nanomsg.org/mangos/v3 v3.3.0
)

require (
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
)

replace github.com/nanobus/go-functions => ../../../go-functions
