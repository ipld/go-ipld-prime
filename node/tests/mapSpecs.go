package tests

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/must"
)

func SpecTestMapStrInt(t *testing.T, np datamodel.NodePrototype) {
	t.Run("map<str,int>, 3 entries", func(t *testing.T) {
		n := buildMapStrIntN3(np)
		t.Run("reads back out", func(t *testing.T) {
			qt.Check(t, n.Length(), qt.Equals, int64(3))

			v, err := n.LookupByString("whee")
			qt.Check(t, err, qt.IsNil)
			v2, err := v.AsInt()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, v2, qt.Equals, int64(1))

			v, err = n.LookupByString("waga")
			qt.Check(t, err, qt.IsNil)
			v2, err = v.AsInt()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, v2, qt.Equals, int64(3))

			v, err = n.LookupByString("woot")
			qt.Check(t, err, qt.IsNil)
			v2, err = v.AsInt()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, v2, qt.Equals, int64(2))
		})
		t.Run("reads via iteration", func(t *testing.T) {
			itr := n.MapIterator()

			qt.Check(t, itr.Done(), qt.IsFalse)
			k, v, err := itr.Next()
			qt.Check(t, err, qt.IsNil)
			k2, err := k.AsString()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, k2, qt.Equals, "whee")
			v2, err := v.AsInt()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, v2, qt.Equals, int64(1))

			qt.Check(t, itr.Done(), qt.IsFalse)
			k, v, err = itr.Next()
			qt.Check(t, err, qt.IsNil)
			k2, err = k.AsString()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, k2, qt.Equals, "woot")
			v2, err = v.AsInt()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, v2, qt.Equals, int64(2))

			qt.Check(t, itr.Done(), qt.IsFalse)
			k, v, err = itr.Next()
			qt.Check(t, err, qt.IsNil)
			k2, err = k.AsString()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, k2, qt.Equals, "waga")
			v2, err = v.AsInt()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, v2, qt.Equals, int64(3))

			qt.Check(t, itr.Done(), qt.IsTrue)
			k, v, err = itr.Next()
			qt.Check(t, err, qt.Equals, datamodel.ErrIteratorOverread{})
			qt.Check(t, k, qt.IsNil)
			qt.Check(t, v, qt.IsNil)
		})
		t.Run("reads for absent keys error sensibly", func(t *testing.T) {
			v, err := n.LookupByString("nope")
			qt.Check(t, err, qt.ErrorAs, &datamodel.ErrNotExists{})
			qt.Check(t, err, qt.ErrorMatches, `key not found: "nope"`)
			qt.Check(t, v, qt.IsNil)
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
			qt.Check(t, err, qt.ErrorAs, &datamodel.ErrRepeatedMapKey{})
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
		qt.Check(t, ma.Finish(), qt.IsNil)
		n := nb.Build()
		qt.Check(t, n.Length(), qt.Equals, int64(1))
		v, err := n.LookupByString("whee")
		qt.Check(t, err, qt.IsNil)
		v2, err := v.AsInt()
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, v2, qt.Equals, int64(1))
	})
	t.Run("builder reset works", func(t *testing.T) {
		// TODO
	})
}

func SpecTestMapStrMapStrInt(t *testing.T, np datamodel.NodePrototype) {
	t.Run("map<str,map<str,int>>", func(t *testing.T) {
		nb := np.NewBuilder()
		ma, err := nb.BeginMap(3)
		must.NotError(err)
		must.NotError(ma.AssembleKey().AssignString("whee"))
		func(ma datamodel.MapAssembler, err error) {
			must.NotError(ma.AssembleKey().AssignString("m1k1"))
			must.NotError(ma.AssembleValue().AssignInt(1))
			must.NotError(ma.AssembleKey().AssignString("m1k2"))
			must.NotError(ma.AssembleValue().AssignInt(2))
			must.NotError(ma.Finish())
		}(ma.AssembleValue().BeginMap(2))
		must.NotError(ma.AssembleKey().AssignString("woot"))
		func(ma datamodel.MapAssembler, err error) {
			must.NotError(ma.AssembleKey().AssignString("m2k1"))
			must.NotError(ma.AssembleValue().AssignInt(3))
			must.NotError(ma.AssembleKey().AssignString("m2k2"))
			must.NotError(ma.AssembleValue().AssignInt(4))
			must.NotError(ma.Finish())
		}(ma.AssembleValue().BeginMap(2))
		must.NotError(ma.AssembleKey().AssignString("waga"))
		func(ma datamodel.MapAssembler, err error) {
			must.NotError(ma.AssembleKey().AssignString("m3k1"))
			must.NotError(ma.AssembleValue().AssignInt(5))
			must.NotError(ma.AssembleKey().AssignString("m3k2"))
			must.NotError(ma.AssembleValue().AssignInt(6))
			must.NotError(ma.Finish())
		}(ma.AssembleValue().BeginMap(2))
		must.NotError(ma.Finish())
		n := nb.Build()

		t.Run("reads back out", func(t *testing.T) {
			qt.Check(t, n.Length(), qt.Equals, int64(3))

			v, err := n.LookupByString("woot")
			qt.Check(t, err, qt.IsNil)
			v2, err := v.LookupByString("m2k1")
			qt.Check(t, err, qt.IsNil)
			v3, err := v2.AsInt()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, v3, qt.Equals, int64(3))
			v2, err = v.LookupByString("m2k2")
			qt.Check(t, err, qt.IsNil)
			v3, err = v2.AsInt()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, v3, qt.Equals, int64(4))
		})
	})
}

func SpecTestMapStrListStr(t *testing.T, np datamodel.NodePrototype) {
	t.Run("map<str,list<str>>", func(t *testing.T) {
		nb := np.NewBuilder()
		ma, err := nb.BeginMap(3)
		must.NotError(err)
		must.NotError(ma.AssembleKey().AssignString("asdf"))
		func(la datamodel.ListAssembler, err error) {
			must.NotError(la.AssembleValue().AssignString("eleven"))
			must.NotError(la.AssembleValue().AssignString("twelve"))
			must.NotError(la.AssembleValue().AssignString("thirteen"))
			must.NotError(la.Finish())
		}(ma.AssembleValue().BeginList(3))
		must.NotError(ma.AssembleKey().AssignString("qwer"))
		func(la datamodel.ListAssembler, err error) {
			must.NotError(la.AssembleValue().AssignString("twentyone"))
			must.NotError(la.AssembleValue().AssignString("twentytwo"))
			must.NotError(la.Finish())
		}(ma.AssembleValue().BeginList(2))
		must.NotError(ma.AssembleKey().AssignString("zxcv"))
		func(la datamodel.ListAssembler, err error) {
			must.NotError(la.AssembleValue().AssignString("thirtyone"))
			must.NotError(la.Finish())
		}(ma.AssembleValue().BeginList(1))
		must.NotError(ma.Finish())
		n := nb.Build()

		t.Run("reads back out", func(t *testing.T) {
			qt.Check(t, n.Length(), qt.Equals, int64(3))

			v, err := n.LookupByString("qwer")
			qt.Check(t, err, qt.IsNil)
			v2, err := v.LookupByIndex(1)
			qt.Check(t, err, qt.IsNil)
			v3, err := v2.AsString()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, v3, qt.Equals, "twentytwo")
		})
	})
}
