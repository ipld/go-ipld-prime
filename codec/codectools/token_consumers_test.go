package codectools

import (
	"io"
	"testing"

	. "github.com/warpfork/go-wish"
)

func TestTokenAssemble(t *testing.T) {
	for _, tcase := range tokenFixtures {
		nb := tcase.value.Prototype().NewBuilder()
		var readerOffset int
		err := TokenAssemble(nb, func() (*Token, error) {
			if readerOffset > len(tcase.sequence) {
				return nil, io.EOF
			}
			readerOffset++
			return &tcase.sequence[readerOffset-1], nil
		}, 1<<10)
		if err != nil {
			t.Error(err)
		}
		Wish(t, nb.Build(), ShouldEqual, tcase.value)
	}
}

func TestTokenAssembler(t *testing.T) {
	for _, tcase := range tokenFixtures {
		nb := tcase.value.Prototype().NewBuilder()
		var ta TokenAssembler
		ta.Initialize(nb, 1<<10)
		for _, tk := range tcase.sequence {
			err := ta.Process(&tk)
			Wish(t, err, ShouldEqual, nil)
		}
		Wish(t, nb.Build(), ShouldEqual, tcase.value)
	}
}
