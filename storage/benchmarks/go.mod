module github.com/ipld/go-ipld-prime/storage/benchmarks

go 1.16

replace github.com/ipld/go-ipld-prime => ../..

replace github.com/ipld/go-ipld-prime/storage/dsadapter => ../dsadapter

require (
	github.com/ipfs/go-ds-flatfs v0.4.5
	github.com/ipld/go-ipld-prime v0.12.3
	github.com/ipld/go-ipld-prime/storage/dsadapter v0.0.0-20211022093231-ebf675a9bd6d
)
