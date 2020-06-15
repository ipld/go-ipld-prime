package methodsets

import (
	"testing"

	. "github.com/warpfork/go-wish"
)

func TestViaUnsafe(t *testing.T) {
	x := Thing{"alpha", "beta"}

	x.Pow()
	Wish(t, x, ShouldEqual, Thing{"base", "beta"})

	x2 := FlipUnsafe(&x)
	Wish(t, x2, ShouldEqual, &Thing2ViaUnsafe{"base", "beta"})

	x2.Pow()
	Wish(t, x2, ShouldEqual, &Thing2ViaUnsafe{"unsafe", "beta"})

	Wish(t, x, ShouldEqual, Thing{"unsafe", "beta"}) // ! effects propagate back to original.

	x.Pow()
	Wish(t, x2, ShouldEqual, &Thing2ViaUnsafe{"base", "beta"}) // ! and still also vice versa.

	// it's not just that we care about retaining mutability (though that's sometimes useful);
	// it's that a 'yes' to that directly implies 'yes' to "can we get this pov without any allocations".
}
