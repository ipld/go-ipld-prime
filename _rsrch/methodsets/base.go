package methodsets

type Thing struct {
	Alpha string
	Beta  string
}

func (x *Thing) Pow() {
	x.Alpha = "base"
}

type ThingPrivate struct {
	a string
	b string
}
