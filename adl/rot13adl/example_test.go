package rot13adl_test

func ExampleWoo() {

	// Output:
}

// An Unmarshal2 function could take a (fairly complex) argument that says what NodePrototype to use in certain positions
//  (is this accomplished with a Selector?  that's a scary large dependency, but maybe.).
//  This can be used for many purposes (efficiency, misc tuning, expecting a *schema* thing part way through...),
//  and it could also be used to say where a NodePrototype for an ADL's substrate should be used.
//
// Schemas, even the repr builders, still ultimately result in "returning" (not really, but, roll with me here) the high-level typed info.
// ...
// This ADL, as currently drafted, does not.
// That's... maybe bad?
// It can be passed through the Reify method and there's a feeling of consistency, there.
// But I'm not sure.
