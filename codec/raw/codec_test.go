package raw

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

var tests = []struct {
	name string
	data []byte
}{
	{"Empty", nil},
	{"Plaintext", []byte("hello there")},
	{"JSON", []byte(`{"foo": "bar"}`)},
	{"NullBytes", []byte("\x00\x00")},
}

func TestRoundtrip(t *testing.T) {
	t.Parallel()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nb := basicnode.Prototype.Bytes.NewBuilder()
			r := bytes.NewBuffer(test.data)

			err := Decode(nb, r)
			qt.Assert(t, err, qt.IsNil)
			node := nb.Build()

			buf := new(bytes.Buffer)
			err = Encode(node, buf)
			qt.Assert(t, err, qt.IsNil)

			qt.Assert(t, buf.Bytes(), qt.DeepEquals, test.data)
		})
	}
}

func TestRoundtripCidlink(t *testing.T) {
	t.Parallel()

	lb := cidlink.LinkBuilder{Prefix: cid.Prefix{
		Version:  1,
		Codec:    rawMulticodec,
		MhType:   0x17,
		MhLength: 4,
	}}
	node := basicnode.NewBytes([]byte("hello there"))

	buf := bytes.Buffer{}
	lnk, err := lb.Build(context.Background(), ipld.LinkContext{}, node,
		func(ipld.LinkContext) (io.Writer, ipld.StoreCommitter, error) {
			return &buf, func(lnk ipld.Link) error { return nil }, nil
		},
	)
	qt.Assert(t, err, qt.IsNil)

	nb := basicnode.Prototype__Any{}.NewBuilder()
	err = lnk.Load(context.Background(), ipld.LinkContext{}, nb,
		func(lnk ipld.Link, _ ipld.LinkContext) (io.Reader, error) {
			return bytes.NewReader(buf.Bytes()), nil
		},
	)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, nb.Build(), qt.DeepEquals, node)
}

// mustOnlyUseRead only exposes Read, hiding Bytes.
type mustOnlyUseRead struct {
	buf *bytes.Buffer
}

func (r mustOnlyUseRead) Read(p []byte) (int, error) {
	return r.buf.Read(p)
}

// mustNotUseRead exposes Bytes and makes Read always error.
type mustNotUseRead struct {
	buf *bytes.Buffer
}

func (r mustNotUseRead) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("must not call Read")
}

func (r mustNotUseRead) Bytes() []byte {
	return r.buf.Bytes()
}

func TestDecodeBuffer(t *testing.T) {
	t.Parallel()

	var err error
	buf := bytes.NewBuffer([]byte("hello there"))

	err = Decode(
		basicnode.Prototype.Bytes.NewBuilder(),
		mustOnlyUseRead{buf},
	)
	qt.Assert(t, err, qt.IsNil)

	err = Decode(
		basicnode.Prototype.Bytes.NewBuilder(),
		mustNotUseRead{buf},
	)
	qt.Assert(t, err, qt.IsNil)
}
