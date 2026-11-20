package dagjson_test

import (
	"bytes"
	"io"
	"math"
	"strings"
	"testing"

	cid "github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	plainjson "github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	rfmtjson "github.com/polydawn/refmt/json"
)

func TestFloatRoundtrip(t *testing.T) {
	for _, tc := range []struct {
		name    string
		value   float64
		encoded string
	}{
		{"zero", 0, "0.0"},
		{"negative-zero", math.Copysign(0, -1), "-0.0"},
		{"one", 1, "1.0"},
		{"negative-whole", -2, "-2.0"},
		{"location", 32, "32.0"},
		{"million", 1e6, "1000000.0"},
		{"beyond-int64", 1e20, "1.0e+20"},
		{"large-exponent", 1e21, "1e+21"},
		{"fraction", 32.1, "32.1"},
		{"small-fixed", 1e-6, "0.000001"},
		{"small-exponent", 1e-7, "1e-7"},
		{"smallest", math.SmallestNonzeroFloat64, "5e-324"},
		{"largest", math.MaxFloat64, "1.7976931348623157e+308"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, enc := range []struct {
				name   string
				encode codec.Encoder
			}{
				{"default", dagjson.Encode},
				{"options", dagjson.EncodeOptions{}.Encode},
			} {
				t.Run(enc.name, func(t *testing.T) {
					data, err := ipld.Encode(basicnode.NewFloat(tc.value), enc.encode)
					if err != nil {
						t.Fatal(err)
					}
					if string(data) != tc.encoded {
						t.Errorf("encoded %q, want %q", data, tc.encoded)
					}
					nb := basicnode.Prototype.Any.NewBuilder()
					if err := dagjson.Decode(nb, bytes.NewReader(data)); err != nil {
						t.Fatalf("decode %q: %v", data, err)
					}
					n := nb.Build()
					if n.Kind() != datamodel.Kind_Float {
						t.Fatalf("decoded %q as %s, want float", data, n.Kind())
					}
					got, err := n.AsFloat()
					if err != nil {
						t.Fatal(err)
					}
					if math.Float64bits(got) != math.Float64bits(tc.value) {
						t.Errorf("round trip changed %g to %g", tc.value, got)
					}
				})
			}
		})
	}
}

func TestFloatNestedKinds(t *testing.T) {
	input := `{"f":32.0,"i":32,"s":"32","mixed":[1.0,2,3.25,-0.0,{"1":4.0}]}`
	nb := basicnode.Prototype.Any.NewBuilder()
	if err := dagjson.Decode(nb, strings.NewReader(input)); err != nil {
		t.Fatal(err)
	}
	node := nb.Build()
	data, err := ipld.Encode(node, dagjson.Encode)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"f":32.0,"i":32,"mixed":[1.0,2,3.25,-0.0,{"1":4.0}],"s":"32"}`
	if string(data) != want {
		t.Errorf("encoded %s, want %s", data, want)
	}
	nb = basicnode.Prototype.Any.NewBuilder()
	if err := dagjson.Decode(nb, bytes.NewReader(data)); err != nil {
		t.Fatal(err)
	}
	got := nb.Build()
	if node.Length() != got.Length() {
		t.Fatal("map length changed on round trip")
	}
	for itr := node.MapIterator(); !itr.Done(); {
		key, value, err := itr.Next()
		if err != nil {
			t.Fatal(err)
		}
		name, err := key.AsString()
		if err != nil {
			t.Fatal(err)
		}
		actual, err := got.LookupByString(name)
		if err != nil {
			t.Fatal(err)
		}
		if !ipld.DeepEqual(value, actual) {
			t.Errorf("%s changed on round trip", name)
		}
	}
}

func TestFloatLinkSystemRoundtrip(t *testing.T) {
	type location struct{ Latitude, Longitude float64 }
	ts := schema.MustTypeSystem(
		schema.SpawnFloat("Float"),
		schema.SpawnStruct("Location", []schema.StructField{
			schema.SpawnStructField("Latitude", "Float", false, false),
			schema.SpawnStructField("Longitude", "Float", false, false),
		}, schema.SpawnStructRepresentationMap(nil)),
	)
	typ := ts.TypeByName("Location")
	proto := bindnode.Prototype((*location)(nil), typ).Representation()
	for _, tc := range []struct {
		name  string
		value location
	}{
		{"fractional", location{32.1, -2.25}},
		{"whole", location{32, -2}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lsys := cidlink.DefaultLinkSystem()
			var block bytes.Buffer
			lsys.StorageWriteOpener = func(linking.LinkContext) (io.Writer, linking.BlockWriteCommitter, error) {
				return &block, func(datamodel.Link) error { return nil }, nil
			}
			lsys.StorageReadOpener = func(linking.LinkContext, datamodel.Link) (io.Reader, error) {
				return bytes.NewReader(block.Bytes()), nil
			}
			lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
				Version: 1, Codec: 0x0129, MhType: 0x12, MhLength: 32,
			}}
			node := bindnode.Wrap(&tc.value, typ).Representation()
			lnk, err := lsys.Store(linking.LinkContext{}, lp, node)
			if err != nil {
				t.Fatal(err)
			}
			got, err := lsys.Load(linking.LinkContext{}, lnk, proto)
			if err != nil {
				t.Fatalf("load stored %s: %v", block.Bytes(), err)
			}
			if !ipld.DeepEqual(node, got) {
				t.Error("typed location changed on store/load")
			}
		})
	}
}

func TestFloatGenericJSONUnchanged(t *testing.T) {
	n := basicnode.NewFloat(32)
	for _, enc := range []codec.Encoder{
		plainjson.Encode,
		func(n datamodel.Node, w io.Writer) error {
			return dagjson.Marshal(n, rfmtjson.NewEncoder(w, rfmtjson.EncodeOptions{}), dagjson.EncodeOptions{})
		},
	} {
		data, err := ipld.Encode(n, enc)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "32" {
			t.Errorf("generic JSON encoded %q, want 32", data)
		}
	}
}
