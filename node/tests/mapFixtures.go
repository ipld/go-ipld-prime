package tests

import (
	"fmt"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/must"
)

var tableStrInt = [25]struct {
	s string
	i int64
}{}

func init() {
	for i := int64(1); i <= 25; i++ {
		tableStrInt[i-1] = struct {
			s string
			i int64
		}{fmt.Sprintf("k%d", i), i}
	}
}

// extracted for reuse between correctness tests and benchmarks
func buildMapStrIntN3(np datamodel.NodePrototype) datamodel.Node {
	nb := np.NewBuilder()
	ma, err := nb.BeginMap(3)
	must.NotError(err)
	must.NotError(ma.AssembleKey().AssignString("whee"))
	must.NotError(ma.AssembleValue().AssignInt(1))
	must.NotError(ma.AssembleKey().AssignString("woot"))
	must.NotError(ma.AssembleValue().AssignInt(2))
	must.NotError(ma.AssembleKey().AssignString("waga"))
	must.NotError(ma.AssembleValue().AssignInt(3))
	must.NotError(ma.Finish())
	return nb.Build()
}

// extracted for reuse across benchmarks
func buildMapStrIntN25(np datamodel.NodePrototype) datamodel.Node {
	nb := np.NewBuilder()
	ma, err := nb.BeginMap(25)
	must.NotError(err)
	for i := 1; i <= 25; i++ {
		must.NotError(ma.AssembleKey().AssignString(tableStrInt[i-1].s))
		must.NotError(ma.AssembleValue().AssignInt(tableStrInt[i-1].i))
	}
	must.NotError(ma.Finish())
	return nb.Build()
}
