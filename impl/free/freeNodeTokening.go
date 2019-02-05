package ipldfree

import (
	"fmt"

	"github.com/polydawn/refmt/shared"
	"github.com/polydawn/refmt/tok"

	"github.com/ipld/go-ipld-prime"
)

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
