package tests

import (
	"bytes"

	refmtjson "github.com/polydawn/refmt/json"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/encoding"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

func mustNodeFromJsonString(nb ipld.NodeBuilder, str string) ipld.Node {
	return must.Node(encoding.Unmarshal(nb, refmtjson.NewDecoder(bytes.NewBufferString(str))))
}
func mustSelectorFromJsonString(nb ipld.NodeBuilder, str string) selector.Selector {
	// Needing an 'nb' parameter here is sort of off-topic, to be honest.
	//  Someday the selector package will probably contain codegen'd nodes of its own schema, and we'll use those unconditionally.
	//  For now... we'll just use whatever node you're already testing, because it oughta work
	//   (and because it avoids hardcoding any other implementation that might cause import cycle funtimes.).
	seldefn := mustNodeFromJsonString(nb, str)
	sel, err := selector.ParseSelector(seldefn)
	must.NotError(err)
	return sel
}
