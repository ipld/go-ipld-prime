package methodsets

type Thing2ViaTypedef Thing

// This also compiles and works if you longhand the entire struct defn again...
//  as long as it's identical, it works.
//  (This does not extend to replacing field types with other-named but structurally identical types.)

func FlipTypedef(x *Thing) *Thing2ViaTypedef {
	return (*Thing2ViaTypedef)(x)
}

func (x *Thing2ViaTypedef) Pow() {
	x.Alpha = "typedef"
}
