package dagcbor

import (
	"errors"
	"fmt"
	"math"

	cid "github.com/ipfs/go-cid"
	"github.com/polydawn/refmt/shared"
	"github.com/polydawn/refmt/tok"

	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

var (
	ErrInvalidMultibase = errors.New("invalid multibase on IPLD link")
)

// This should be identical to the general feature in the parent package,
// except for the `case tok.TBytes` block,
// which has dag-cbor's special sauce for detecting schemafree links.

func Unmarshal(nb ipld.NodeBuilder, tokSrc shared.TokenSource) (ipld.Node, error) {
	var tk tok.Token
	done, err := tokSrc.Step(&tk)
	if err != nil {
		return nil, err
	}
	if done && !tk.Type.IsValue() {
		return nil, err
	}
	return unmarshal(nb, tokSrc, &tk)
}

// starts with the first token already primed.  Necessary to get recursion
//  to flow right without a peek+unpeek system.
func unmarshal(nb ipld.NodeBuilder, tokSrc shared.TokenSource, tk *tok.Token) (ipld.Node, error) {
	// FUTURE: check for typed.NodeBuilder that's going to parse a Link (they can slurp any token kind they want).
	switch tk.Type {
	case tok.TMapOpen:
		mb, err := nb.CreateMap()
		if err != nil {
			return nil, err
		}
		expectLen := tk.Length
		if tk.Length == -1 {
			expectLen = math.MaxInt32
		}
		observedLen := 0
		var k string
		var v ipld.Node
		for {
			_, err := tokSrc.Step(tk)
			if err != nil {
				return nil, err
			}
			switch tk.Type {
			case tok.TMapClose:
				if expectLen != math.MaxInt32 && observedLen != expectLen {
					return nil, fmt.Errorf("unexpected mapClose before declared length")
				}
				return mb.Build()
			case tok.TString:
				// continue
			default:
				return nil, fmt.Errorf("unexpected %s token while expecting map key", tk.Type)
			}
			observedLen++
			if observedLen > expectLen {
				return nil, fmt.Errorf("unexpected continuation of map elements beyond declared length")
			}
			k = tk.Str
			v, err = Unmarshal(mb.BuilderForValue(k), tokSrc)
			if err != nil {
				return nil, err
			}
			kn, err := mb.BuilderForKeys().CreateString(k)
			if err != nil {
				return nil, fmt.Errorf("value rejected as key: %s", err)
			}
			if err := mb.Insert(kn, v); err != nil {
				return nil, err
			}
		}
	case tok.TMapClose:
		return nil, fmt.Errorf("unexpected mapClose token")
	case tok.TArrOpen:
		lb, err := nb.CreateList()
		if err != nil {
			return nil, err
		}
		expectLen := tk.Length
		if tk.Length == -1 {
			expectLen = math.MaxInt32
		}
		i := 0
		for {
			_, err := tokSrc.Step(tk)
			if err != nil {
				return nil, err
			}
			switch tk.Type {
			case tok.TArrClose:
				if expectLen != math.MaxInt32 && i != expectLen {
					return nil, fmt.Errorf("unexpected arrClose before declared length")
				}
				return lb.Build()
			default:
				if i >= expectLen {
					return nil, fmt.Errorf("unexpected continuation of array elements beyond declared length")
				}
				v, err := unmarshal(lb.BuilderForValue(i), tokSrc, tk)
				i++
				if err != nil {
					return nil, err
				}
				lb.Append(v)
			}
		}
	case tok.TArrClose:
		return nil, fmt.Errorf("unexpected arrClose token")
	case tok.TNull:
		return nb.CreateNull()
	case tok.TString:
		return nb.CreateString(tk.Str)
	case tok.TBytes:
		if !tk.Tagged {
			return nb.CreateBytes(tk.Bytes)
		}
		switch tk.Tag {
		case linkTag:
			if tk.Bytes[0] != 0 {
				return nil, ErrInvalidMultibase
			}
			elCid, err := cid.Cast(tk.Bytes[1:])
			if err != nil {
				return nil, err
			}
			return nb.CreateLink(cidlink.Link{elCid})
		default:
			return nil, fmt.Errorf("unhandled cbor tag %d", tk.Tag)
		}
	case tok.TBool:
		return nb.CreateBool(tk.Bool)
	case tok.TInt:
		return nb.CreateInt(int(tk.Int)) // FIXME overflow check
	case tok.TUint:
		return nb.CreateInt(int(tk.Uint)) // FIXME overflow check
	case tok.TFloat64:
		return nb.CreateFloat(tk.Float64)
	default:
		panic("unreachable")
	}
}
