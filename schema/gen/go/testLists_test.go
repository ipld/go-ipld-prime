package gengo

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestListsContainingMaybe(t *testing.T) {
	t.Parallel()

	for _, engine := range []*genAndCompileEngine{
		{
			subtestName: "maybe-using-embed",
			prefix:      "lists-embed",
			adjCfg: AdjunctCfg{
				maybeUsesPtr: map[schema.TypeName]bool{"String": false},
			},
		},
		{
			subtestName: "maybe-using-ptr",
			prefix:      "lists-mptr",
			adjCfg: AdjunctCfg{
				maybeUsesPtr: map[schema.TypeName]bool{"String": false},
			},
		},
	} {
		t.Run(engine.subtestName, func(t *testing.T) {
			tests.SchemaTestListsContainingMaybe(t, engine)
		})
	}

}

func TestListsContainingLists(t *testing.T) {
	t.Parallel()

	engine := &genAndCompileEngine{prefix: "lists-of-lists"}
	tests.SchemaTestListsContainingLists(t, engine)
}
