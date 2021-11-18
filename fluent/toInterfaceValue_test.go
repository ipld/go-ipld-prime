package fluent_test

import (
	"testing"

	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestListValue(t *testing.T) {
	a := []string{"a", "b", "c"}
	n, err := fluent.Reflect(basicnode.Prototype.Any, a)
	if err != nil {
		t.Fatal(err)
	}
	out, err := fluent.ToInterfaceValue(n)
	if err != nil {
		t.Fatal(err)
	}
	outArr := out.([]interface{})

	if len(a) != len(outArr) {
		t.Errorf("Mismatch in array size")
	}

	for i, v := range outArr {
		if a[i] != v {
			t.Errorf("expected %v, got %v at index %v", a[i], v, i)
		}
	}
}

func TestMapValue(t *testing.T) {
	a := map[string]interface{}{"a": "1", "b": int64(2), "c": 3.14}
	n, err := fluent.Reflect(basicnode.Prototype.Any, a)
	if err != nil {
		t.Fatal(err)
	}
	out, err := fluent.ToInterfaceValue(n)
	if err != nil {
		t.Fatal(err)
	}
	outM := out.(map[string]interface{})

	if len(a) != len(outM) {
		t.Errorf("Mismatch in size")
	}

	for k, v := range outM {
		if v != a[k] {
			t.Errorf("expected %v, got %v at key %v", a[k], v, k)
		}
	}

}
