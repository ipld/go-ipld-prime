package typesystem

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/impl/free"
)

func TestSimpleTypes(t *testing.T) {
	t.Run("string alone", func(t *testing.T) {
		n1 := &ipldfree.Node{}
		n1.SetString("asdf")
		t1 := TypeString{
			anyType{name: "Foo"},
		}
		Wish(t,
			Validate(Universe{}, t1, n1),
			ShouldEqual, []error(nil))
	})
}
