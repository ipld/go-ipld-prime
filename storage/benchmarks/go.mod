module github.com/ipld/go-ipld-prime/storage/benchmarks

go 1.16

replace github.com/ipld/go-ipld-prime => ../..

replace github.com/ipld/go-ipld-prime/storage/dsadapter => ../dsadapter

require (
	github.com/ipfs/go-ds-flatfs v0.5.0
	github.com/ipld/go-ipld-prime v0.14.1
	github.com/ipld/go-ipld-prime/storage/dsadapter v0.0.0-20211130004103-85b37597b213
)
