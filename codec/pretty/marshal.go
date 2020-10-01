package pretty

import (
	"bytes"
	"encoding/hex"
	"io"
	"strconv"
	"strings"

	ipld "github.com/ipld/go-ipld-prime"
	codectools "github.com/ipld/go-ipld-prime/codec/tools"
	"github.com/ipld/go-ipld-prime/schema"
)

// Marshal in the "pretty" package is meant to provide a pleasing human-readable textual representation of data.
// It includes indentation, prose descriptions of the data kinds for clarity,
// and even some additional information about schema types if present.
//
// This is meant for use in debugging, pretty-printing, and examples.
// There is no unmarshaller for this format.
//
// Note that this is *not* a multicodec -- it inserts several piece of information
// in the output stream which are not pure IPLD Data Model information,
// and the presense of that information means that this function does not satisfy
// the constraints required for something to be considered a multicodec.
func Marshal(n ipld.Node, w io.Writer) error {
	return EncodeConfig{
		Line:   []byte{'\n'},
		Indent: []byte{'\t'},
		Sep:    []byte{},
	}.Marshal(n, w)
}

func MarshalToString(n ipld.Node) (string, error) {
	var buf bytes.Buffer
	err := Marshal(n, &buf)
	return buf.String(), err
}

type EncodeConfig struct {
	Line   []byte // typically: '\n'
	Indent []byte // typically: '\t'
	Sep    []byte // typically: ',' or nothing

	// the following limits can be used to configure elision:
	// (exceeding them causes "..." in the output)
	// (not yet implemented)

	DepthLimit int // every child beyond this is replaced with "...".
	CountLimit int // this counter is per depth level, not in total (use TotalLimit for that).
	TotalLimit int // this counter is the total number of elements emitted.  map start, list start, and each scalar counts as one.  map keys do not count.
}

func (cfg EncodeConfig) Marshal(n ipld.Node, w io.Writer) error {
	st := encodeState{cfg: cfg, w: w, count: make([]int, 0, 10)}
	return st.marshal(n)
}

type encodeState struct {
	cfg     EncodeConfig
	w       io.Writer
	err     error
	scratch [64]byte
	count   []int
	total   int
	// depth int // = len(count)
}

func (st *encodeState) marshal(n ipld.Node) error {
	st.total++
	if tn, ok := n.(schema.TypedNode); ok {
		st.write("(")
		st.write(string(tn.Type().Name()))
		st.write(") ")
	}
	switch n.ReprKind() {
	case ipld.ReprKind_Null:
		if n.IsAbsent() {
			st.write("absent")
		} else {
			st.write("null")
		}
		return st.checkErr()
	case ipld.ReprKind_Map:
		st.count = append(st.count, 0)
		st.write("map node {")
		st.writeBytes(st.cfg.Line)
		for itr := n.MapIterator(); !itr.Done(); {
			k, v, err := itr.Next()
			if err != nil {
				return err
			}
			st.count[len(st.count)-1]++
			st.writeIndent()
			if err := st.emitFlat(k); err != nil {
				return err
			}
			st.write(": ")
			if err := st.marshal(v); err != nil {
				return err
			}
			st.writeBytes(st.cfg.Sep)
			st.writeBytes(st.cfg.Line)
			if err := st.checkErr(); err != nil {
				return err
			}
		}
		st.count = st.count[0 : len(st.count)-1]
		st.writeIndent()
		st.write("}")
		return st.checkErr()
	case ipld.ReprKind_List:
		st.count = append(st.count, 0)
		st.write("list node [")
		st.writeBytes(st.cfg.Line)
		for itr := n.ListIterator(); !itr.Done(); {
			i, v, err := itr.Next()
			if err != nil {
				return err
			}
			st.count[len(st.count)-1]++
			st.writeIndent()
			b := strconv.AppendInt(st.scratch[:0], int64(i), 10)
			st.writeBytes(b)
			st.write(": ")
			if err := st.marshal(v); err != nil {
				return err
			}
			st.writeBytes(st.cfg.Sep)
			st.writeBytes(st.cfg.Line)
			if err := st.checkErr(); err != nil {
				return err
			}
		}
		st.count = st.count[0 : len(st.count)-1]
		st.writeIndent()
		st.write("]")
		return st.checkErr()
	case ipld.ReprKind_Bool:
		v, _ := n.AsBool()
		switch v {
		case true:
			st.write("bool node: true")
			return st.checkErr()
		case false:
			st.write("bool node: false")
			return st.checkErr()
		default:
			panic("unreachable")
		}
	case ipld.ReprKind_Int:
		v, _ := n.AsInt()
		st.write("int node: ")
		b := strconv.AppendInt(st.scratch[:0], int64(v), 10)
		st.writeBytes(b)
		return st.checkErr()
	case ipld.ReprKind_Float:
		v, _ := n.AsFloat()
		st.write("float node: ")
		b := strconv.AppendFloat(st.scratch[:0], v, 'g', -1, 64)
		st.writeBytes(b)
		return st.checkErr()
	case ipld.ReprKind_String:
		v, _ := n.AsString()
		st.write("string node: ")
		return codectools.WriteQuotedEscapedString(v, st.w)
	case ipld.ReprKind_Bytes: // we're going to use hexdump for this, which is pretty wild.  Not meant for parsing!
		v, _ := n.AsBytes()
		st.write("bytes node:")
		st.count = append(st.count, 0)
		for _, s := range strings.Split(hex.Dump(v), "\n") {
			if s == "" {
				continue
			}
			st.writeBytes(st.cfg.Line) // this will not be particularly readable if Line is nothing.
			st.writeIndent()
			st.writeBytes([]byte{'|', ' '})
			st.write(s)
		}
		st.count = st.count[0 : len(st.count)-1]
		return st.checkErr()
	case ipld.ReprKind_Link:
		v, _ := n.AsLink()
		st.write("link node: ")
		st.write(v.String())
		return st.checkErr()
	case ipld.ReprKind_Invalid:
		panic("invalid node encountered")
	default:
		panic("unreachable")
	}
}

// emitFlat is used for map keys.
// Mostly map keys are strings... but they can be structs, too, if we're dealing with schemas and they have a string representation.
// We don't want to do a full marshal recursively for this for two reasons:
//  - if it's not a string, with indentation it'd be a mess;
//  - we don't to repeat the type info for every key in a map because it's too bulky and not interesting.
func (st *encodeState) emitFlat(n ipld.Node) error {
	switch n.ReprKind() {
	case ipld.ReprKind_String:
		s, _ := n.AsString()
		return codectools.WriteQuotedEscapedString(s, st.w)
	case ipld.ReprKind_Map:
		panic("nyi")
	}
	panic("emitFlat: unhandled node kind")
}

func (st *encodeState) writeIndent() {
	for i := 0; i < len(st.count); i++ {
		st.writeBytes(st.cfg.Indent)
	}
}

func (st *encodeState) write(s string) {
	if st.err == nil {
		_, st.err = st.w.Write([]byte(s))
	}
}
func (st *encodeState) writeBytes(b []byte) {
	if st.err == nil {
		_, st.err = st.w.Write(b)
	}
}
func (st *encodeState) checkErr() error {
	return st.err
}
