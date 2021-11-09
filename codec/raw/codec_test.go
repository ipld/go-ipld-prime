package raw

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	nodetests "github.com/ipld/go-ipld-prime/node/tests"
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

	lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,
		Codec:    rawMulticodec,
		MhType:   0x13,
		MhLength: 4,
	}}
	node := basicnode.NewBytes([]byte("hello there"))

	lsys := cidlink.DefaultLinkSystem()

	buf := bytes.Buffer{}
	lsys.StorageWriteOpener = func(lnkCtx linking.LinkContext) (io.Writer, linking.BlockWriteCommitter, error) {
		return &buf, func(lnk datamodel.Link) error { return nil }, nil
	}
	lsys.StorageReadOpener = func(lnkCtx linking.LinkContext, lnk datamodel.Link) (io.Reader, error) {
		return bytes.NewReader(buf.Bytes()), nil
	}
	lnk, err := lsys.Store(linking.LinkContext{}, lp, node)

	qt.Assert(t, err, qt.IsNil)

	newNode, err := lsys.Load(linking.LinkContext{}, lnk, basicnode.Prototype.Any)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, newNode, nodetests.NodeContentEquals, node)
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
