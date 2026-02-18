module github.com/ipld/go-ipld-prime/storage/benchmarks

go 1.24.0

replace github.com/ipld/go-ipld-prime => ../..

replace github.com/ipld/go-ipld-prime/storage/dsadapter => ../dsadapter

require (
	github.com/ipfs/go-ds-flatfs v0.6.0
	github.com/ipld/go-ipld-prime v0.20.0
	github.com/ipld/go-ipld-prime/storage/dsadapter v0.20.0
)

require (
	github.com/alexbrainman/goissue34681 v0.0.0-20191006012335-3fc7a47baff5 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/ipfs/go-datastore v0.9.1 // indirect
	github.com/ipfs/go-log/v2 v2.9.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
)
