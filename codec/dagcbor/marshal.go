package dagcbor

import (
	"fmt"
	"sort"

	"github.com/polydawn/refmt/shared"
	"github.com/polydawn/refmt/tok"

	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

// This file should be identical to the general feature in the parent package,
// except for the `case ipld.Kind_Link` block,
// which is dag-cbor's special sauce for schemafree links.

type MarshalOptions struct {
	// If true, allow encoding of Link nodes as CBOR tag(42), otherwise reject
	// them as unencodable
	AllowLinks bool
}

func Marshal(n ipld.Node, sink shared.TokenSink, options MarshalOptions) error {
	var tk tok.Token
	return marshal(n, &tk, sink, options)
}

func marshal(n ipld.Node, tk *tok.Token, sink shared.TokenSink, options MarshalOptions) error {
	switch n.Kind() {
	case ipld.Kind_Invalid:
		return fmt.Errorf("cannot traverse a node that is absent")
	case ipld.Kind_Null:
		tk.Type = tok.TNull
		_, err := sink.Step(tk)
		return err
	case ipld.Kind_Map:
		// Emit start of map.
		tk.Type = tok.TMapOpen
		tk.Length = int(n.Length()) // TODO: overflow check
		if _, err := sink.Step(tk); err != nil {
			return err
		}
		// Collect map entries, then sort by key
		type entry struct {
			key   string
			value ipld.Node
		}
		entries := []entry{}
		for itr := n.MapIterator(); !itr.Done(); {
			k, v, err := itr.Next()
			if err != nil {
				return err
			}
			keyStr, err := k.AsString()
			if err != nil {
				return err
			}
			entries = append(entries, entry{keyStr, v})
		}
		// RFC7049 style sort as per DAG-CBOR spec
		sort.Slice(entries, func(i, j int) bool {
			li, lj := len(entries[i].key), len(entries[j].key)
			if li == lj {
				return entries[i].key < entries[j].key
			}
			return li < lj
		})
		// Emit map contents (and recurse).
		for _, e := range entries {
			tk.Type = tok.TString
			tk.Str = e.key
			if _, err := sink.Step(tk); err != nil {
				return err
			}
			if err := marshal(e.value, tk, sink, options); err != nil {
				return err
			}
		}
		// Emit map close.
		tk.Type = tok.TMapClose
		_, err := sink.Step(tk)
		return err
	case ipld.Kind_List:
		// Emit start of list.
		tk.Type = tok.TArrOpen
		l := n.Length()
		tk.Length = int(l) // TODO: overflow check
		if _, err := sink.Step(tk); err != nil {
			return err
		}
		// Emit list contents (and recurse).
		for i := int64(0); i < l; i++ {
			v, err := n.LookupByIndex(i)
			if err != nil {
				return err
			}
			if err := marshal(v, tk, sink, options); err != nil {
				return err
			}
		}
		// Emit list close.
		tk.Type = tok.TArrClose
		_, err := sink.Step(tk)
		return err
	case ipld.Kind_Bool:
		v, err := n.AsBool()
		if err != nil {
			return err
		}
		tk.Type = tok.TBool
		tk.Bool = v
		_, err = sink.Step(tk)
		return err
	case ipld.Kind_Int:
		v, err := n.AsInt()
		if err != nil {
			return err
		}
		tk.Type = tok.TInt
		tk.Int = int64(v)
		_, err = sink.Step(tk)
		return err
	case ipld.Kind_Float:
		v, err := n.AsFloat()
		if err != nil {
			return err
		}
		tk.Type = tok.TFloat64
		tk.Float64 = v
		_, err = sink.Step(tk)
		return err
	case ipld.Kind_String:
		v, err := n.AsString()
		if err != nil {
			return err
		}
		tk.Type = tok.TString
		tk.Str = v
		_, err = sink.Step(tk)
		return err
	case ipld.Kind_Bytes:
		v, err := n.AsBytes()
		if err != nil {
			return err
		}
		tk.Type = tok.TBytes
		tk.Bytes = v
		_, err = sink.Step(tk)
		return err
	case ipld.Kind_Link:
		if !options.AllowLinks {
			return fmt.Errorf("cannot Marshal ipld links to CBOR")
		}
		v, err := n.AsLink()
		if err != nil {
			return err
		}
		switch lnk := v.(type) {
		case cidlink.Link:
			tk.Type = tok.TBytes
			tk.Bytes = append([]byte{0}, lnk.Bytes()...)
			tk.Tagged = true
			tk.Tag = linkTag
			_, err = sink.Step(tk)
			tk.Tagged = false
			return err
		default:
			return fmt.Errorf("schemafree link emission only supported by this codec for CID type links")
		}
	default:
		panic("unreachable")
	}
}
