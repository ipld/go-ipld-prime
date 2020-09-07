package tests

import (
	"strings"

	refmtjson "github.com/polydawn/refmt/json"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

// various benchmarks assign their final result here,
// in order to defuse the possibility of their work being elided.
var sink interface{}

func mustNodeFromJsonString(np ipld.NodePrototype, str string) ipld.Node {
	nb := np.NewBuilder()
	must.NotError(codec.Unmarshal(nb, refmtjson.NewDecoder(strings.NewReader(str))))
	return nb.Build()
}

func mustSelectorFromJsonString(np ipld.NodePrototype, str string) selector.Selector {
	// Needing an 'ns' parameter here is sort of off-topic, to be honest.
	//  Someday the selector package will probably contain codegen'd nodes of its own schema, and we'll use those unconditionally.
	//  For now... we'll just use whatever node you're already testing, because it oughta work
	//   (and because it avoids hardcoding any other implementation that might cause import cycle funtimes.).
	seldefn := mustNodeFromJsonString(np, str)
	sel, err := selector.ParseSelector(seldefn)
	must.NotError(err)
	return sel
}
