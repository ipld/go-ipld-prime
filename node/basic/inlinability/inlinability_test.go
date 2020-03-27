// Compile with '-gcflags -S' and grep the assembly for `"".Test`.
// You'll find that both methods produce identical bodies modulo line numbers.
package inlinability

import (
	"testing"

	ipld "github.com/ipld/go-ipld-prime"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

var sink ipld.NodeBuilder

func TestStructReference(t *testing.T) {
	nb := basicnode.Style__String{}.NewBuilder()
	sink = nb
}

func TestVarReference(t *testing.T) {
	nb := basicnode.Style.String.NewBuilder()
	sink = nb
}
