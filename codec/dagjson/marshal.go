package dagjson

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/shared"
	"github.com/polydawn/refmt/tok"

	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

// This should be identical to the general feature in the parent package,
// except for the `case datamodel.Kind_Link` block,
// which is dag-json's special sauce for schemafree links.

// EncodeOptions can be used to customize the behavior of an encoding function.
// The Encode method on this struct fits the codec.Encoder function interface.
type EncodeOptions struct {
	// If true, will encode nodes with a Link kind using the DAG-JSON
	// `{"/":"cid string"}` form.
	EncodeLinks bool

	// If true, will encode nodes with a Bytes kind using the DAG-JSON
	// `{"/":{"bytes":"base64 bytes..."}}` form.
	EncodeBytes bool

	// Control the sorting of map keys, using one of the `codec.MapSortMode_*` constants.
	MapSortMode codec.MapSortMode
}

// Encode walks the given datamodel.Node and serializes it to the given io.Writer.
// Encode fits the codec.Encoder function interface.
//
// The behavior of the encoder can be customized by setting fields in the EncodeOptions struct before calling this method.
func (cfg EncodeOptions) Encode(n datamodel.Node, w io.Writer) error {
	enc := &floatEncoder{w: w}
	enc.Encoder = json.NewEncoder(enc, json.EncodeOptions{})
	return Marshal(n, enc, cfg)
}

// floatEncoder preserves the distinction between floats and integers in DAG-JSON.
// refmt emits a float as a single numeric write, but may omit its decimal point.
// Keep its formatting, adding a decimal point only when that write would look
// like an integer. The token type keeps integer values and numeric strings intact.
type floatEncoder struct {
	*json.Encoder
	w       io.Writer
	isFloat bool
	value   float64
	buffer  [64]byte
}

func (enc *floatEncoder) Step(tk *tok.Token) (bool, error) {
	enc.isFloat = tk.Type == tok.TFloat64
	enc.value = tk.Float64
	return enc.Encoder.Step(tk)
}

func (enc *floatEncoder) Write(p []byte) (int, error) {
	size := len(p)
	// A float's Step can also write array separators and whitespace.
	if enc.isFloat && size > 0 && (p[0] == '-' || p[0] >= '0' && p[0] <= '9') &&
		!bytes.ContainsAny(p, ".eE") {
		p = append(enc.buffer[:0], p...)
		if enc.value >= 1<<63 || enc.value < -(1<<63) {
			// refmt tries ParseInt before ParseFloat and can reject a large
			// integer prefix even when it is followed by a decimal point.
			p = strconv.AppendFloat(p[:0], enc.value, 'e', -1, 64)
			if bytes.IndexByte(p, '.') < 0 {
				i := bytes.IndexByte(p, 'e')
				p = append(p, 0, 0)
				copy(p[i+2:], p[i:len(p)-2])
				p[i], p[i+1] = '.', '0'
			}
		} else {
			p = append(p, '.', '0')
		}
	}
	n, err := enc.w.Write(p)
	if n < len(p) && err == nil {
		err = io.ErrShortWrite
	}
	if n == len(p) {
		n = size
	} else if n > size {
		n = size
	}
	return n, err
}

// Future work: we would like to remove the Marshal function,
// and in particular, stop seeing types from refmt (like shared.TokenSink) be visible.
// Right now, some kinds of configuration (e.g. for whitespace and prettyprint) are only available through interacting with the refmt types;
// we should improve our API so that this can be done with only our own types in this package.

// Marshal is a deprecated function.
// Please consider switching to EncodeOptions.Encode instead.
func Marshal(n datamodel.Node, sink shared.TokenSink, options EncodeOptions) error {
	var tk tok.Token
	switch n.Kind() {
	case datamodel.Kind_Invalid:
		return fmt.Errorf("cannot traverse a node that is absent")
	case datamodel.Kind_Null:
		tk.Type = tok.TNull
		_, err := sink.Step(&tk)
		return err
	case datamodel.Kind_Map:
		// Emit start of map.
		tk.Type = tok.TMapOpen
		expectedLength := int(n.Length())
		tk.Length = expectedLength // TODO: overflow check
		if _, err := sink.Step(&tk); err != nil {
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
			if len(entries) != expectedLength {
				return fmt.Errorf("map Length() does not match number of MapIterator() entries")
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
			var entryCount int
			for _, e := range entries {
				tk.Type = tok.TString
				tk.Str = e.key
				entryCount++
				if _, err := sink.Step(&tk); err != nil {
					return err
				}
				if err := Marshal(e.value, sink, options); err != nil {
					return err
				}
			}
			if entryCount != expectedLength {
				return fmt.Errorf("map Length() does not match number of MapIterator() entries")
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
	case datamodel.Kind_List:
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
	case datamodel.Kind_Bool:
		v, err := n.AsBool()
		if err != nil {
			return err
		}
		tk.Type = tok.TBool
		tk.Bool = v
		_, err = sink.Step(&tk)
		return err
	case datamodel.Kind_Int:
		v, err := n.AsInt()
		if err != nil {
			return err
		}
		tk.Type = tok.TInt
		tk.Int = int64(v)
		_, err = sink.Step(&tk)
		return err
	case datamodel.Kind_Float:
		v, err := n.AsFloat()
		if err != nil {
			return err
		}
		tk.Type = tok.TFloat64
		tk.Float64 = v
		_, err = sink.Step(&tk)
		return err
	case datamodel.Kind_String:
		v, err := n.AsString()
		if err != nil {
			return err
		}
		tk.Type = tok.TString
		tk.Str = v
		_, err = sink.Step(&tk)
		return err
	case datamodel.Kind_Bytes:
		if !options.EncodeBytes {
			return fmt.Errorf("cannot marshal IPLD bytes to this codec")
		}
		v, err := n.AsBytes()
		if err != nil {
			return err
		}
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
	case datamodel.Kind_Link:
		if !options.EncodeLinks {
			return fmt.Errorf("cannot marshal IPLD links to this codec")
		}
		v, err := n.AsLink()
		if err != nil {
			return err
		}
		switch lnk := v.(type) {
		case cidlink.Link:
			if !lnk.Cid.Defined() {
				return fmt.Errorf("encoding undefined CIDs are not supported by this codec")
			}
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
			return fmt.Errorf("schemafree link emission only supported by this codec for CID type links; got type %T", lnk)
		}
	default:
		panic("unreachable")
	}
}
