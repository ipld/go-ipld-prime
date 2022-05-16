package bindnode_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/multiformats/go-multihash"

	qt "github.com/frankban/quicktest"
)

type BoolSubst int

const (
	BoolSubst_Yes = 100
	BoolSubst_No  = -100
)

func BoolSubstFromBool(b bool) (interface{}, error) {
	if b {
		return BoolSubst_Yes, nil
	}
	return BoolSubst_No, nil
}

func BoolToBoolSubst(b interface{}) (bool, error) {
	bp, ok := b.(*BoolSubst)
	if !ok {
		return true, fmt.Errorf("expected *BoolSubst value")
	}
	switch *bp {
	case BoolSubst_Yes:
		return true, nil
	case BoolSubst_No:
		return false, nil
	default:
		return true, fmt.Errorf("bad BoolSubst")
	}
}

type IntSubst string

func IntSubstFromInt(i int64) (interface{}, error) {
	if i == 1000 {
		return "one thousand", nil
	} else if i == 2000 {
		return "two thousand", nil
	}
	return nil, fmt.Errorf("unexpected value of IntSubst")
}

func IntToIntSubst(i interface{}) (int64, error) {
	ip, ok := i.(*IntSubst)
	if !ok {
		return 0, fmt.Errorf("expected *IntSubst value")
	}
	switch *ip {
	case "one thousand":
		return 1000, nil
	case "two thousand":
		return 2000, nil
	default:
		return 0, fmt.Errorf("bad IntSubst")
	}
}

type BigFloat struct{ *big.Float }

func BigFloatFromFloat(f float64) (interface{}, error) {
	bf := big.NewFloat(f)
	return &BigFloat{bf}, nil
}

func FloatFromBigFloat(f interface{}) (float64, error) {
	fp, ok := f.(*BigFloat)
	if !ok {
		return 0, fmt.Errorf("expected *BigFloat value")
	}
	f64, _ := fp.Float64()
	return f64, nil
}

type ByteArray [][]byte

func ByteArrayFromString(s string) (interface{}, error) {
	sa := strings.Split(s, "|")
	ba := make([][]byte, 0)
	for _, a := range sa {
		ba = append(ba, []byte(a))
	}
	return ba, nil
}

func StringFromByteArray(b interface{}) (string, error) {
	bap, ok := b.(*ByteArray)
	if !ok {
		return "", fmt.Errorf("expected *ByteArray value")
	}
	sb := strings.Builder{}
	for i, b := range *bap {
		sb.WriteString(string(b))
		if i != len(*bap)-1 {
			sb.WriteString("|")
		}
	}
	return sb.String(), nil
}

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

func BoopFromBytes(b []byte) (interface{}, error) {
	return NewBoop(b), nil
}

func BoopToBytes(iface interface{}) ([]byte, error) {
	if boop, ok := iface.(*Boop); ok {
		return boop.Bytes(), nil
	}
	return nil, fmt.Errorf("did not get expected type")
}

func FropFromBytes(b []byte) (interface{}, error) {
	return NewFropFromBytes(b), nil
}

func FropToBytes(iface interface{}) ([]byte, error) {
	if frop, ok := iface.(*Frop); ok {
		return frop.Bytes(), nil
	}
	return nil, fmt.Errorf("did not get expected type")
}

// Bitcoin's version of "links" is a hex form of the dbl-sha2-256 digest reversed
type BtcId string

func FromCidToBtcId(c cid.Cid) (interface{}, error) {
	if c.Prefix().Codec != cid.BitcoinBlock { // should be able to do BitcoinTx too .. but ..
		return nil, fmt.Errorf("can only convert IDs for BitcoinBlock codecs")
	}
	// and multihash must be dbl-sha2-256
	dig, err := multihash.Decode(c.Hash())
	if err != nil {
		return nil, err
	}
	hid := make([]byte, 0)
	for i := len(dig.Digest) - 1; i >= 0; i-- {
		hid = append(hid, dig.Digest[i])
	}
	return BtcId(hex.EncodeToString(hid)), nil
}

func FromBtcIdToCid(iface interface{}) (cid.Cid, error) {
	bid, ok := iface.(*BtcId)
	if !ok {
		return cid.Undef, fmt.Errorf("expected *BtcId value")
	}
	dig := make([]byte, 0)
	hid, err := hex.DecodeString(string(*bid))
	if err != nil {
		return cid.Undef, err
	}
	for i := len(hid) - 1; i >= 0; i-- {
		dig = append(dig, hid[i])
	}
	mh, err := multihash.Encode(dig, multihash.DBL_SHA2_256)
	if err != nil {
		return cid.Undef, err
	}
	return cid.NewCidV1(cid.BitcoinBlock, mh), nil
}

type Boom struct {
	S    string
	St   ByteArray
	B    Boop
	Bo   BoolSubst
	Bptr *Boop
	F    Frop
	Fl   BigFloat
	I    int
	In   IntSubst
	L    BtcId
}

const boomSchema = `
type Boom struct {
	S String
	St String
	B Bytes
	Bo Bool
	Bptr nullable Bytes
	F Bytes
	Fl Float
	I Int
	In Int
	L &Any
} representation map
`

const boomFixtureDagJson = `{"B":{"/":{"bytes":"dGhlc2UgYXJlIGJ5dGVz"}},"Bo":false,"Bptr":{"/":{"bytes":"dGhlc2UgYXJlIHBvaW50ZXIgYnl0ZXM"}},"F":{"/":{"bytes":"AAH3fubjrGlwOMpClAkh/ro13L5Uls4/CtI"}},"Fl":1.12,"I":10101,"In":2000,"L":{"/":"bagyacvra2e6qt2fohajauxceox55t3gedsyqap2phmv7q2qaaaaaaaaaaaaa"},"S":"a string here","St":"a|byte|array"}`

var boomFixtureInstance = Boom{
	B:    *NewBoop([]byte("these are bytes")),
	Bo:   BoolSubst_No,
	Bptr: NewBoop([]byte("these are pointer bytes")),
	F:    NewFropFromString("12345678901234567891234567890123456789012345678901234567890"),
	Fl:   BigFloat{big.NewFloat(1.12)},
	I:    10101,
	In:   IntSubst("two thousand"),
	S:    "a string here",
	St:   ByteArray([][]byte{[]byte("a"), []byte("byte"), []byte("array")}),
	L:    BtcId("00000000000000006af82b3b4f3f00b11cc4ecd9fb75445c0a1238aee8093dd1"),
}

func TestCustom(t *testing.T) {
	opts := []bindnode.Option{
		bindnode.AddCustomTypeBytesConverter(&Boop{}, BoopFromBytes, BoopToBytes),
		bindnode.AddCustomTypeBytesConverter(&Frop{}, FropFromBytes, FropToBytes),
		bindnode.AddCustomTypeBoolConverter(BoolSubst(0), BoolSubstFromBool, BoolToBoolSubst),
		bindnode.AddCustomTypeIntConverter(IntSubst(""), IntSubstFromInt, IntToIntSubst),
		bindnode.AddCustomTypeFloatConverter(&BigFloat{}, BigFloatFromFloat, FloatFromBigFloat),
		bindnode.AddCustomTypeStringConverter(&ByteArray{}, ByteArrayFromString, StringFromByteArray),
		bindnode.AddCustomTypeLinkConverter(BtcId(""), FromCidToBtcId, FromBtcIdToCid),
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
		cmp.Comparer(func(x, y BigFloat) bool { return x.String() == y.String() }),
	)
	qt.Assert(t, *inst, cmpr, boomFixtureInstance)

	tn := bindnode.Wrap(inst, schemaType, opts...)
	var buf bytes.Buffer
	err = dagjson.Encode(tn.Representation(), &buf)
	qt.Assert(t, err, qt.IsNil)

	qt.Assert(t, buf.String(), qt.Equals, boomFixtureDagJson)
}
