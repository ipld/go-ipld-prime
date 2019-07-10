package utils

import (
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
)

// CreateNestedPathSelectorNode is a utility method to create a deep nested path selector
// from a set of strings
func CreateNestedPathSelectorNode(pathSegments []string) (ipld.Node, error) {
	fnb := fluent.WrapNodeBuilder(ipldfree.NodeBuilder())
	var selector ipld.Node
	err := fluent.Recover(func() {
		selector = createNestedPathSelectorNode(pathSegments, fnb)
	})
	return selector, err
}

func createNestedPathSelectorNode(pathSegments []string, fnb fluent.NodeBuilder) ipld.Node {
	if len(pathSegments) == 0 {
		return fnb.CreateBool(true)
	}
	return fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
		mb.Insert(
			knb.CreateString("f"),
			vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(pathSegments[0]), createNestedPathSelectorNode(pathSegments[1:], vnb))
			}),
		)
	})
}
