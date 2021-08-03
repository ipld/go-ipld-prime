package dagjson

import (
	"encoding/base64"
	"fmt"
	"sort"

	"github.com/polydawn/refmt/shared"
	"github.com/polydawn/refmt/tok"

	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

// This should be identical to the general feature in the parent package,
// except for the `case ipld.Kind_Link` block,
// which is dag-json's special sauce for schemafree links.

type MarshalOptions struct {
	// If true, will encode nodes with a Link kind using the DAG-JSON
	// `{"/":"cid string"}` form.
	EncodeLinks bool

	// If true, will encode nodes with a Bytes kind using the DAG-JSON
	// `{"/":{"bytes":"base64 bytes..."}}` form.
	EncodeBytes bool

	// If true, will sort map keys prior to encoding using plain bytewise
	// comparison.
	SortMapKeys bool
}

func Marshal(n ipld.Node, sink shared.TokenSink, options MarshalOptions) error {
	var tk tok.Token
	switch n.Kind() {
	case ipld.Kind_Invalid:
		return fmt.Errorf("cannot traverse a node that is absent")
	case ipld.Kind_Null:
		tk.Type = tok.TNull
		_, err := sink.Step(&tk)
		return err
	case ipld.Kind_Map:
		// Emit start of map.
		tk.Type = tok.TMapOpen
		tk.Length = int(n.Length()) // TODO: overflow check
		if _, err := sink.Step(&tk); err != nil {
			return err
		}
		if options.SortMapKeys {
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
			sort.Slice(entries, func(i, j int) bool { return entries[i].key < entries[j].key })
			// Emit map contents (and recurse).
			for _, e := range entries {
				tk.Type = tok.TString
				tk.Str = e.key
				if _, err := sink.Step(&tk); err != nil {
					return err
				}
				if err := Marshal(e.value, sink, options); err != nil {
					return err
				}
			}
		} else {
			// Don't sort map, emit map contents (and recurse).
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
				if err := Marshal(v, sink, options); err != nil {
					return err
				}
			}
		}
		// Emit map close.
		tk.Type = tok.TMapClose
		_, err := sink.Step(&tk)
		return err
	case ipld.Kind_List:
		// Emit start of list.
		tk.Type = tok.TArrOpen
		l := n.Length()
		tk.Length = int(l) // TODO: overflow check
		if _, err := sink.Step(&tk); err != nil {
			return err
		}
		// Emit list contents (and recurse).
		for i := int64(0); i < l; i++ {
			v, err := n.LookupByIndex(i)
			if err != nil {
				return err
			}
			if err := Marshal(v, sink, options); err != nil {
				return err
			}
		}
		// Emit list close.
		tk.Type = tok.TArrClose
		_, err := sink.Step(&tk)
		return err
	case ipld.Kind_Bool:
		v, err := n.AsBool()
		if err != nil {
			return err
		}
		tk.Type = tok.TBool
		tk.Bool = v
		_, err = sink.Step(&tk)
		return err
	case ipld.Kind_Int:
		v, err := n.AsInt()
		if err != nil {
			return err
		}
		tk.Type = tok.TInt
		tk.Int = int64(v)
		_, err = sink.Step(&tk)
		return err
	case ipld.Kind_Float:
		v, err := n.AsFloat()
		if err != nil {
			return err
		}
		tk.Type = tok.TFloat64
		tk.Float64 = v
		_, err = sink.Step(&tk)
		return err
	case ipld.Kind_String:
		v, err := n.AsString()
		if err != nil {
			return err
		}
		tk.Type = tok.TString
		tk.Str = v
		_, err = sink.Step(&tk)
		return err
	case ipld.Kind_Bytes:
		v, err := n.AsBytes()
		if err != nil {
			return err
		}
		if options.EncodeBytes {
			// Precisely seven tokens to emit:
			tk.Type = tok.TMapOpen
			tk.Length = 1
			if _, err = sink.Step(&tk); err != nil {
				return err
			}
			tk.Type = tok.TString
			tk.Str = "/"
			if _, err = sink.Step(&tk); err != nil {
				return err
			}
			tk.Type = tok.TMapOpen
			tk.Length = 1
			if _, err = sink.Step(&tk); err != nil {
				return err
			}
			tk.Type = tok.TString
			tk.Str = "bytes"
			if _, err = sink.Step(&tk); err != nil {
				return err
			}
			tk.Str = base64.RawStdEncoding.EncodeToString(v)
			if _, err = sink.Step(&tk); err != nil {
				return err
			}
			tk.Type = tok.TMapClose
			if _, err = sink.Step(&tk); err != nil {
				return err
			}
			tk.Type = tok.TMapClose
			if _, err = sink.Step(&tk); err != nil {
				return err
			}
			return nil
		} else {
			tk.Type = tok.TBytes
			tk.Bytes = v
			_, err = sink.Step(&tk)
			return err
		}
	case ipld.Kind_Link:
		if !options.EncodeLinks {
			return fmt.Errorf("cannot Marshal ipld links to JSON")
		}
		v, err := n.AsLink()
		if err != nil {
			return err
		}
		switch lnk := v.(type) {
		case cidlink.Link:
			// Precisely four tokens to emit:
			tk.Type = tok.TMapOpen
			tk.Length = 1
			if _, err = sink.Step(&tk); err != nil {
				return err
			}
			tk.Type = tok.TString
			tk.Str = "/"
			if _, err = sink.Step(&tk); err != nil {
				return err
			}
			tk.Str = lnk.Cid.String()
			if _, err = sink.Step(&tk); err != nil {
				return err
			}
			tk.Type = tok.TMapClose
			if _, err = sink.Step(&tk); err != nil {
				return err
			}
			return nil
		default:
			return fmt.Errorf("schemafree link emission only supported by this codec for CID type links")
		}
	default:
		panic("unreachable")
	}
}
