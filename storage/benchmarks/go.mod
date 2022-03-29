module github.com/ipld/go-ipld-prime/storage/benchmarks

go 1.17

replace github.com/ipld/go-ipld-prime => ../..

replace github.com/ipld/go-ipld-prime/storage/dsadapter => ../dsadapter

require (
	github.com/ipfs/go-ds-flatfs v0.4.5
	github.com/ipld/go-ipld-prime v0.12.3
	github.com/ipld/go-ipld-prime/storage/dsadapter v0.0.0-20211022093231-ebf675a9bd6d
)

require (
	github.com/alexbrainman/goissue34681 v0.0.0-20191006012335-3fc7a47baff5 // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/ipfs/go-datastore v0.4.6 // indirect
	github.com/ipfs/go-log v1.0.3 // indirect
	github.com/ipfs/go-log/v2 v2.0.3 // indirect
	github.com/jbenet/goprocess v0.1.4 // indirect
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
)
