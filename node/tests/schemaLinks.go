package tests

import (
	"bytes"
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/storage/memstore"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/ipld/go-ipld-prime/traversal/selector"
	"github.com/ipld/go-ipld-prime/traversal/selector/builder"
)

var store = memstore.Store{}

func encode(n datamodel.Node) (datamodel.Node, datamodel.Link) {
	lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,
		Codec:    0x0129,
		MhType:   0x13,
		MhLength: 4,
	}}
	lsys := cidlink.DefaultLinkSystem()
	lsys.SetWriteStorage(&store)

	lnk, err := lsys.Store(linking.LinkContext{}, lp, n)
	if err != nil {
		panic(err)
	}
	return n, lnk
}

func SchemaTestLinks(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()
	ts.Accumulate(schema.SpawnInt("Int"))
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnList("ListOfStrings", "String", false))

	ts.Accumulate(schema.SpawnLink("Link"))                                        // &Any
	ts.Accumulate(schema.SpawnLinkReference("IntLink", "Int"))                     // &Int
	ts.Accumulate(schema.SpawnLinkReference("StringLink", "String"))               // &String
	ts.Accumulate(schema.SpawnLinkReference("ListOfStringsLink", "ListOfStrings")) // &ListOfStrings

	ts.Accumulate(schema.SpawnStruct("LinkStruct",
		[]schema.StructField{
			schema.SpawnStructField("any", "Link", false, false),
			schema.SpawnStructField("int", "IntLink", false, false),
			schema.SpawnStructField("str", "StringLink", false, false),
			schema.SpawnStructField("strlist", "ListOfStringsLink", false, false),
		},
		schema.SpawnStructRepresentationMap(map[string]string{}),
	))

	engine.Init(t, ts)

	t.Run("typed linkage traversal", func(t *testing.T) {
		_, intNodeLnk := func() (datamodel.Node, datamodel.Link) {
			np := engine.PrototypeByName("Int")
			nb := np.NewBuilder()
			nb.AssignInt(101)
			return encode(nb.Build())
		}()
		_, stringNodeLnk := encode(fluent.MustBuild(engine.PrototypeByName("String"), func(na fluent.NodeAssembler) {
			na.AssignString("a string")
		}))
		_, listOfStringsNodeLnk := encode(fluent.MustBuildList(engine.PrototypeByName("ListOfStrings"), 3, func(la fluent.ListAssembler) {
			la.AssembleValue().AssignString("s1")
			la.AssembleValue().AssignString("s2")
			la.AssembleValue().AssignString("s3")
		}))
		linkStructNode, _ := encode(fluent.MustBuildMap(engine.PrototypeByName("LinkStruct"), 4, func(ma fluent.MapAssembler) {
			ma.AssembleEntry("any").AssignLink(stringNodeLnk)
			ma.AssembleEntry("int").AssignLink(intNodeLnk)
			ma.AssembleEntry("str").AssignLink(stringNodeLnk)
			ma.AssembleEntry("strlist").AssignLink(listOfStringsNodeLnk)
		}))

		ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype.Any)
		ss := ssb.ExploreRecursive(selector.RecursionLimitDepth(3), ssb.ExploreUnion(
			ssb.Matcher(),
			ssb.ExploreAll(ssb.ExploreRecursiveEdge()),
		))
		s, err := ss.Selector()
		qt.Check(t, err, qt.IsNil)

		var order int
		lsys := cidlink.DefaultLinkSystem()
		lsys.SetReadStorage(&store)
		err = traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem: lsys,
				LinkTargetNodePrototypeChooser: func(lnk datamodel.Link, lnkCtx linking.LinkContext) (datamodel.NodePrototype, error) {
					if tlnkNd, ok := lnkCtx.LinkNode.(schema.TypedLinkNode); ok {
						return tlnkNd.LinkTargetNodePrototype(), nil
					}
					return basicnode.Prototype.Any, nil
				},
			},
		}.WalkMatching(linkStructNode, s, func(prog traversal.Progress, n datamodel.Node) error {
			buf := new(bytes.Buffer)
			dagjson.Encode(n, buf)
			fmt.Printf("Walked %d: %v\n", order, buf.String())
			switch order {
			case 0: // root
				qt.Check(t, n.Prototype(), qt.Equals, engine.PrototypeByName("LinkStruct"))
			case 1: // from an &Any
				qt.Check(t, n.Prototype(), qt.Equals, basicnode.Prototype__String{})
			case 2: // &Int
				qt.Check(t, n.Prototype(), qt.Equals, engine.PrototypeByName("Int"))
			case 3: // &String
				qt.Check(t, n.Prototype(), qt.Equals, engine.PrototypeByName("String"))
			case 4: // &ListOfStrings
				qt.Check(t, n.Prototype(), qt.Equals, engine.PrototypeByName("ListOfStrings"))
			case 5:
				fallthrough
			case 6:
				fallthrough
			case 7:
				qt.Check(t, n.Prototype(), qt.Equals, engine.PrototypeByName("String"))
			}
			order++
			return nil
		})
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, order, qt.Equals, 8)
	})
}
