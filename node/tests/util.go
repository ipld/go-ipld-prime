package tests

import (
	"strings"

	"github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

func mustNodeFromJsonString(np datamodel.NodePrototype, str string) datamodel.Node {
	nb := np.NewBuilder()
	must.NotError(json.Decode(nb, strings.NewReader(str)))
	return nb.Build()
}

func mustSelectorFromJsonString(np datamodel.NodePrototype, str string) selector.Selector {
	// Needing an 'ns' parameter here is sort of off-topic, to be honest.
	//  Someday the selector package will probably contain codegen'd nodes of its own schema, and we'll use those unconditionally.
	//  For now... we'll just use whatever node you're already testing, because it oughta work
	//   (and because it avoids hardcoding any other implementation that might cause import cycle funtimes.).
	seldefn := mustNodeFromJsonString(np, str)
	sel, err := selector.ParseSelector(seldefn)
	must.NotError(err)
	return sel
}
