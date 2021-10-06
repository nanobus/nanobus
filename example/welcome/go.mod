module github.com/nanobus/nanobus/example/welcome

go 1.17

require (
	github.com/nanobus/go-functions v0.0.0-20210930143304-4e5a7c52d459
	github.com/oklog/run v1.1.0
)

require (
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
)

replace github.com/nanobus/go-functions => ../../../go-functions
