package tests

import (
	"testing"

	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/must"
)

func SpecTestMapStrInt(t *testing.T, np ipld.NodePrototype) {
	t.Run("map<str,int>, 3 entries", func(t *testing.T) {
		n := buildMapStrIntN3(np)
		t.Run("reads back out", func(t *testing.T) {
			Wish(t, n.Length(), ShouldEqual, int64(3))

			v, err := n.LookupByString("whee")
			Wish(t, err, ShouldEqual, nil)
			v2, err := v.AsInt()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, v2, ShouldEqual, int64(1))

			v, err = n.LookupByString("waga")
			Wish(t, err, ShouldEqual, nil)
			v2, err = v.AsInt()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, v2, ShouldEqual, int64(3))

			v, err = n.LookupByString("woot")
			Wish(t, err, ShouldEqual, nil)
			v2, err = v.AsInt()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, v2, ShouldEqual, int64(2))
		})
		t.Run("reads via iteration", func(t *testing.T) {
			itr := n.MapIterator()

			Wish(t, itr.Done(), ShouldEqual, false)
			k, v, err := itr.Next()
			Wish(t, err, ShouldEqual, nil)
			k2, err := k.AsString()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, k2, ShouldEqual, "whee")
			v2, err := v.AsInt()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, v2, ShouldEqual, int64(1))

			Wish(t, itr.Done(), ShouldEqual, false)
			k, v, err = itr.Next()
			Wish(t, err, ShouldEqual, nil)
			k2, err = k.AsString()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, k2, ShouldEqual, "woot")
			v2, err = v.AsInt()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, v2, ShouldEqual, int64(2))

			Wish(t, itr.Done(), ShouldEqual, false)
			k, v, err = itr.Next()
			Wish(t, err, ShouldEqual, nil)
			k2, err = k.AsString()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, k2, ShouldEqual, "waga")
			v2, err = v.AsInt()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, v2, ShouldEqual, int64(3))

			Wish(t, itr.Done(), ShouldEqual, true)
			k, v, err = itr.Next()
			Wish(t, err, ShouldEqual, ipld.ErrIteratorOverread{})
			Wish(t, k, ShouldEqual, nil)
			Wish(t, v, ShouldEqual, nil)
		})
		t.Run("reads for absent keys error sensibly", func(t *testing.T) {
			v, err := n.LookupByString("nope")
			Wish(t, err, ShouldBeSameTypeAs, ipld.ErrNotExists{})
			Wish(t, err.Error(), ShouldEqual, `key not found: "nope"`)
			Wish(t, v, ShouldEqual, nil)
		})
	})
	t.Run("repeated key should error", func(t *testing.T) {
		nb := np.NewBuilder()
		ma, err := nb.BeginMap(3)
		if err != nil {
			panic(err)
		}
		if err := ma.AssembleKey().AssignString("whee"); err != nil {
			panic(err)
		}
		if err := ma.AssembleValue().AssignInt(1); err != nil {
			panic(err)
		}
		if err := ma.AssembleKey().AssignString("whee"); err != nil {
			Wish(t, err, ShouldBeSameTypeAs, ipld.ErrRepeatedMapKey{})
			// No string assertion at present -- how that should be presented for typed stuff is unsettled
			//  (and if it's clever, it'll differ from untyped, which will mean no assertion possible!).
		}
	})
	t.Run("using expired child assemblers should panic", func(t *testing.T) {
		nb := np.NewBuilder()
		ma, err := nb.BeginMap(3)
		must.NotError(err)

		// Assemble a key, and then try to assign it again.  Latter should fail.
		ka := ma.AssembleKey()
		must.NotError(ka.AssignString("whee"))
		func() {
			defer func() { recover() }()
			ka.AssignString("woo")
			t.Fatal("must not be reached")
		}()

		// Assemble a value, and then try to assign it again.  Latter should fail.
		// (This does assume your system can continue after disregarding the last error.)
		va := ma.AssembleValue()
		must.NotError(va.AssignInt(1))
		func() {
			defer func() { recover() }()
			va.AssignInt(2)
			t.Fatal("must not be reached")
		}()

		// ... and neither of these should've had visible effects!
		Wish(t, ma.Finish(), ShouldEqual, nil)
		n := nb.Build()
		Wish(t, n.Length(), ShouldEqual, int64(1))
		v, err := n.LookupByString("whee")
		Wish(t, err, ShouldEqual, nil)
		v2, err := v.AsInt()
		Wish(t, err, ShouldEqual, nil)
		Wish(t, v2, ShouldEqual, int64(1))
	})
	t.Run("builder reset works", func(t *testing.T) {
		// TODO
	})
}

func SpecTestMapStrMapStrInt(t *testing.T, np ipld.NodePrototype) {
	t.Run("map<str,map<str,int>>", func(t *testing.T) {
		nb := np.NewBuilder()
		ma, err := nb.BeginMap(3)
		must.NotError(err)
		must.NotError(ma.AssembleKey().AssignString("whee"))
		func(ma ipld.MapAssembler, err error) {
			must.NotError(ma.AssembleKey().AssignString("m1k1"))
			must.NotError(ma.AssembleValue().AssignInt(1))
			must.NotError(ma.AssembleKey().AssignString("m1k2"))
			must.NotError(ma.AssembleValue().AssignInt(2))
			must.NotError(ma.Finish())
		}(ma.AssembleValue().BeginMap(2))
		must.NotError(ma.AssembleKey().AssignString("woot"))
		func(ma ipld.MapAssembler, err error) {
			must.NotError(ma.AssembleKey().AssignString("m2k1"))
			must.NotError(ma.AssembleValue().AssignInt(3))
			must.NotError(ma.AssembleKey().AssignString("m2k2"))
			must.NotError(ma.AssembleValue().AssignInt(4))
			must.NotError(ma.Finish())
		}(ma.AssembleValue().BeginMap(2))
		must.NotError(ma.AssembleKey().AssignString("waga"))
		func(ma ipld.MapAssembler, err error) {
			must.NotError(ma.AssembleKey().AssignString("m3k1"))
			must.NotError(ma.AssembleValue().AssignInt(5))
			must.NotError(ma.AssembleKey().AssignString("m3k2"))
			must.NotError(ma.AssembleValue().AssignInt(6))
			must.NotError(ma.Finish())
		}(ma.AssembleValue().BeginMap(2))
		must.NotError(ma.Finish())
		n := nb.Build()

		t.Run("reads back out", func(t *testing.T) {
			Wish(t, n.Length(), ShouldEqual, int64(3))

			v, err := n.LookupByString("woot")
			Wish(t, err, ShouldEqual, nil)
			v2, err := v.LookupByString("m2k1")
			Wish(t, err, ShouldEqual, nil)
			v3, err := v2.AsInt()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, v3, ShouldEqual, int64(3))
			v2, err = v.LookupByString("m2k2")
			Wish(t, err, ShouldEqual, nil)
			v3, err = v2.AsInt()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, v3, ShouldEqual, int64(4))
		})
	})
}

func SpecTestMapStrListStr(t *testing.T, np ipld.NodePrototype) {
	t.Run("map<str,list<str>>", func(t *testing.T) {
		nb := np.NewBuilder()
		ma, err := nb.BeginMap(3)
		must.NotError(err)
		must.NotError(ma.AssembleKey().AssignString("asdf"))
		func(la ipld.ListAssembler, err error) {
			must.NotError(la.AssembleValue().AssignString("eleven"))
			must.NotError(la.AssembleValue().AssignString("twelve"))
			must.NotError(la.AssembleValue().AssignString("thirteen"))
			must.NotError(la.Finish())
		}(ma.AssembleValue().BeginList(3))
		must.NotError(ma.AssembleKey().AssignString("qwer"))
		func(la ipld.ListAssembler, err error) {
			must.NotError(la.AssembleValue().AssignString("twentyone"))
			must.NotError(la.AssembleValue().AssignString("twentytwo"))
			must.NotError(la.Finish())
		}(ma.AssembleValue().BeginList(2))
		must.NotError(ma.AssembleKey().AssignString("zxcv"))
		func(la ipld.ListAssembler, err error) {
			must.NotError(la.AssembleValue().AssignString("thirtyone"))
			must.NotError(la.Finish())
		}(ma.AssembleValue().BeginList(1))
		must.NotError(ma.Finish())
		n := nb.Build()

		t.Run("reads back out", func(t *testing.T) {
			Wish(t, n.Length(), ShouldEqual, int64(3))

			v, err := n.LookupByString("qwer")
			Wish(t, err, ShouldEqual, nil)
			v2, err := v.LookupByIndex(1)
			Wish(t, err, ShouldEqual, nil)
			v3, err := v2.AsString()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, v3, ShouldEqual, "twentytwo")
		})
	})
}
