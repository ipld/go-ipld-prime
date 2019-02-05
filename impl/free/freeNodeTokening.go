package ipldfree

import (
	"fmt"

	"github.com/polydawn/refmt/shared"
	"github.com/polydawn/refmt/tok"

	"github.com/ipld/go-ipld-prime"
)

// wishlist: if we could reconstruct the ipld.Path of an error while
//  *unwinding* from that error... that'd be nice.
//   (trying to build it proactively would waste tons of allocs on the happy path.)
//  we can probably do this; it just requires well-typed errors.

var (
	_ ipld.NodeUnmarshaller = Unmarshal
)

func Unmarshal(tokSrc shared.TokenSource) (ipld.Node, error) {
	var tk tok.Token
	done, err := tokSrc.Step(&tk)
	if done {
		return &Node{}, nil // invalid node, but not exactly an error
	}
	if err != nil {
		return nil, err
	}
	return unmarshal(tokSrc, &tk)
}

// starts with the first token already primed.  Necessary to get recursion
//  to flow right without a peek+unpeek system.
func unmarshal(tokSrc shared.TokenSource, tk *tok.Token) (ipld.Node, error) {
	var n Node
	switch tk.Type {
	case tok.TMapOpen:
		n.coerceType(ipld.ReprKind_Map)
		for {
			done, err := tokSrc.Step(tk)
			if done {
				return &n, fmt.Errorf("unexpected EOF")
			}
			if err != nil {
				return &n, err
			}
			switch tk.Type {
			case tok.TMapClose:
				return &n, nil
			case tok.TString:
				// continue
			default:
				return &n, fmt.Errorf("unexpected %s token while expecting map key", tk.Type)
			}
			k := tk.Str
			v, err := Unmarshal(tokSrc)
			if err != nil {
				return &n, err
			}
			if v.Kind() == ipld.ReprKind_Invalid {
				return &n, fmt.Errorf("unexpected EOF")
			}
			if _, exists := n._map[k]; exists {
				return &n, fmt.Errorf("repeated map key %q", tk)
			}
			n._mapOrd = append(n._mapOrd, k)
			n._map[k] = v
		}
	case tok.TMapClose:
		return nil, fmt.Errorf("unexpected mapClose token")
	case tok.TArrOpen:
		n.coerceType(ipld.ReprKind_List)
		for {
			done, err := tokSrc.Step(tk)
			if done {
				return &n, fmt.Errorf("unexpected EOF")
			}
			if err != nil {
				return &n, err
			}
			switch tk.Type {
			case tok.TArrClose:
				return &n, nil
			default:
				v, err := unmarshal(tokSrc, tk)
				if err != nil {
					return &n, err
				}
				if v.Kind() == ipld.ReprKind_Invalid {
					return &n, fmt.Errorf("unexpected EOF")
				}
				n._arr = append(n._arr, v)
			}
		}
		return &n, nil
	case tok.TArrClose:
		return nil, fmt.Errorf("unexpected arrClose token")
	case tok.TNull:
		n.coerceType(ipld.ReprKind_Null)
		return &n, nil
	case tok.TString:
		n.coerceType(ipld.ReprKind_String)
		n._str = tk.Str
		return &n, nil
	case tok.TBytes:
		n.coerceType(ipld.ReprKind_Bytes)
		n._bytes = tk.Bytes
		return &n, nil
		// TODO should also check tags to produce CIDs.
		//  n.b. with schemas, we can comprehend links without tags;
		//   but without schemas, tags are the only disambiguator.
	case tok.TBool:
		n.coerceType(ipld.ReprKind_Bool)
		n._bool = tk.Bool
		return &n, nil
	case tok.TInt:
		n.coerceType(ipld.ReprKind_Int)
		n._int = int(tk.Int) // FIXME overflow check
		return &n, nil
	case tok.TUint:
		n.coerceType(ipld.ReprKind_Int)
		n._int = int(tk.Uint) // FIXME overflow check
		return &n, nil
	case tok.TFloat64:
		n.coerceType(ipld.ReprKind_Float)
		n._float = tk.Float64
		return &n, nil
	default:
		panic("unreachable")
	}
}

func (n *Node) PushTokens(sink shared.TokenSink) error {
	var tk tok.Token
	switch n.kind {
	case ipld.ReprKind_Invalid:
		return fmt.Errorf("cannot traverse a node that is undefined")
	case ipld.ReprKind_Null:
		tk.Type = tok.TNull
		_, err := sink.Step(&tk)
		return err
	case ipld.ReprKind_Map:
		// Emit start of map.
		tk.Type = tok.TMapOpen
		tk.Length = len(n._map)
		if _, err := sink.Step(&tk); err != nil {
			return err
		}
		// Emit map contents (and recurse).
		for k, v := range n._map { // FIXME map order
			tk.Type = tok.TString
			tk.Str = k
			if _, err := sink.Step(&tk); err != nil {
				return err
			}
			switch v2 := v.(type) {
			case ipld.TokenizableNode:
				if err := v2.PushTokens(sink); err != nil {
					return err
				}
			default:
				panic("todo generic node tokenizer fallback")
			}
		}
		// Emit map close.
		tk.Type = tok.TMapClose
		_, err := sink.Step(&tk)
		return err
	case ipld.ReprKind_List:
		// Emit start of list.
		tk.Type = tok.TArrOpen
		tk.Length = len(n._arr)
		if _, err := sink.Step(&tk); err != nil {
			return err
		}
		// Emit list contents (and recurse).
		for _, v := range n._arr {
			switch v2 := v.(type) {
			case ipld.TokenizableNode:
				if err := v2.PushTokens(sink); err != nil {
					return err
				}
			default:
				panic("todo generic node tokenizer fallback")
			}
		}
		// Emit list close.
		tk.Type = tok.TArrClose
		_, err := sink.Step(&tk)
		return err
	case ipld.ReprKind_Bool:
		tk.Type = tok.TBool
		tk.Bool = n._bool
		_, err := sink.Step(&tk)
		return err
	case ipld.ReprKind_Int:
		tk.Type = tok.TInt
		tk.Int = int64(n._int)
		_, err := sink.Step(&tk)
		return err
	case ipld.ReprKind_Float:
		tk.Type = tok.TFloat64
		tk.Float64 = n._float
		_, err := sink.Step(&tk)
		return err
	case ipld.ReprKind_String:
		tk.Type = tok.TString
		tk.Str = n._str
		_, err := sink.Step(&tk)
		return err
	case ipld.ReprKind_Bytes:
		tk.Type = tok.TBytes
		tk.Bytes = n._bytes
		_, err := sink.Step(&tk)
		return err
	case ipld.ReprKind_Link:
		panic("todo link emission")
	default:
		panic("unreachable")
	}
}
