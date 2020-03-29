package nozeros

type internal struct {
	x string
}

func (x *internal) Pow() {}

// func (x ExportedPtr) Pow() {} // no, "invalid receiver type ExportedPtr (ExportedPtr is a pointer type)".
// wait, no, this was only illegal for `type ExportedPtr *internal` -- AHHAH

type ExportedAlias = internal

type ExportedPtr = *internal

// type ExportedPtr *internal // ALMOST works...
//  except somewhat bizarrely, 'Pow()' becomes undefined and invisible on the exported type.
//  As far as I can tell, aside from that, it's identical to using the alias.

func NewExportedPtr(v string) ExportedPtr {
	return &internal{v}
}

// ---

type FooData struct {
	x string
}

func (x Foo) Pow() {
	x.x = "waht"
}
func (x Foo) Read() string {
	return x.x
}

type foo FooData

type Foo = *foo

func NewFoo(v string) Foo {
	return &foo{v}
}
