package typegen

import (
	"io"
	"text/template"

	declaration "github.com/ipld/go-ipld-prime/typed/declaration"
	wish "github.com/warpfork/go-wish"
)

type generationMonad struct {
	// - all the basic structs and typedefs
	// - includes all the consts for enums.
	// - includes reflective Type() methods.
	typesFile io.Writer

	// - all the ipld.Node accessors.
	// - all the ipld.NodeBuilder types and their methods.
	hypergenericInterfacesFile io.Writer

	// - all the schlep for enums and unions to be closed interfaces.
	closedMembershipBoilerplateFile io.Writer
}

func (gm generationMonad) writeType(name declaration.TypeName, dt declaration.Type) {
	var tmpl string
	switch dt.(type) {
	case declaration.TypeBool:
		tmpl = wish.Dedent(`
			type {{ .Name }} bool
		`) + "\n"
	case declaration.TypeString:
		tmpl = wish.Dedent(`
			type {{ .Name }} string
		`) + "\n"
	case declaration.TypeBytes:
		// punt, not required for typedecl bootstrapping
	case declaration.TypeInt:
		// punt, not required for typedecl bootstrapping
	case declaration.TypeFloat:
		// punt, not required for typedecl bootstrapping
	case declaration.TypeMap:
		tmpl = wish.Dedent(`
			type {{ .Name }} struct {
				val map[{{ .Type.KeyType }}]{{ .Type.ValueType }}
				ord []{{ .Type.KeyType }}
			}
		`) + "\n"
	case declaration.TypeList:
		// punt, not required for typedecl bootstrapping
	case declaration.TypeLink:
		// punt, not required for typedecl bootstrapping
	case declaration.TypeUnion:
		// TODO
	case declaration.TypeStruct:
		tmpl = wish.Dedent(`
			type {{ .Name }} struct {
				{{- range .Fields -}}
				// wow, we need a very limited set of things for ranging to work.
				// a channel (but it'd have to be fully pre-buffered, or need a goroutine, so no)
				// or a slice, basically.  that's it.  there are no options for generatives.
				// one long shot is to use block or define plus with to do recursion, but, lol.
				// all of these require converter methods of some kind if Node only has generators.
				// I guess we will indeed have to keep an immediate-mode keys list for ranging purposes!
				{{- end -}}
			}
		`) + "\n"
	case declaration.TypeEnum:
		// TODO
	}
	template.Must(template.New("").Parse(tmpl)).Execute(gm.typesFile, map[string]interface{}{
		"Name": name,
		"Type": dt,
	})
}

/*
	FUTURE
	------

	### Add reflection'ish Type() getters:

	```
	func (x {{ .Name }}) Type() typesystem.Type { return typeOf{{.Name}}; }

	var typeOf{{ .Name }} = typesystem.Type{

	}
	```

	(This'll come pretty late in the game; it'll require full construction
	to be working, not least of which because there will be cycles to break
	using the two-phase approach again!)
*/
