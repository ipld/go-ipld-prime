package tests

import (
	"testing"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/must"
)

func SpecBenchmarkMapStrInt_3n_AssembleStandard(b *testing.B, np datamodel.NodePrototype) {
	for i := 0; i < b.N; i++ {
		sink = buildMapStrIntN3(np)
	}
}

func SpecBenchmarkMapStrInt_3n_AssembleEntry(b *testing.B, np datamodel.NodePrototype) {
	for i := 0; i < b.N; i++ {
		nb := np.NewBuilder()
		ma, err := nb.BeginMap(3)
		if err != nil {
			panic(err)
		}
		if va, err := ma.AssembleEntry("whee"); err != nil {
			panic(err)
		} else {
			must.NotError(va.AssignInt(1))
		}
		if va, err := ma.AssembleEntry("woot"); err != nil {
			panic(err)
		} else {
			must.NotError(va.AssignInt(2))
		}
		if va, err := ma.AssembleEntry("waga"); err != nil {
			panic(err)
		} else {
			must.NotError(va.AssignInt(3))
		}
		must.NotError(ma.Finish())
		sink = nb.Build()
	}
}

func SpecBenchmarkMapStrInt_3n_Iteration(b *testing.B, np datamodel.NodePrototype) {
	n := buildMapStrIntN3(np)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		itr := n.MapIterator()
		for k, v, _ := itr.Next(); !itr.Done(); k, v, _ = itr.Next() {
			sink = k
			sink = v
		}
	}
}

// n25 -->

func SpecBenchmarkMapStrInt_25n_AssembleStandard(b *testing.B, np datamodel.NodePrototype) {
	for i := 0; i < b.N; i++ {
		sink = buildMapStrIntN25(np)
	}
}

func SpecBenchmarkMapStrInt_25n_AssembleEntry(b *testing.B, np datamodel.NodePrototype) {
	for i := 0; i < b.N; i++ {
		nb := np.NewBuilder()
		ma, err := nb.BeginMap(25)
		if err != nil {
			panic(err)
		}
		for i := 1; i <= 25; i++ {
			if va, err := ma.AssembleEntry(tableStrInt[i-1].s); err != nil {
				panic(err)
			} else {
				must.NotError(va.AssignInt(tableStrInt[i-1].i))
			}
		}
		must.NotError(ma.Finish())
		sink = nb.Build()
	}
}

func SpecBenchmarkMapStrInt_25n_Iteration(b *testing.B, np datamodel.NodePrototype) {
	n := buildMapStrIntN25(np)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		itr := n.MapIterator()
		for k, v, _ := itr.Next(); !itr.Done(); k, v, _ = itr.Next() {
			sink = k
			sink = v
		}
	}
}
