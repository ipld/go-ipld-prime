module github.com/ipld/go-ipld-prime/storage/benchmarks

go 1.19

replace github.com/ipld/go-ipld-prime => ../..

replace github.com/ipld/go-ipld-prime/storage/dsadapter => ../dsadapter

require (
	github.com/ipfs/go-ds-flatfs v0.5.1
	github.com/ipld/go-ipld-prime v0.20.0
	github.com/ipld/go-ipld-prime/storage/dsadapter v0.20.0
)

require (
	github.com/alexbrainman/goissue34681 v0.0.0-20191006012335-3fc7a47baff5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/ipfs/go-datastore v0.6.0 // indirect
	github.com/ipfs/go-log v1.0.3 // indirect
	github.com/ipfs/go-log/v2 v2.0.3 // indirect
	github.com/jbenet/goprocess v0.1.4 // indirect
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
)
