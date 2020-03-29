package mixins

import (
	"io"

	ipld "github.com/ipld/go-ipld-prime"
)

// kindTraitsGenerator is a embedded in all the other mixins,
// and handles all the method generation which is a pure function of the kind.
//
// OVERRIDE THE METHODS THAT DO APPLY TO YOUR KIND;
// the default method bodies produced by this mixin are those that return errors,
// and that is not what you want for the methods that *are* interesting for your kind.
// The kindTraitsGenerator methods will panic if called for a kind that should've overriden them.
//
// The other types in this package embed kindTraitsGenerator with a name,
// and only forward the methods to it that don't apply for their kind;
// this means when they're used as an anonymous embed, they grant
// all the appropriate dummy methods to their container,
// while leaving the ones that are still needed entirely absent,
// so the compiler helpfully tells you to finish rather than waiting until
// runtime to panic if a should-have-been-overriden method slips through.
type kindTraitsGenerator struct {
	PkgName    string
	TypeName   string // as will be printed in messages (e.g. can be goosed up a bit, like "Thing.Repr" instead of "_Thing__Repr").
	TypeSymbol string // the identifier in code (sometimes is munged internals like "_Thing__Repr" corresponding to no publicly admitted schema.Type.Name).
	Kind       ipld.ReprKind
}

func (g kindTraitsGenerator) emitNodeMethodLookupString(w io.Writer) {
	if ipld.ReprKindSet_JustMap.Contains(g.Kind) {
		panic("gen internals error: you should've overriden this")
	}
	doTemplate(`
		func ({{ .TypeSymbol }}) LookupString(string) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .TypeName }}", MethodName: "LookupString", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_{{ .Kind }}}
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodLookup(w io.Writer) {
	if ipld.ReprKindSet_JustMap.Contains(g.Kind) {
		panic("gen internals error: you should've overriden this")
	}
	doTemplate(`
		func ({{ .TypeSymbol }}) Lookup(ipld.Node) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .TypeName }}", MethodName: "Lookup", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_{{ .Kind }}}
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodLookupIndex(w io.Writer) {
	if ipld.ReprKindSet_JustList.Contains(g.Kind) {
		panic("gen internals error: you should've overriden this")
	}
	doTemplate(`
		func ({{ .TypeSymbol }}) LookupIndex(idx int) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .TypeName }}", MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_{{ .Kind }}}
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodLookupSegment(w io.Writer) {
	if ipld.ReprKindSet_Recursive.Contains(g.Kind) {
		panic("gen internals error: you should've overriden this")
	}
	doTemplate(`
		func ({{ .TypeSymbol }}) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .TypeName }}", MethodName: "LookupSegment", AppropriateKind: ipld.ReprKindSet_Recursive, ActualKind: ipld.ReprKind_{{ .Kind }}}
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodMapIterator(w io.Writer) {
	if ipld.ReprKindSet_JustMap.Contains(g.Kind) {
		panic("gen internals error: you should've overriden this")
	}
	doTemplate(`
		func ({{ .TypeSymbol }}) MapIterator() ipld.MapIterator {
			return nil
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodListIterator(w io.Writer) {
	if ipld.ReprKindSet_JustList.Contains(g.Kind) {
		panic("gen internals error: you should've overriden this")
	}
	doTemplate(`
		func ({{ .TypeSymbol }}) ListIterator() ipld.ListIterator {
			return nil
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodLength(w io.Writer) {
	if ipld.ReprKindSet_Recursive.Contains(g.Kind) {
		panic("gen internals error: you should've overriden this")
	}
	doTemplate(`
		func ({{ .TypeSymbol }}) Length() int {
			return -1
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodIsUndefined(w io.Writer) {
	doTemplate(`
		func ({{ .TypeSymbol }}) IsUndefined() bool {
			return false
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodIsNull(w io.Writer) {
	doTemplate(`
		func ({{ .TypeSymbol }}) IsNull() bool {
			return false
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodAsBool(w io.Writer) {
	if ipld.ReprKindSet_JustBool.Contains(g.Kind) {
		panic("gen internals error: you should've overriden this")
	}
	doTemplate(`
		func ({{ .TypeSymbol }}) AsBool() (bool, error) {
			return false, ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .TypeName }}", MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_{{ .Kind }}}
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodAsInt(w io.Writer) {
	if ipld.ReprKindSet_JustInt.Contains(g.Kind) {
		panic("gen internals error: you should've overriden this")
	}
	doTemplate(`
		func ({{ .TypeSymbol }}) AsInt() (int, error) {
			return 0, ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .TypeName }}", MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: ipld.ReprKind_{{ .Kind }}}
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodAsFloat(w io.Writer) {
	if ipld.ReprKindSet_JustFloat.Contains(g.Kind) {
		panic("gen internals error: you should've overriden this")
	}
	doTemplate(`
		func ({{ .TypeSymbol }}) AsFloat() (float64, error) {
			return 0, ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .TypeName }}", MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_{{ .Kind }}}
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodAsString(w io.Writer) {
	if ipld.ReprKindSet_JustString.Contains(g.Kind) {
		panic("gen internals error: you should've overriden this")
	}
	doTemplate(`
		func ({{ .TypeSymbol }}) AsString() (string, error) {
			return "", ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .TypeName }}", MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: ipld.ReprKind_{{ .Kind }}}
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodAsBytes(w io.Writer) {
	if ipld.ReprKindSet_JustBytes.Contains(g.Kind) {
		panic("gen internals error: you should've overriden this")
	}
	doTemplate(`
		func ({{ .TypeSymbol }}) AsBytes() ([]byte, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .TypeName }}", MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_{{ .Kind }}}
		}
	`, w, g)
}

func (g kindTraitsGenerator) emitNodeMethodAsLink(w io.Writer) {
	if ipld.ReprKindSet_JustLink.Contains(g.Kind) {
		panic("gen internals error: you should've overriden this")
	}
	doTemplate(`
		func ({{ .TypeSymbol }}) AsLink() (ipld.Link, error) {
			return nil, ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .TypeName }}", MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_{{ .Kind }}}
		}
	`, w, g)
}
