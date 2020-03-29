package methodsets

type AliasThing = Thing

// func (x *AliasThing) Pow() {}
// NOPE!
// ./viaAliases.go:5:22: (*Thing).Pow redeclared in this block
//	previous declaration at ./base.go:8:6

type AliasPtr = *Thing

// ^ Oddly, works.

// func (x *AliasPtr) Pow() {}
// NOPE!
// ./aliases.go:14:6: invalid receiver type **Thing (*Thing is not a defined type)

// func (x AliasPtr) Pow() {}
// NOPE!
// ./aliases.go:18:19: (*Thing).Pow redeclared in this block
// 	previous declaration at ./base.go:8:6

/*
	Conclusion: no joy.
	Aliases really are a syntactic sugar thing, and do not seem to enable
	any interesting tricks that would not otherwise be possible,
	and certainly don't appear to get us closer to the "methodsets" semantic I yearn for.
*/
