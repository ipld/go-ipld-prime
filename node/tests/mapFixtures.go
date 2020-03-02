package tests

import (
	"strconv"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/must"
)

var tableStrInt = [25]struct {
	s string
	i int
}{}

func init() {
	for i := 1; i <= 25; i++ {
		tableStrInt[i-1] = struct {
			s string
			i int
		}{"k" + strconv.Itoa(i), i}
	}
}

// extracted for reuse between correctness tests and benchmarks
func buildMapStrIntN3(ns ipld.NodeStyle) ipld.Node {
	nb := ns.NewBuilder()
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
func buildMapStrIntN25(ns ipld.NodeStyle) ipld.Node {
	nb := ns.NewBuilder()
	ma, err := nb.BeginMap(25)
	must.NotError(err)
	for i := 1; i <= 25; i++ {
		must.NotError(ma.AssembleKey().AssignString(tableStrInt[i-1].s))
		must.NotError(ma.AssembleValue().AssignInt(tableStrInt[i-1].i))
	}
	must.NotError(ma.Finish())
	return nb.Build()
}
