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
			Name: "Foo",
		}
		Wish(t,
			Validate(nil, t1, n1),
			ShouldEqual, []error(nil))
	})
}
