package realgen

import (
	"github.com/ipld/go-ipld-prime/schema"
)

// Code generated go-ipld-prime DO NOT EDIT.

const (
	midvalue = schema.Maybe(4)
)

type maState uint8

const (
	maState_initial maState = iota
	maState_midKey
	maState_expectValue
	maState_midValue
	maState_finished
)
