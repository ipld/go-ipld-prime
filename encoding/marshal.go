package encoding

import (
	"fmt"

	"github.com/polydawn/refmt/shared"
	"github.com/polydawn/refmt/tok"

	"github.com/ipld/go-ipld-prime"
)

// FUTURE there are very open questions on how to handle detection and special-track'ing for advLayout nodes when we get to that feature.

// Marshal provides a very general node-to-tokens marshalling feature.
// It can handle either cbor or json by being combined with a refmt TokenSink.
//
// It is valid for all the data model types except links, which are only
// supported if the nodes are typed and provide additional information
// to clarify how links should be encoded through their type info.
// (The dag-cbor and dag-json formats can be used if links are of CID
// implementation and need to be encoded in a schemafree way.)
func Marshal(n ipld.Node, sink shared.TokenSink) error {
	var tk tok.Token
	switch n.ReprKind() {
	case ipld.ReprKind_Invalid:
		return fmt.Errorf("cannot traverse a node that is undefined")
	case ipld.ReprKind_Null:
		tk.Type = tok.TNull
		_, err := sink.Step(&tk)
		return err
	case ipld.ReprKind_Map:
		// Emit start of map.
		tk.Type = tok.TMapOpen
		tk.Length = n.Length()
		if _, err := sink.Step(&tk); err != nil {
			return err
		}
		// Emit map contents (and recurse).
		for itr := n.MapIterator(); !itr.Done(); {
			k, v, err := itr.Next()
			if err != nil {
				return err
			}
			tk.Type = tok.TString
			tk.Str, err = k.AsString()
			if err != nil {
				return err
			}
			if _, err := sink.Step(&tk); err != nil {
				return err
			}
			if err := Marshal(v, sink); err != nil {
				return err
			}
		}
		// Emit map close.
		tk.Type = tok.TMapClose
		_, err := sink.Step(&tk)
		return err
	case ipld.ReprKind_List:
		// Emit start of list.
		tk.Type = tok.TArrOpen
		l := n.Length()
		tk.Length = l
		if _, err := sink.Step(&tk); err != nil {
			return err
		}
		// Emit list contents (and recurse).
		for i := 0; i < l; i++ {
			v, err := n.LookupIndex(i)
			if err != nil {
				return err
			}
			if err := Marshal(v, sink); err != nil {
				return err
			}
		}
		// Emit list close.
		tk.Type = tok.TArrClose
		_, err := sink.Step(&tk)
		return err
	case ipld.ReprKind_Bool:
		v, err := n.AsBool()
		if err != nil {
			return err
		}
		tk.Type = tok.TBool
		tk.Bool = v
		_, err = sink.Step(&tk)
		return err
	case ipld.ReprKind_Int:
		v, err := n.AsInt()
		if err != nil {
			return err
		}
		tk.Type = tok.TInt
		tk.Int = int64(v)
		_, err = sink.Step(&tk)
		return err
	case ipld.ReprKind_Float:
		v, err := n.AsFloat()
		if err != nil {
			return err
		}
		tk.Type = tok.TFloat64
		tk.Float64 = v
		_, err = sink.Step(&tk)
		return err
	case ipld.ReprKind_String:
		v, err := n.AsString()
		if err != nil {
			return err
		}
		tk.Type = tok.TString
		tk.Str = v
		_, err = sink.Step(&tk)
		return err
	case ipld.ReprKind_Bytes:
		v, err := n.AsBytes()
		if err != nil {
			return err
		}
		tk.Type = tok.TBytes
		tk.Bytes = v
		_, err = sink.Step(&tk)
		return err
	case ipld.ReprKind_Link:
		return fmt.Errorf("link emission not supported by this codec without a schema!  (maybe you want dag-cbor or dag-json)")
	default:
		panic("unreachable")
	}
}
