package whee

import (
	"bytes"
	"testing"

	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/tok"

	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/encoding"
	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	"github.com/ipld/go-ipld-prime/schema"
)

// TokenSourceBucket acts like a TokenSource by yielding tokens from a pre-made
// slice; and also keeps track of how far it's been read.
type TokenSourceBucket struct {
	tokens []tok.Token
	read   int
}

func (tb *TokenSourceBucket) Step(yield *tok.Token) (done bool, err error) {
	*yield = tb.tokens[tb.read]
	tb.read++
	return tb.read > len(tb.tokens), nil
}

func plz(n ipld.Node, e error) ipld.Node {
	if e != nil {
		panic(e)
	}
	return n
}

func erp(n ipld.Node, e error) interface{} {
	if e != nil {
		return e
	}
	return n
}

func TestScalarUnmarshal(t *testing.T) {
	t.Run("string node", func(t *testing.T) {
		tb := &TokenSourceBucket{tokens: []tok.Token{
			{Type: tok.TString, Str: "zooooom"},
		}}
		nb := String__NodeBuilder()
		n, err := encoding.Unmarshal(nb, tb)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_String)
		Wish(t, fluent.WrapNode(n).AsString(), ShouldEqual, "zooooom")
		Wish(t, tb.read, ShouldEqual, 1)
	})
}

// Test generated structs (n, nb, rn, rnb).
//
// Does not exercise iterators (marshal/unmarshal tests do that heavily anyway).
//
// Yes, it's big.  Proding all cases around optionals and nullables is fun.
func TestGeneratedStructs(t *testing.T) {
	t.Run("struct with map repr", func(t *testing.T) {
		var (
			v0, v1, v2, v3, v4 schema.TypedNode
		)
		t.Run("type-level build and read", func(t *testing.T) {
			t.Run("all fields set", func(t *testing.T) {
				mb, err := Stroct__NodeBuilder().CreateMap()
				Require(t, err, ShouldEqual, nil)
				mb.Insert(ipldfree.String("f1"), plz(String__NodeBuilder().CreateString("a")))
				mb.Insert(ipldfree.String("f2"), plz(String__NodeBuilder().CreateString("b")))
				mb.Insert(ipldfree.String("f3"), plz(String__NodeBuilder().CreateString("c")))
				mb.Insert(ipldfree.String("f4"), plz(String__NodeBuilder().CreateString("d")))
				n, err := mb.Build()
				v0 = n.(schema.TypedNode)

				Wish(t, err, ShouldEqual, nil)
				Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, plz(n.LookupString("f1")), ShouldEqual, plz(String__NodeBuilder().CreateString("a")))
				Wish(t, plz(n.LookupString("f2")), ShouldEqual, plz(String__NodeBuilder().CreateString("b")))
				Wish(t, plz(n.LookupString("f3")), ShouldEqual, plz(String__NodeBuilder().CreateString("c")))
				Wish(t, plz(n.LookupString("f4")), ShouldEqual, plz(String__NodeBuilder().CreateString("d")))
			})
			t.Run("using null nullable", func(t *testing.T) {
				mb, err := Stroct__NodeBuilder().CreateMap()
				Require(t, err, ShouldEqual, nil)
				mb.Insert(ipldfree.String("f1"), plz(String__NodeBuilder().CreateString("a")))
				mb.Insert(ipldfree.String("f2"), plz(String__NodeBuilder().CreateString("b")))
				mb.Insert(ipldfree.String("f3"), plz(String__NodeBuilder().CreateString("c")))
				mb.Insert(ipldfree.String("f4"), ipld.Null)
				n, err := mb.Build()
				v1 = n.(schema.TypedNode)

				Wish(t, err, ShouldEqual, nil)
				Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, 4)
				Wish(t, plz(n.LookupString("f1")), ShouldEqual, plz(String__NodeBuilder().CreateString("a")))
				Wish(t, plz(n.LookupString("f2")), ShouldEqual, plz(String__NodeBuilder().CreateString("b")))
				Wish(t, plz(n.LookupString("f3")), ShouldEqual, plz(String__NodeBuilder().CreateString("c")))
				Wish(t, plz(n.LookupString("f4")), ShouldEqual, ipld.Null)
			})
			t.Run("using null optional nullable", func(t *testing.T) {
				mb, err := Stroct__NodeBuilder().CreateMap()
				Require(t, err, ShouldEqual, nil)
				mb.Insert(ipldfree.String("f1"), plz(String__NodeBuilder().CreateString("a")))
				mb.Insert(ipldfree.String("f2"), plz(String__NodeBuilder().CreateString("b")))
				mb.Insert(ipldfree.String("f3"), ipld.Null)
				mb.Insert(ipldfree.String("f4"), plz(String__NodeBuilder().CreateString("d")))
				n, err := mb.Build()
				v2 = n.(schema.TypedNode)

				Wish(t, err, ShouldEqual, nil)
				Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, 4)
				Wish(t, plz(n.LookupString("f1")), ShouldEqual, plz(String__NodeBuilder().CreateString("a")))
				Wish(t, plz(n.LookupString("f2")), ShouldEqual, plz(String__NodeBuilder().CreateString("b")))
				Wish(t, plz(n.LookupString("f3")), ShouldEqual, ipld.Null)
				Wish(t, plz(n.LookupString("f4")), ShouldEqual, plz(String__NodeBuilder().CreateString("d")))
			})
			t.Run("using skipped optional", func(t *testing.T) {
				mb, err := Stroct__NodeBuilder().CreateMap()
				Require(t, err, ShouldEqual, nil)
				mb.Insert(ipldfree.String("f1"), plz(String__NodeBuilder().CreateString("a")))
				mb.Insert(ipldfree.String("f3"), plz(String__NodeBuilder().CreateString("c")))
				mb.Insert(ipldfree.String("f4"), plz(String__NodeBuilder().CreateString("d")))
				n, err := mb.Build()
				v3 = n.(schema.TypedNode)

				Wish(t, err, ShouldEqual, nil)
				Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, 4)
				Wish(t, plz(n.LookupString("f1")), ShouldEqual, plz(String__NodeBuilder().CreateString("a")))
				Wish(t, plz(n.LookupString("f2")), ShouldEqual, ipld.Undef)
				Wish(t, plz(n.LookupString("f3")), ShouldEqual, plz(String__NodeBuilder().CreateString("c")))
				Wish(t, plz(n.LookupString("f4")), ShouldEqual, plz(String__NodeBuilder().CreateString("d")))
			})
			t.Run("using skipped optional nullable", func(t *testing.T) {
				mb, err := Stroct__NodeBuilder().CreateMap()
				Require(t, err, ShouldEqual, nil)
				mb.Insert(ipldfree.String("f1"), plz(String__NodeBuilder().CreateString("a")))
				mb.Insert(ipldfree.String("f2"), plz(String__NodeBuilder().CreateString("b")))
				mb.Insert(ipldfree.String("f4"), plz(String__NodeBuilder().CreateString("d")))
				n, err := mb.Build()
				v4 = n.(schema.TypedNode)

				Wish(t, err, ShouldEqual, nil)
				Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, 4)
				Wish(t, plz(n.LookupString("f1")), ShouldEqual, plz(String__NodeBuilder().CreateString("a")))
				Wish(t, plz(n.LookupString("f2")), ShouldEqual, plz(String__NodeBuilder().CreateString("b")))
				Wish(t, plz(n.LookupString("f3")), ShouldEqual, ipld.Undef)
				Wish(t, plz(n.LookupString("f4")), ShouldEqual, plz(String__NodeBuilder().CreateString("d")))
			})
		})
		t.Run("representation read", func(t *testing.T) {
			t.Run("all fields set", func(t *testing.T) {
				n := v0.Representation()

				Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, 4)
				Wish(t, plz(n.LookupString("f1")), ShouldEqual, plz(String__NodeBuilder().CreateString("a")))
				Wish(t, plz(n.LookupString("f2")), ShouldEqual, plz(String__NodeBuilder().CreateString("b")))
				Wish(t, plz(n.LookupString("f3")), ShouldEqual, plz(String__NodeBuilder().CreateString("c")))
				Wish(t, plz(n.LookupString("f4")), ShouldEqual, plz(String__NodeBuilder().CreateString("d")))
			})
			t.Run("using null nullable", func(t *testing.T) {
				n := v1.Representation()

				Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, 4)
				Wish(t, plz(n.LookupString("f1")), ShouldEqual, plz(String__NodeBuilder().CreateString("a")))
				Wish(t, plz(n.LookupString("f2")), ShouldEqual, plz(String__NodeBuilder().CreateString("b")))
				Wish(t, plz(n.LookupString("f3")), ShouldEqual, plz(String__NodeBuilder().CreateString("c")))
				Wish(t, plz(n.LookupString("f4")), ShouldEqual, ipld.Null)
			})
			t.Run("using null optional nullable", func(t *testing.T) {
				n := v2.Representation()

				Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, 4)
				Wish(t, plz(n.LookupString("f1")), ShouldEqual, plz(String__NodeBuilder().CreateString("a")))
				Wish(t, plz(n.LookupString("f2")), ShouldEqual, plz(String__NodeBuilder().CreateString("b")))
				Wish(t, plz(n.LookupString("f3")), ShouldEqual, ipld.Null)
				Wish(t, plz(n.LookupString("f4")), ShouldEqual, plz(String__NodeBuilder().CreateString("d")))
			})
			t.Run("using skipped optional", func(t *testing.T) {
				n := v3.Representation()

				Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, 3) // note this is shorter, even though it's not at the type level!
				Wish(t, plz(n.LookupString("f1")), ShouldEqual, plz(String__NodeBuilder().CreateString("a")))
				Wish(t, erp(n.LookupString("f2")), ShouldEqual, ipld.ErrNotExists{ipld.PathSegmentOfString("f2")})
				Wish(t, plz(n.LookupString("f3")), ShouldEqual, plz(String__NodeBuilder().CreateString("c")))
				Wish(t, plz(n.LookupString("f4")), ShouldEqual, plz(String__NodeBuilder().CreateString("d")))
			})
			t.Run("using skipped optional nullable", func(t *testing.T) {
				n := v4.Representation()

				Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, 3) // note this is shorter, even though it's not at the type level!
				Wish(t, plz(n.LookupString("f1")), ShouldEqual, plz(String__NodeBuilder().CreateString("a")))
				Wish(t, plz(n.LookupString("f2")), ShouldEqual, plz(String__NodeBuilder().CreateString("b")))
				Wish(t, erp(n.LookupString("f3")), ShouldEqual, ipld.ErrNotExists{ipld.PathSegmentOfString("f3")})
				Wish(t, plz(n.LookupString("f4")), ShouldEqual, plz(String__NodeBuilder().CreateString("d")))
			})
			// TODO will need even more cases to probe implicits
		})
		t.Run("representation build", func(t *testing.T) {
			t.Run("all fields set", func(t *testing.T) {
				mb, err := Stroct__ReprBuilder().CreateMap()
				Require(t, err, ShouldEqual, nil)
				mb.Insert(ipldfree.String("f1"), plz(String__NodeBuilder().CreateString("a")))
				mb.Insert(ipldfree.String("f2"), plz(String__NodeBuilder().CreateString("b")))
				mb.Insert(ipldfree.String("f3"), plz(String__NodeBuilder().CreateString("c")))
				mb.Insert(ipldfree.String("f4"), plz(String__NodeBuilder().CreateString("d")))
				n, err := mb.Build()

				Wish(t, err, ShouldEqual, nil)
				Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n, ShouldEqual, v0)
			})
			t.Run("using null nullable", func(t *testing.T) {
				mb, err := Stroct__ReprBuilder().CreateMap()
				Require(t, err, ShouldEqual, nil)
				mb.Insert(ipldfree.String("f1"), plz(String__NodeBuilder().CreateString("a")))
				mb.Insert(ipldfree.String("f2"), plz(String__NodeBuilder().CreateString("b")))
				mb.Insert(ipldfree.String("f3"), plz(String__NodeBuilder().CreateString("c")))
				mb.Insert(ipldfree.String("f4"), ipld.Null)
				n, err := mb.Build()

				Wish(t, err, ShouldEqual, nil)
				Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
				Wish(t, n.Length(), ShouldEqual, 4)
				Wish(t, n, ShouldEqual, v1)
			})
			t.Run("using null optional nullable", func(t *testing.T) {
				mb, err := Stroct__ReprBuilder().CreateMap()
				Require(t, err, ShouldEqual, nil)
				mb.Insert(ipldfree.String("f1"), plz(String__NodeBuilder().CreateString("a")))
				mb.Insert(ipldfree.String("f2"), plz(String__NodeBuilder().CreateString("b")))
				mb.Insert(ipldfree.String("f3"), ipld.Null)
				mb.Insert(ipldfree.String("f4"), plz(String__NodeBuilder().CreateString("d")))
				n, err := mb.Build()

				Wish(t, err, ShouldEqual, nil)
				Wish(t, n, ShouldEqual, v2)
			})
			t.Run("using skipped optional", func(t *testing.T) {
				mb, err := Stroct__ReprBuilder().CreateMap()
				Require(t, err, ShouldEqual, nil)
				mb.Insert(ipldfree.String("f1"), plz(String__NodeBuilder().CreateString("a")))
				mb.Insert(ipldfree.String("f3"), plz(String__NodeBuilder().CreateString("c")))
				mb.Insert(ipldfree.String("f4"), plz(String__NodeBuilder().CreateString("d")))
				n, err := mb.Build()

				Wish(t, err, ShouldEqual, nil)
				Wish(t, n, ShouldEqual, v3)
			})
			t.Run("using skipped optional nullable", func(t *testing.T) {
				mb, err := Stroct__ReprBuilder().CreateMap()
				Require(t, err, ShouldEqual, nil)
				mb.Insert(ipldfree.String("f1"), plz(String__NodeBuilder().CreateString("a")))
				mb.Insert(ipldfree.String("f2"), plz(String__NodeBuilder().CreateString("b")))
				mb.Insert(ipldfree.String("f4"), plz(String__NodeBuilder().CreateString("d")))
				n, err := mb.Build()

				Wish(t, err, ShouldEqual, nil)
				Wish(t, n, ShouldEqual, v4)

			})
			// TODO will need even more cases to probe implicits
		})
	})
}

func TestStructUnmarshal(t *testing.T) {
	t.Run("stroct", func(t *testing.T) {
		t.Run("all fields set", func(t *testing.T) {
			tb := &TokenSourceBucket{tokens: []tok.Token{
				{Type: tok.TMapOpen, Length: 4},
				{Type: tok.TString, Str: "f1"}, {Type: tok.TString, Str: "a"},
				{Type: tok.TString, Str: "f2"}, {Type: tok.TString, Str: "b"},
				{Type: tok.TString, Str: "f3"}, {Type: tok.TString, Str: "c"},
				{Type: tok.TString, Str: "f4"}, {Type: tok.TString, Str: "d"},
				{Type: tok.TMapClose},
			}}
			nb := Stroct__NodeBuilder()
			n, err := encoding.Unmarshal(nb, tb)

			Require(t, err, ShouldEqual, nil)
			Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
			Wish(t, plz(n.LookupString("f1")), ShouldEqual, plz(String__NodeBuilder().CreateString("a")))
			Wish(t, plz(n.LookupString("f2")), ShouldEqual, plz(String__NodeBuilder().CreateString("b")))
			Wish(t, plz(n.LookupString("f3")), ShouldEqual, plz(String__NodeBuilder().CreateString("c")))
			Wish(t, plz(n.LookupString("f4")), ShouldEqual, plz(String__NodeBuilder().CreateString("d")))
		})
	})
}

func BenchmarkStructUnmarshal(b *testing.B) {
	bs := []byte(`{"f1":"a","f2":"b","f3":"c","f4":"d"}`)
	for i := 0; i < b.N; i++ {
		nb := Stroct__NodeBuilder()
		n, err := encoding.Unmarshal(nb, json.NewDecoder(bytes.NewReader(bs)))
		if err != nil {
			panic(err)
		}
		sink = n
	}
}

var sink interface{}
