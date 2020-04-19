package gengo

import (
	"testing"

	//. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestMapsContainingMaybe(t *testing.T) {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnMap("Map__String__String",
		ts.TypeByName("String"), ts.TypeByName("String"), false))
	ts.Accumulate(schema.SpawnMap("Map__String__nullableString",
		ts.TypeByName("String"), ts.TypeByName("String"), true))

	t.Run("maybe-using-embed", func(t *testing.T) {
		adjCfg.maybeUsesPtr["String"] = false

		prefix := "maps-embed"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
			//test(t, getStyleByName("Stroct"), getStyleByName("Stroct.Repr"))
		})
	})
	t.Run("maybe-using-ptr", func(t *testing.T) {
		adjCfg.maybeUsesPtr["String"] = true

		prefix := "maps-mptr"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
			//test(t, getStyleByName("Stroct"), getStyleByName("Stroct.Repr"))
		})
	})
}

func TestMapsContainingMaps(t *testing.T) {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnMap("Map__String__String",
		ts.TypeByName("String"), ts.TypeByName("String"), false))
	ts.Accumulate(schema.SpawnMap("Map__String__Map__String__String",
		ts.TypeByName("String"), ts.TypeByName("Map__String__String"), true))

	prefix := "maps-recursive"
	pkgName := "main"
	genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
		//test(t, getStyleByName("Stroct"), getStyleByName("Stroct.Repr"))
	})
}
