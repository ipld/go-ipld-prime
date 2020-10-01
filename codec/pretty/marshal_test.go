package pretty

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/fluent"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func Test(t *testing.T) {
	testOne := func(t *testing.T, data interface{}, expect string) {
		t.Helper()

		n, _ := fluent.Reflect(basicnode.Prototype.Any, data)
		pretty, err := MarshalToString(n)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, pretty+"\n", ShouldEqual, Dedent(expect))
	}
	t.Run("SimpleMap", func(t *testing.T) {
		testOne(t,
			map[string]interface{}{
				"asdf": "jkl",
				"werf": "vbn",
			}, `
			map node {
				"asdf": string node: "jkl"
				"werf": string node: "vbn"
			}
		`)
	})
	t.Run("RaggedMapAndList", func(t *testing.T) {
		testOne(t,
			map[string]interface{}{
				"asdf": []string{"ert", "rty"},
				"werf": map[string]interface{}{
					"zot": "zif",
				},
			}, `
			map node {
				"asdf": list node [
					0: string node: "ert"
					1: string node: "rty"
				]
				"werf": map node {
					"zot": string node: "zif"
				}
			}
		`)
	})
	t.Run("Bytes", func(t *testing.T) {
		testOne(t,
			[]byte{0x0, 0x1, 0x2},
			`
			bytes node:
				| 00000000  00 01 02                                          |...|
		`)
	})
	t.Run("MapContainingBytes", func(t *testing.T) {
		testOne(t,
			map[string]interface{}{
				"asdf": []byte{0x0, 0x1, 0x2},
				"werf": []byte{0x20, 0x21, 0x22},
			},
			`
			map node {
				"asdf": bytes node:
					| 00000000  00 01 02                                          |...|
				"werf": bytes node:
					| 00000000  20 21 22                                          | !"|
			}
		`)
	})
	t.Run("LongBytes", func(t *testing.T) {
		testOne(t,
			[]byte{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
				0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
			},
			`
			bytes node:
				| 00000000  00 01 02 03 04 05 06 07  10 11 12 13 14 15 16 17  |................|
				| 00000010  20 21 22 23 24 25 26 27                           | !"#$%&'|
		`)
	})
}
