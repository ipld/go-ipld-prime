package dagcbor

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/ipld/go-ipld-prime/datamodel"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/testutil/garbage"
)

func calculateActualLength(t *testing.T, n datamodel.Node) int64 {
	var buf bytes.Buffer
	err := Encode(n, &buf)
	qt.Assert(t, err, qt.IsNil)
	return int64(buf.Len())
}

func verifyEstimatedSize(t *testing.T, n datamodel.Node) {
	estimatedLength, err := EncodedLength(n)
	qt.Assert(t, err, qt.IsNil)
	actualLength := calculateActualLength(t, n)
	qt.Assert(t, estimatedLength, qt.Equals, actualLength)
}

func TestEncodedLength(t *testing.T) {
	t.Run("int boundaries", func(t *testing.T) {
		for ii := 0; ii < 4; ii++ {
			verifyEstimatedSize(t, basicnode.NewInt(int64(lengthBoundaries[ii].upperBound)))
			verifyEstimatedSize(t, basicnode.NewInt(int64(lengthBoundaries[ii].upperBound)-1))
			verifyEstimatedSize(t, basicnode.NewInt(int64(lengthBoundaries[ii].upperBound)+1))
			verifyEstimatedSize(t, basicnode.NewInt(-1*int64(lengthBoundaries[ii].upperBound)))
			verifyEstimatedSize(t, basicnode.NewInt(-1*int64(lengthBoundaries[ii].upperBound)-1))
			verifyEstimatedSize(t, basicnode.NewInt(-1*int64(lengthBoundaries[ii].upperBound)+1))
		}
	})

	t.Run("small garbage", func(t *testing.T) {
		seed := time.Now().Unix()
		t.Logf("randomness seed: %v\n", seed)
		rnd := rand.New(rand.NewSource(seed))
		for i := 0; i < 1000; i++ {
			gbg := garbage.Generate(rnd, garbage.TargetBlockSize(1<<6))
			verifyEstimatedSize(t, gbg)
		}
	})

	t.Run("medium garbage", func(t *testing.T) {
		seed := time.Now().Unix()
		t.Logf("randomness seed: %v\n", seed)
		rnd := rand.New(rand.NewSource(seed))
		for i := 0; i < 100; i++ {
			gbg := garbage.Generate(rnd, garbage.TargetBlockSize(1<<16))
			verifyEstimatedSize(t, gbg)
		}
	})

	t.Run("large garbage", func(t *testing.T) {
		seed := time.Now().Unix()
		t.Logf("randomness seed: %v\n", seed)
		rnd := rand.New(rand.NewSource(seed))
		for i := 0; i < 10; i++ {
			gbg := garbage.Generate(rnd, garbage.TargetBlockSize(1<<20))
			verifyEstimatedSize(t, gbg)
		}
	})
}
