module github.com/ipld/go-ipld-prime/storage/benchmarks

go 1.23

replace github.com/ipld/go-ipld-prime => ../..

replace github.com/ipld/go-ipld-prime/storage/dsadapter => ../dsadapter

require (
	github.com/ipfs/go-ds-flatfs v0.5.5
	github.com/ipld/go-ipld-prime v0.20.0
	github.com/ipld/go-ipld-prime/storage/dsadapter v0.20.0
)

require (
	github.com/alexbrainman/goissue34681 v0.0.0-20191006012335-3fc7a47baff5 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/ipfs/go-datastore v0.8.2 // indirect
	github.com/ipfs/go-log/v2 v2.5.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
	golang.org/x/sys v0.30.0 // indirect
)
