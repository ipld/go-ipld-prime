package encoding

import (
	"fmt"
	"math"

	"github.com/polydawn/refmt/shared"
	"github.com/polydawn/refmt/tok"

	"github.com/ipld/go-ipld-prime"
)

// wishlist: if we could reconstruct the ipld.Path of an error while
//  *unwinding* from that error... that'd be nice.
//   (trying to build it proactively would waste tons of allocs on the happy path.)
//  we can do this; it just requires well-typed errors and a bunch of work.

// Tests for all this are in the ipld.Node impl tests!
//  They're effectively doing double duty: testing the builders, too.
//   (Is that sensible?  Should it be refactored?  Not sure; maybe!)

// Unmarshal provides a very general tokens-to-node unmarshalling feature.
// It can handle either cbor or json by being combined with a refmt TokenSink.
//
// It is valid for all the data model types except links, which are only
// supported if the nodes are typed and provide additional information
// to clarify how links should be decoded through their type info.
// (The dag-cbor and dag-json formats can be used if links are of CID
// implementation and need to be decoded in a schemafree way.)
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
			// FUTURE: check for typed.NodeBuilder; need to specialize before recursing if so.
			// FUTURE: similar specialization needed for bind.Node as well -- perhaps this actually needs to live on NodeBuilder.
			v, err = Unmarshal(nb, tokSrc)
			if err != nil {
				return nil, err
			}
			kn, err := nb.CreateString(k)
			if err != nil {
				panic(err) // TODO: I'm no longer sure Insert should take a Node instead of string, but not recursing into reviewing that choice now.
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
		observedLen := 0
		for {
			_, err := tokSrc.Step(tk)
			if err != nil {
				return nil, err
			}
			switch tk.Type {
			case tok.TArrClose:
				if expectLen != math.MaxInt32 && observedLen != expectLen {
					return nil, fmt.Errorf("unexpected arrClose before declared length")
				}
				return lb.Build()
			default:
				observedLen++
				if observedLen > expectLen {
					return nil, fmt.Errorf("unexpected continuation of array elements beyond declared length")
				}
				// FUTURE: check for typed.NodeBuilder; need to specialize before recursing if so.
				//  N.B. when considering optionals for tuple-represented structs, keep in mind how murky that will get here.
				v, err := unmarshal(nb, tokSrc, tk)
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
		return nb.CreateBytes(tk.Bytes)
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
