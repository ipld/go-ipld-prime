package gengraphql

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func EmitFileHeader(w io.Writer) {
	// For now, no minima needed.
}

func EmitFileCompletion(w io.Writer, ts schema.TypeSystem) {
	writeTemplate(`
	type Root {
		{{- range $t := .GetTypes }}{{ if $t | IsComplex }}
		{{ $t.Name }}(id: ID): {{ $t | TypeSymbol }}
		{{ end }}{{- end}}
	}
	schema {
		query: Root
	}
	`, w, ts)
}
