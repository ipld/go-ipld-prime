package fluent

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/ipld/go-ipld-prime"
)

func Reflect(np ipld.NodePrototype, i interface{}) (ipld.Node, error) {
	return defaultReflector.Reflect(np, i)
}

func ReflectIntoAssembler(na ipld.NodeAssembler, i interface{}) error {
	return defaultReflector.ReflectIntoAssembler(na, i)
}

var defaultReflector = Reflector{
	MapOrder: func(x, y string) bool {
		return x < y
	},
}

type Reflector struct {
	// MapOrder is used to decide a deterministic order for inserting entries to maps.
	// (This is used when converting golang maps, since their iteration order is randomized;
	// it is not used when converting other types such as structs, since those have a stable order.)
	// MapOrder should return x < y in the same way as sort.Interface.Less.
	MapOrder func(x, y string) bool
}

func (rcfg Reflector) Reflect(np ipld.NodePrototype, i interface{}) (ipld.Node, error) {
	nb := np.NewBuilder()
	if err := rcfg.ReflectIntoAssembler(nb, i); err != nil {
		return nil, err
	}
	return nb.Build(), nil
}

// ReflectIntoAssembler is a handy method for converting some basic golang types into Nodes.
//
// This plays fast and loose in general -- it's meant for demos and simple hacking, not for serious use.
// For example, in reflecting on structs, Reflect assumes no anonymous fields or other complications.
// There is no support for configuring converting struct fields with different names or other transformations.
// And so forth.
// If you need more control: this function is not what you should be using.
func (rcfg Reflector) ReflectIntoAssembler(na ipld.NodeAssembler, i interface{}) error {
	// Cover the most common values with a type-switch, as it's faster than reflection.
	switch x := i.(type) {
	case map[string]string:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Sort(sortableStrings{keys, rcfg.MapOrder})
		ma, err := na.BeginMap(len(x))
		if err != nil {
			return err
		}
		for _, k := range keys {
			va, err := ma.AssembleEntry(k)
			if err != nil {
				return err
			}
			if err := va.AssignString(x[k]); err != nil {
				return err
			}
		}
		return ma.Finish()
	case map[string]interface{}:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Sort(sortableStrings{keys, rcfg.MapOrder})
		ma, err := na.BeginMap(len(x))
		if err != nil {
			return err
		}
		for _, k := range keys {
			va, err := ma.AssembleEntry(k)
			if err != nil {
				return err
			}
			if err := rcfg.ReflectIntoAssembler(va, x[k]); err != nil {
				return err
			}
		}
		return ma.Finish()
	case []string:
		la, err := na.BeginList(len(x))
		if err != nil {
			return err
		}
		for _, v := range x {
			if err := la.AssembleValue().AssignString(v); err != nil {
				return err
			}
		}
		return la.Finish()
	case []interface{}:
		la, err := na.BeginList(len(x))
		if err != nil {
			return err
		}
		for _, v := range x {
			if err := rcfg.ReflectIntoAssembler(la.AssembleValue(), v); err != nil {
				return err
			}
		}
		return la.Finish()
	case string:
		return na.AssignString(x)
	case int:
		return na.AssignInt(x)
	case nil:
		return na.AssignNull()
	}
	// That didn't fly?  Reflection time.
	rv := reflect.ValueOf(i)
	switch rv.Kind() {
	case reflect.Bool:
		return na.AssignBool(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return na.AssignInt(int(rv.Int()))
	case reflect.Float32, reflect.Float64:
		return na.AssignFloat(rv.Float())
	case reflect.String:
		return na.AssignString(rv.String())
	case reflect.Slice, reflect.Array:
		l := rv.Len()
		la, err := na.BeginList(l)
		if err != nil {
			return err
		}
		for i := 0; i < l; i++ {
			if err := rcfg.ReflectIntoAssembler(la.AssembleValue(), rv.Index(i).Interface()); err != nil {
				return err
			}
		}
		return la.Finish()
	case reflect.Map:
		// the keys slice for sorting keeps things in reflect.Value form, because unboxing is cheap,
		//  but re-boxing is not cheap, and the MapIndex method requires reflect.Value again later.
		keys := make([]reflect.Value, 0, rv.Len())
		itr := rv.MapRange()
		for itr.Next() {
			k := itr.Key()
			if k.Kind() != reflect.String {
				return fmt.Errorf("cannot convert a map with non-string keys (%T)", i)
			}
			keys = append(keys, k)
		}
		sort.Sort(sortableReflectStrings{keys, rcfg.MapOrder})
		ma, err := na.BeginMap(rv.Len())
		if err != nil {
			return err
		}
		for _, k := range keys {
			va, err := ma.AssembleEntry(k.String())
			if err != nil {
				return err
			}
			if err := rcfg.ReflectIntoAssembler(va, rv.MapIndex(k).Interface()); err != nil {
				return err
			}
		}
		return ma.Finish()
	case reflect.Struct:
		l := rv.NumField()
		ma, err := na.BeginMap(l)
		if err != nil {
			return err
		}
		for i := 0; i < l; i++ {
			fn := rv.Type().Field(i).Name
			fv := rv.Field(i)
			va, err := ma.AssembleEntry(fn)
			if err != nil {
				return err
			}
			if err := rcfg.ReflectIntoAssembler(va, fv.Interface()); err != nil {
				return err
			}
		}
		return ma.Finish()
	case reflect.Ptr:
		if rv.IsNil() {
			return na.AssignNull()
		}
		return rcfg.ReflectIntoAssembler(na, rv.Elem())
	case reflect.Interface:
		return rcfg.ReflectIntoAssembler(na, rv.Elem())
	}
	// Some kints of values -- like Uintptr, Complex64/128, Channels, etc -- are not supported by this function.
	return fmt.Errorf("fluent.Reflect: unsure how to handle type %T (kind: %v)", i, rv.Kind())
}

type sortableStrings struct {
	a    []string
	less func(x, y string) bool
}

func (a sortableStrings) Len() int           { return len(a.a) }
func (a sortableStrings) Swap(i, j int)      { a.a[i], a.a[j] = a.a[j], a.a[i] }
func (a sortableStrings) Less(i, j int) bool { return a.less(a.a[i], a.a[j]) }

type sortableReflectStrings struct {
	a    []reflect.Value
	less func(x, y string) bool
}

func (a sortableReflectStrings) Len() int           { return len(a.a) }
func (a sortableReflectStrings) Swap(i, j int)      { a.a[i], a.a[j] = a.a[j], a.a[i] }
func (a sortableReflectStrings) Less(i, j int) bool { return a.less(a.a[i].String(), a.a[j].String()) }
