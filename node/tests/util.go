package tests

import (
	"bytes"

	refmtjson "github.com/polydawn/refmt/json"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

// various benchmarks assign their final result here,
// in order to defuse the possibility of their work being elided.
var sink interface{}

func mustNodeFromJsonString(ns ipld.NodeStyle, str string) ipld.Node {
	nb := ns.NewBuilder()
	must.NotError(codec.Unmarshal(nb, refmtjson.NewDecoder(bytes.NewBufferString(str))))
	return nb.Build()
}

func mustSelectorFromJsonString(ns ipld.NodeStyle, str string) selector.Selector {
	// Needing an 'ns' parameter here is sort of off-topic, to be honest.
	//  Someday the selector package will probably contain codegen'd nodes of its own schema, and we'll use those unconditionally.
	//  For now... we'll just use whatever node you're already testing, because it oughta work
	//   (and because it avoids hardcoding any other implementation that might cause import cycle funtimes.).
	seldefn := mustNodeFromJsonString(ns, str)
	sel, err := selector.ParseSelector(seldefn)
	must.NotError(err)
	return sel
}
