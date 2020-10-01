package fluent_test

import (
	"os"

	"github.com/ipld/go-ipld-prime/fluent"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/pretty"
)

// ExampleReflect_Map demonstrates how fluent.Reflect works on maps.
// Notice that the order of keys in the IPLD map nodes is reliably sorted,
// even though golang maps provide no order guarantee.
func ExampleReflect_Map() {
	n, _ := fluent.Reflect(basicnode.Prototype.Any, map[string]interface{}{
		"k1": "fine",
		"k2": "super",
		"k3": map[string]string{
			"k31": "thanks",
			"k32": "for",
			"k33": "asking",
		},
	})
	pretty.Marshal(n, os.Stdout)

	// Output:
	// map node {
	// 	"k1": string node: "fine"
	// 	"k2": string node: "super"
	// 	"k3": map node {
	// 		"k31": string node: "thanks"
	// 		"k32": string node: "for"
	// 		"k33": string node: "asking"
	// 	}
	// }
}

// ExampleReflect_Struct demonstrates how fluent.Reflect works on structs.
// Notice that the order of keys in the IPLD map nodes retains the struct field order.
func ExampleReflect_Struct() {
	type Woo struct {
		A string
		B string
	}
	type Whee struct {
		X string
		Z string
		M Woo
	}
	n, _ := fluent.Reflect(basicnode.Prototype.Any, Whee{
		X: "fine",
		Z: "super",
		M: Woo{"thanks", "really"},
	})
	pretty.Marshal(n, os.Stdout)

	// Output:
	// map node {
	// 	"X": string node: "fine"
	// 	"Z": string node: "super"
	// 	"M": map node {
	// 		"A": string node: "thanks"
	// 		"B": string node: "really"
	// 	}
	// }
}
