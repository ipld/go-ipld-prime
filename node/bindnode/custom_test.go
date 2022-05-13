package bindnode_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/node/bindnode"

	qt "github.com/frankban/quicktest"
)

// similar to cid/Cid, go-address/Address, go-graphsync/RequestID
type Boop struct{ str string }

func NewBoop(b []byte) *Boop {
	return &Boop{string(b)}
}

func (b Boop) Bytes() []byte {
	return []byte(b.str)
}

func (b Boop) String() string {
	return b.str
}

// similar to go-state-types/big/Int
type Frop struct{ *big.Int }

func NewFropFromString(str string) Frop {
	v, _ := big.NewInt(0).SetString(str, 10)
	return Frop{v}
}

func NewFropFromBytes(buf []byte) *Frop {
	var negative bool
	switch buf[0] {
	case 0:
		negative = false
	case 1:
		negative = true
	default:
		panic("can't handle this")
	}

	i := big.NewInt(0).SetBytes(buf[1:])
	if negative {
		i.Neg(i)
	}

	return &Frop{i}
}

func (b *Frop) Bytes() []byte {
	switch {
	case b.Sign() > 0:
		return append([]byte{0}, b.Int.Bytes()...)
	case b.Sign() < 0:
		return append([]byte{1}, b.Int.Bytes()...)
	default:
		return []byte{}
	}
}

type Boom struct {
	S    string
	B    Boop
	Bptr *Boop
	F    Frop
	I    int
}

const boomSchema = `
type Boom struct {
	S String
	B Bytes
	Bptr nullable Bytes
	F Bytes
	I Int
} representation map
`

const boomFixtureDagJson = `{"B":{"/":{"bytes":"dGhlc2UgYXJlIGJ5dGVz"}},"Bptr":{"/":{"bytes":"dGhlc2UgYXJlIHBvaW50ZXIgYnl0ZXM"}},"F":{"/":{"bytes":"AAH3fubjrGlwOMpClAkh/ro13L5Uls4/CtI"}},"I":10101,"S":"a string here"}`

var boomFixtureInstance = Boom{
	S:    "a string here",
	B:    *NewBoop([]byte("these are bytes")),
	Bptr: NewBoop([]byte("these are pointer bytes")),
	F:    NewFropFromString("12345678901234567891234567890123456789012345678901234567890"),
	I:    10101,
}

func BoopFromBytes(b []byte) (*Boop, error) {
	return NewBoop(b), nil
}

func BoopToBytes(boop *Boop) ([]byte, error) {
	return boop.Bytes(), nil
}

func FropFromBytes(b []byte) (*Frop, error) {
	return NewFropFromBytes(b), nil
}

func FropToBytes(frop *Frop) ([]byte, error) {
	return frop.Bytes(), nil
}

func TestCustom(t *testing.T) {
	opts := []bindnode.Option{
		bindnode.AddCustomTypeBytesConverter(BoopFromBytes, BoopToBytes),
		bindnode.AddCustomTypeBytesConverter(FropFromBytes, FropToBytes),
	}

	typeSystem, err := ipld.LoadSchemaBytes([]byte(boomSchema))
	qt.Assert(t, err, qt.IsNil)
	schemaType := typeSystem.TypeByName("Boom")
	proto := bindnode.Prototype(&Boom{}, schemaType, opts...)

	builder := proto.Representation().NewBuilder()
	err = dagjson.Decode(builder, bytes.NewReader([]byte(boomFixtureDagJson)))
	qt.Assert(t, err, qt.IsNil)

	typ := bindnode.Unwrap(builder.Build())
	inst, ok := typ.(*Boom)
	qt.Assert(t, ok, qt.IsTrue)

	cmpr := qt.CmpEquals(
		cmp.Comparer(func(x, y Boop) bool { return x.String() == y.String() }),
		cmp.Comparer(func(x, y Frop) bool { return x.String() == y.String() }),
	)
	qt.Assert(t, *inst, cmpr, boomFixtureInstance)

	tn := bindnode.Wrap(inst, schemaType, opts...)
	var buf bytes.Buffer
	err = dagjson.Encode(tn.Representation(), &buf)
	qt.Assert(t, err, qt.IsNil)

	qt.Assert(t, buf.String(), qt.Equals, boomFixtureDagJson)
}
