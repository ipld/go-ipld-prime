package rot13adl_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/polydawn/refmt/json"

	"github.com/ipld/go-ipld-prime/adl/rot13adl"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/must"
)

func ExampleUnmarshallingToADL() {
	// Create a NodeBuilder for the ADL's substrate.
	//  Unmarshalling into this memory structure is optimal,
	//   because it immediately puts data into the right memory layout for the ADL code to work on,
	//  but you could use any other kind of NodeBuilder just as well and still get correct results.
	nb := rot13adl.Prototype.SubstrateRoot.NewBuilder()

	// Unmarshal -- using the substrate's nodebuilder just like you'd unmarshal with any other nodebuilder.
	err := dagjson.Unmarshal(nb, json.NewDecoder(strings.NewReader(`"n pbby fgevat"`)))
	fmt.Printf("unmarshal error: %v\n", err)

	// Use `Reify` to get the synthetic high-level view of the ADL data.
	substrateNode := nb.Build()
	syntheticView, err := rot13adl.Reify(substrateNode)
	fmt.Printf("reify error: %v\n", err)

	// We can inspect the synthetic ADL node like any other node!
	fmt.Printf("adl node kind: %v\n", syntheticView.Kind())
	fmt.Printf("adl view value: %q\n", must.String(syntheticView))

	// Output:
	// unmarshal error: <nil>
	// reify error: <nil>
	// adl node kind: string
	// adl view value: "a cool string"
}

func ExampleCreatingViaADL() {
	// Create a NodeBuilder for the ADL -- the high-level synthesized thing (not the substrate).
	nb := rot13adl.Prototype.Node.NewBuilder()

	// Create a ADL node via its builder.  This is just like creating any other node in IPLD.
	nb.AssignString("woohoo")
	n := nb.Build()

	// We can inspect the synthetic ADL node like any other node!
	fmt.Printf("adl node kind: %v\n", n.Kind())
	fmt.Printf("adl view value: %q\n", must.String(n))

	// We can get the substrate view and examine that as a node too.
	// (This requires a cast to see that we have an ADL, though.  Not all IPLD nodes have a 'Substrate' property.)
	substrateNode := n.(rot13adl.R13String).Substrate()
	fmt.Printf("substrate node kind: %v\n", substrateNode.Kind())
	fmt.Printf("substrate value: %q\n", must.String(substrateNode))

	// To marshal the ADL, just use marshal methods on its substrate as normal:
	var marshalBuffer bytes.Buffer
	err := dagjson.Marshal(substrateNode, json.NewEncoder(&marshalBuffer, json.EncodeOptions{}), true)
	fmt.Printf("marshalled: %v\n", marshalBuffer.String())
	fmt.Printf("marshal error: %v\n", err)

	// Output:
	// adl node kind: string
	// adl view value: "woohoo"
	// substrate node kind: string
	// substrate value: "jbbubb"
	// marshalled: "jbbubb"
	// marshal error: <nil>
}

// It's worth noting that the builders for an ADL substrate node still return the substrate.
// (This is interesting in contrast to Schemas, where codegenerated representation-level builders
// yield the type-level node values (and not the representation level node).)
//
// To convert the substrate node to the high level synthesized view of the ADL,
// use Reify as normal -- it's the same whether you've used the substrate type
// or if you've used any other node implementation to hold the data.
//

// Future work: unmarshalling which can invoke an ADL mid-structure,
// and automatically places the reified ADL in place in the larger structure.
//
// There will be several ways to do this (it hinges around "the signalling problem",
// discussed in https://github.com/ipld/specs/issues/130 ):
//
// The first way is to use IPLD Schemas, which provide a signalling mechanism
// by leaning on the schema, and the matching of shape of surrounding data to the schema,
// as a way to determine where an ADL is expected to appear.
//
// A second mechanism could involve new unmarshal function contracts
// which would ake a (fairly complex) argument that says what NodePrototype to use in certain positions.
// This could be accomplished by use of Selectors.
// (This would also have many other potential purposes -- implementing this in terms of NodePrototype selection is very multi-purpose,
// and could be used for efficiency and misc tuning purposes,
// for expecting a *schema* thing part way through, and so forth.)
//
