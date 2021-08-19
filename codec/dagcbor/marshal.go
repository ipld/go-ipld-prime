package dagcbor

import (
	"fmt"
	"io"
	"sort"

	"github.com/polydawn/refmt/cbor"
	"github.com/polydawn/refmt/shared"
	"github.com/polydawn/refmt/tok"

	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

// This file should be identical to the general feature in the parent package,
// except for the `case datamodel.Kind_Link` block,
// which is dag-cbor's special sauce for schemafree links.

// EncodeOptions can be used to customize the behavior of an encoding function.
// The Encode method on this struct fits the codec.Encoder function interface.
type EncodeOptions struct {
	// If true, allow encoding of Link nodes as CBOR tag(42);
	// otherwise, reject them as unencodable.
	AllowLinks bool

	// Control the sorting of map keys, using one of the `codec.MapSortMode_*` constants.
	MapSortMode codec.MapSortMode
}

// Encode walks the given datamodel.Node and serializes it to the given io.Writer.
// Encode fits the codec.Encoder function interface.
//
// The behavior of the encoder can be customized by setting fields in the EncodeOptions struct before calling this method.
func (cfg EncodeOptions) Encode(n datamodel.Node, w io.Writer) error {
	// Probe for a builtin fast path.  Shortcut to that if possible.
	type detectFastPath interface {
		EncodeDagCbor(io.Writer) error
	}
	if n2, ok := n.(detectFastPath); ok {
		return n2.EncodeDagCbor(w)
	}
	// Okay, generic inspection path.
	return Marshal(n, cbor.NewEncoder(w), cfg)
}

// Future work: we would like to remove the Marshal function,
// and in particular, stop seeing types from refmt (like shared.TokenSink) be visible.
// Right now, some kinds of configuration (e.g. for whitespace and prettyprint) are only available through interacting with the refmt types;
// we should improve our API so that this can be done with only our own types in this package.

// Marshal is a deprecated function.
// Please consider switching to EncodeOptions.Encode instead.
func Marshal(n datamodel.Node, sink shared.TokenSink, options EncodeOptions) error {
	var tk tok.Token
	return marshal(n, &tk, sink, options)
}

func marshal(n datamodel.Node, tk *tok.Token, sink shared.TokenSink, options EncodeOptions) error {
	switch n.Kind() {
	case datamodel.Kind_Invalid:
		return fmt.Errorf("cannot traverse a node that is absent")
	case datamodel.Kind_Null:
		tk.Type = tok.TNull
		_, err := sink.Step(tk)
		return err
	case datamodel.Kind_Map:
		return marshalMap(n, tk, sink, options)
	case datamodel.Kind_List:
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
	case datamodel.Kind_Bool:
		v, err := n.AsBool()
		if err != nil {
			return err
		}
		tk.Type = tok.TBool
		tk.Bool = v
		_, err = sink.Step(tk)
		return err
	case datamodel.Kind_Int:
		v, err := n.AsInt()
		if err != nil {
			return err
		}
		tk.Type = tok.TInt
		tk.Int = int64(v)
		_, err = sink.Step(tk)
		return err
	case datamodel.Kind_Float:
		v, err := n.AsFloat()
		if err != nil {
			return err
		}
		tk.Type = tok.TFloat64
		tk.Float64 = v
		_, err = sink.Step(tk)
		return err
	case datamodel.Kind_String:
		v, err := n.AsString()
		if err != nil {
			return err
		}
		tk.Type = tok.TString
		tk.Str = v
		_, err = sink.Step(tk)
		return err
	case datamodel.Kind_Bytes:
		v, err := n.AsBytes()
		if err != nil {
			return err
		}
		tk.Type = tok.TBytes
		tk.Bytes = v
		_, err = sink.Step(tk)
		return err
	case datamodel.Kind_Link:
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

func marshalMap(n datamodel.Node, tk *tok.Token, sink shared.TokenSink, options EncodeOptions) error {
	// Emit start of map.
	tk.Type = tok.TMapOpen
	tk.Length = int(n.Length()) // TODO: overflow check
	if _, err := sink.Step(tk); err != nil {
		return err
	}
	if options.MapSortMode != codec.MapSortMode_None {
		// Collect map entries, then sort by key
		type entry struct {
			key   string
			value datamodel.Node
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
		// Apply the desired sort function.
		switch options.MapSortMode {
		case codec.MapSortMode_Lexical:
			sort.Slice(entries, func(i, j int) bool {
				return entries[i].key < entries[j].key
			})
		case codec.MapSortMode_RFC7049:
			sort.Slice(entries, func(i, j int) bool {
				// RFC7049 style sort as per DAG-CBOR spec
				li, lj := len(entries[i].key), len(entries[j].key)
				if li == lj {
					return entries[i].key < entries[j].key
				}
				return li < lj
			})
		}
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
	} else { // no sorting
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
			if _, err := sink.Step(tk); err != nil {
				return err
			}
			if err := marshal(v, tk, sink, options); err != nil {
				return err
			}
		}
	}
	// Emit map close.
	tk.Type = tok.TMapClose
	_, err := sink.Step(tk)
	return err
}
