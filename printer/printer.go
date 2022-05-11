package printer

import (
	"encoding/hex"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
)

// Print emits a textual description of the node tree straight to stdout.
// All printer configuration will be the default;
// links will be printed, and will not be traversed.
func Print(n datamodel.Node) {
	Config{}.Print(n)
}

// Sprint returns a textual description of the node tree.
// All printer configuration will be the default;
// links will be printed, and will not be traversed.
func Sprint(n datamodel.Node) string {
	return Config{}.Sprint(n)
}

// Fprint accepts an io.Writer to which a textual description of the node tree will be written.
// All printer configuration will be the default;
// links will be printed, and will not be traversed.
func Fprint(w io.Writer, n datamodel.Node) {
	Config{}.Fprint(w, n)
}

// Print emits a textual description of the node tree straight to stdout.
// The configuration structure this method is attached to can be used to specified details for how the printout will be formatted.
func (cfg Config) Print(n datamodel.Node) {
	cfg.Fprint(os.Stdout, n)
}

// Sprint returns a textual description of the node tree.
// The configuration structure this method is attached to can be used to specified details for how the printout will be formatted.
func (cfg Config) Sprint(n datamodel.Node) string {
	var buf strings.Builder
	cfg.Fprint(&buf, n)
	return buf.String()
}

// Fprint accepts an io.Writer to which a textual description of the node tree will be written.
// The configuration structure this method is attached to can be used to specified details for how the printout will be formatted.
func (cfg Config) Fprint(w io.Writer, n datamodel.Node) {
	pr := printBuf{w, cfg}
	pr.Config.init()
	pr.doString(0, printState_normal, n)
}

type Config struct {
	// If true, long strings and long byte sequences will truncated, and will include ellipses instead.
	//
	// Not yet supported.
	Abbreviate bool

	// If set, the indentation to use.
	// If nil, it will be treated as a default "\t".
	Indentation []byte

	// Probably does exactly what you think it does.
	StartingIndent []byte

	// Set to true if you like verbosity, I guess.
	// If false, strings will only have kind+type markings if they're typed.
	//
	// Not yet supported.
	AlwaysMarkStrings bool

	// Set to true if you want type info to be skipped for any type that's in the Prelude
	// (e.g. instead of `string<String>{` seeing only `string{` is preferred, etc).
	//
	// Not yet supported.
	ElidePreludeTypeInfo bool

	// Set to true if you want maps to use "complex"-style printouts:
	// meaning they will print their keys on separate lines than their values,
	// and keys may spread across mutiple lines if appropriate.
	//
	// If not set, a heuristic will be used based on if the map is known to
	// have keys that are complex enough that rendering them as oneline seems likely to overload.
	// See Config.useCmplxKeys for exactly how that's deteremined.
	UseMapComplexStyleAlways bool

	// For maps to use "complex"-style printouts (or not) per type.
	// See docs on UseMapComplexStyleAlways for the overview of what "complex"-style means.
	UseMapComplexStyleOnType map[schema.TypeName]bool
}

func (cfg *Config) init() {
	if cfg.Indentation == nil {
		cfg.Indentation = []byte{'\t'}
	}
}

// oneline decides if a value should be flatted into printing on a single,
// or if it's allowed to spread out over multiple lines.
// Note that this will not be asked if something outside of a value has already declared it's
// doing a oneline rendering; that railroads everything within it into that mode too.
func (cfg Config) oneline(typ schema.Type, isInKey bool) bool {
	return isInKey // Future: this could become customizable, with some kind of Always|OnlyInKeys|Never option enum per type.
}

/* TODO: not implemented or used
// useRepr decides if a value should be printed using its representation.
// Sometimes configuring this to be true for structs or unions with stringy representations
// will cause map printouts using them as keys to become drastically more readable
// (if with some loss of informativeness, or at least loss of explicitness).
func (cfg Config) useRepr(typ schema.Type, isInKey bool) bool {
	return false
}
*/

// useCmplxKeys decides if a map should print itself using a multi-line and extra-indented style for keys.
func (cfg Config) useCmplxKeys(mapn datamodel.Node) bool {
	if cfg.UseMapComplexStyleAlways {
		return true
	}
	tn, ok := mapn.(schema.TypedNode)
	if !ok {
		return false
	}
	tnt := tn.Type()
	if tnt == nil {
		return false
	}
	force, ok := cfg.UseMapComplexStyleOnType[tnt.Name()]
	if ok {
		return force
	}
	ti, ok := tnt.(*schema.TypeMap)
	if !ok { // Probably should never even have been asked, then?
		panic("how did you get here?")
	}
	return !cfg.oneline(ti.KeyType(), true)
}

// FUTURE: one could imagine putting an optional LinkSystem param into the Config, too, and some recursion control.
// It's definitely going to be the default to do zero recursion across links, though,
// as doing that requires creating graph visualizations, and that is both possible, yet to do well becomes rather nontrivial.
// Also, often a single node's tree visualization has been enough to get started debugging whatever I need to debug so far.

type printBuf struct {
	wr io.Writer

	Config
}

func (z *printBuf) writeString(s string) {
	z.wr.Write([]byte(s))
}

func (z *printBuf) doIndent(indentLevel int) {
	z.wr.Write(z.Config.StartingIndent)
	for i := 0; i < indentLevel; i++ {
		z.wr.Write(z.Config.Indentation)
	}
}

const (
	printState_normal       uint8 = iota
	printState_isKey              // may sometimes entersen or stringify things harder.
	printState_isValue            // signals that we're continuing a line that started with a key (so, don't emit indent).
	printState_isCmplxKey         // used to ask something to use multiline form, and an extra indent -- the opposite of what isKey does.
	printState_isCmplxValue       // we're continuing a line (so don't emit indent), and we're stuck in complex mode (so keep telling your children to stay in this state too).
)

func (z *printBuf) doString(indentLevel int, printState uint8, n datamodel.Node) {
	// First: indent.
	switch printState {
	case printState_normal, printState_isKey, printState_isCmplxKey:
		z.doIndent(indentLevel)
	}
	// Second: the typekind and type name; or, just the kind, if there's no type.
	//  Note: this can be somewhat overbearing -- for example, typed strings are going to get called out as `string<String>{"value"}`.
	//   This is rather agonizingly verbose, but also accurate; I'm not sure if we'd want to elide information about typed-vs-untyped entirely.
	if tn, ok := n.(schema.TypedNode); ok {
		var tnk schema.TypeKind
		var tntName string
		// Defensively check for nil node type
		if tnt := tn.Type(); tnt == nil {
			tntName = "?!nil"
			tnk = schema.TypeKind_Invalid
		} else {
			tntName = tnt.Name()
			tnk = tnt.TypeKind()
		}
		z.writeString(tnk.String())
		z.writeString("<")
		z.writeString(tntName)
		z.writeString(">")
		switch tnk {
		case schema.TypeKind_Invalid:
			z.writeString("{?!}")
		case schema.TypeKind_Map:
			// continue -- the data-model driven behavior is sufficient to handle the content.
		case schema.TypeKind_List:
			// continue -- the data-model driven behavior is sufficient to handle the content.
		case schema.TypeKind_Unit:
			return // that's it!  there's no content data for a unit type.
		case schema.TypeKind_Bool:
			// continue -- the data-model driven behavior is sufficient to handle the content.
		case schema.TypeKind_Int:
			// continue -- the data-model driven behavior is sufficient to handle the content.
		case schema.TypeKind_Float:
			// continue -- the data-model driven behavior is sufficient to handle the content.
		case schema.TypeKind_String:
			// continue -- the data-model driven behavior is sufficient to handle the content.
		case schema.TypeKind_Bytes:
			// continue -- the data-model driven behavior is sufficient to handle the content.
		case schema.TypeKind_Link:
			// continue -- the data-model driven behavior is sufficient to handle the content.
		case schema.TypeKind_Struct:
			// Very similar to a map, but keys aren't quoted.
			// Also, because it's possible for structs to be keys in a map themselves, they potentially need oneline emission.
			// Or, to customize emission in another direction if being a key in a map that's printing in "complex" mode.
			// FUTURE: there should also probably be some way to configure instructions to use their representation form instead.
			oneline :=
				printState == printState_isCmplxValue ||
					printState != printState_isCmplxKey && z.Config.oneline(tn.Type(), printState == printState_isKey)
			deepen := 1
			if printState == printState_isCmplxKey {
				deepen = 2
			}
			childState := printState_isValue
			if oneline {
				childState = printState_isCmplxValue
			}
			z.writeString("{")
			if !oneline && n.Length() > 0 {
				z.writeString("\n")
			}
			for itr := n.MapIterator(); !itr.Done(); {
				k, v, _ := itr.Next()
				if !oneline {
					z.doIndent(indentLevel + deepen)
				}
				fn, _ := k.AsString()
				z.writeString(fn)
				z.writeString(": ")
				z.doString(indentLevel+deepen, childState, v)
				if oneline {
					if !itr.Done() {
						z.writeString(", ")
					}
				} else {
					z.writeString("\n")
				}
			}
			if !oneline {
				z.doIndent(indentLevel)
			}
			z.writeString("}")
			return
		case schema.TypeKind_Union:
			// There will only be one thing in it, but we still have to use an iterator
			//  to figure out what that is if we're doing this generically.
			//  We can ignore the key and just look at the value type again though (even though those are the same in practice).
			_, v, _ := n.MapIterator().Next()
			z.writeString("{")
			z.doString(indentLevel, printState_isValue, v)
			z.writeString("}")
			return
		case schema.TypeKind_Enum:
			panic("TODO")
		default:
			panic("unreachable")
		}
	} else {
		if n.IsAbsent() {
			z.writeString("absent")
			return
		}
		z.writeString(n.Kind().String())
	}
	// Third: all the actual content.
	// FUTURE: this is probably gonna become... somewhat more conditional, and may end up being a sub-function to be reasonably wieldy.
	switch n.Kind() {
	case datamodel.Kind_Map:
		// Maps have to decide if they have complex keys and want to use an additionally-intended pattern to make that readable.
		// "Complex" here means roughly: if you try to cram them into one line, it doesn't look good.
		// This choice starts at the map but is mostly executed during the printing of the key:
		//  the key will start itself at normal indentation,
		//  but should then doubly indent all its nested values (assuming it has any).
		cmplxKeys := z.Config.useCmplxKeys(n)
		childKeyState := printState_isKey
		if cmplxKeys {
			childKeyState = printState_isCmplxKey
		}
		z.writeString("{")
		if n.Length() > 0 {
			z.writeString("\n")
		} else {
			z.writeString("}")
			return
		}
		for itr := n.MapIterator(); !itr.Done(); {
			k, v, err := itr.Next()
			if err != nil {
				z.doIndent(indentLevel + 1)
				z.writeString("!! map iteration step yielded error: ")
				z.writeString(err.Error())
				z.writeString("\n")
				break
			}
			z.doString(indentLevel+1, childKeyState, k)
			z.writeString(": ")
			z.doString(indentLevel+1, printState_isValue, v)
			z.writeString("\n")
		}
		z.doIndent(indentLevel)
		z.writeString("}")
	case datamodel.Kind_List:
		z.writeString("{")
		if n.Length() > 0 {
			z.writeString("\n")
		} else {
			z.writeString("}")
			return
		}
		for itr := n.ListIterator(); !itr.Done(); {
			idx, v, err := itr.Next()
			if err != nil {
				z.doIndent(indentLevel + 1)
				z.writeString("!! list iteration step yielded error: ")
				z.writeString(err.Error())
				z.writeString("\n")
				break
			}
			z.doIndent(indentLevel + 1)
			z.writeString(strconv.FormatInt(idx, 10))
			z.writeString(": ")
			z.doString(indentLevel+1, printState_isValue, v)
			z.writeString("\n")
		}
		z.doIndent(indentLevel)
		z.writeString("}")
	case datamodel.Kind_Null:
		// nothing: we already wrote the word "null" when we wrote the kind info prefix.
	case datamodel.Kind_Bool:
		z.writeString("{")
		if b, _ := n.AsBool(); b {
			z.writeString("true")
		} else {
			z.writeString("false")
		}
		z.writeString("}")
	case datamodel.Kind_Int:
		x, _ := n.AsInt()
		z.writeString("{")
		z.writeString(strconv.FormatInt(x, 10))
		z.writeString("}")
	case datamodel.Kind_Float:
		x, _ := n.AsFloat()
		z.writeString("{")
		z.writeString(strconv.FormatFloat(x, 'f', -1, 64))
		z.writeString("}")
	case datamodel.Kind_String:
		x, _ := n.AsString()
		z.writeString("{")
		z.writeString(strconv.QuoteToGraphic(x))
		z.writeString("}")
	case datamodel.Kind_Bytes:
		x, _ := n.AsBytes()
		z.writeString("{")
		dst := make([]byte, hex.EncodedLen(len(x)))
		hex.Encode(dst, x)
		z.writeString(string(dst))
		z.writeString("}")
	case datamodel.Kind_Link:
		x, _ := n.AsLink()
		z.writeString("{")
		z.writeString(x.String())
		z.writeString("}")
	}
}
