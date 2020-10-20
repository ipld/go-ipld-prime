package codectools

import (
	"io"
	"testing"

	. "github.com/warpfork/go-wish"
)

func TestTokenWalk(t *testing.T) {
	for _, tcase := range tokenFixtures {
		var result []Token
		err := TokenWalk(tcase.value, func(tk *Token) error {
			result = append(result, *tk)
			return nil
		})
		if err != nil {
			t.Error(err)
		}
		Wish(t, stringifyTokens(result), ShouldEqual, stringifyTokens(tcase.sequence))
	}
}

func TestNodeTokenizer(t *testing.T) {
	for _, tcase := range tokenFixtures {
		var nt NodeTokenizer
		var result []Token
		nt.Initialize(tcase.value)
		for {
			tk, err := nt.ReadToken()
			if err == nil {
				result = append(result, *tk)
			} else if err == io.EOF {
				break
			} else {
				t.Error(err)
				break
			}
		}
		Wish(t, stringifyTokens(result), ShouldEqual, stringifyTokens(tcase.sequence))
	}
}
