package bindnode_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/node/bindnode"
)

func TestEnumError(t *testing.T) {
	type Action string
	const (
		ActionPresent = Action("p")
		ActionMissing = Action("m")
	)
	type S struct{ Action Action }

	schema := `
		type S struct {
			Action Action
		} representation tuple
		type Action enum {
			| Present             ("p")
			| Missing             ("m")
		} representation string
 	`

	typeSystem, err := ipld.LoadSchemaBytes([]byte(schema))
	qt.Assert(t, err, qt.IsNil)
	schemaType := typeSystem.TypeByName("S")

	node := bindnode.Wrap(&S{Action: ActionPresent}, schemaType).Representation()
	_, err = ipld.Encode(node, dagcbor.Encode)
	qt.Assert(t, err, qt.IsNotNil)
	qt.Assert(t, err.Error(), qt.Equals, `AsString: "p" is not a valid member of enum Action (bindnode works at the type level; did you mean "Present"?)`)
}
