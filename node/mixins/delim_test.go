package mixins

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestSplitExact(t *testing.T) {
	type expect struct {
		value []string
		err   error
	}
	type tcase struct {
		s      string
		sep    string
		count  int
		expect expect
	}
	for _, ent := range []tcase{
		{"", "", 0, expect{[]string{}, nil}},
		{"", ":", 1, expect{[]string{""}, nil}},
		{"x", ":", 1, expect{[]string{"x"}, nil}},
		{"x:y", ":", 2, expect{[]string{"x", "y"}, nil}},
		{"x:y:", ":", 2, expect{nil, fmt.Errorf("expected 1 instances of the delimiter, found 2")}},
		{":x:y", ":", 2, expect{nil, fmt.Errorf("expected 1 instances of the delimiter, found 2")}},
		{"x:y:", ":", 3, expect{[]string{"x", "y", ""}, nil}},
	} {
		value, err := SplitExact(ent.s, ent.sep, ent.count)
		ent2 := tcase{ent.s, ent.sep, ent.count, expect{value, err}}
		qt.Check(t, ent2, qt.CmpEquals(cmp.Exporter(func(reflect.Type) bool { return true })), ent)
	}
}
