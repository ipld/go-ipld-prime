package tests

import (
	"bytes"

	refmtjson "github.com/polydawn/refmt/json"

	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
	"github.com/ipld/go-ipld-prime/_rsrch/nodesolution/codec"
)

// various benchmarks assign their final result here,
// in order to defuse the possibility of their work being elided.
var sink interface{}

func mustNodeFromJsonString(nb ipld.NodeBuilder, str string) ipld.Node {
	err := codec.Unmarshal(nb, refmtjson.NewDecoder(bytes.NewBufferString(str)))
	if err != nil {
		panic(err)
	}
	return nb.Build()
}
