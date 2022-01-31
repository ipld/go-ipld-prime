package garbage

import (
	"bytes"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestGarbageProducesAllKinds(t *testing.T) {
	kindCount := make(map[datamodel.Kind]int)
	for i := 0; i < 1000; i++ {
		gbg := Garbage()
		kindCount[gbg.Kind()]++
	}
	for _, kind := range append(datamodel.KindSet_Scalar, datamodel.KindSet_Recursive...) {
		qt.Assert(t, kindCount[kind], qt.Not(qt.Equals), 0)
	}
}

func TestGarbageProducesValidNodes(t *testing.T) {
	// round-trip through a codec should pick up most possible problems with Node validity
	for i := 0; i < 1000; i++ {
		var buf bytes.Buffer
		gbg := Garbage()
		err := dagcbor.Encode(gbg, &buf)
		qt.Assert(t, err, qt.IsNil)
		nb := basicnode.Prototype.Any.NewBuilder()
		err = dagcbor.Decode(nb, &buf)
		qt.Assert(t, err, qt.IsNil)
		ipld.DeepEqual(gbg, nb.Build())
	}
}

func TestGarbageProducesSingleKind(t *testing.T) {
	for _, kind := range append(datamodel.KindSet_Scalar, datamodel.KindSet_Recursive...) {
		t.Run(kind.String(), func(t *testing.T) {
			kindCount := make(map[datamodel.Kind]int)
			for i := 0; i < 1000; i++ {
				gbg := Garbage(InitialWeights(map[datamodel.Kind]int{kind: 1}))
				kindCount[gbg.Kind()]++
			}
			for _, k := range append(datamodel.KindSet_Scalar, datamodel.KindSet_Recursive...) {
				if k == kind {
					qt.Assert(t, kindCount[k], qt.Equals, 1000)
				} else {
					qt.Assert(t, kindCount[k], qt.Equals, 0)
				}
			}
		})
	}
}
