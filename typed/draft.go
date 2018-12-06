package typed

var (
	// Prelude types
	tString = TypeString{
		Name: "String",
	}
	// User's types
	/*
		struct Foo {
			f1 String
			f2 [nullable String]
			f3 optional String
		}
	*/
	tFoo = TypeObject{
		Name: "Foo",
		Fields: []ObjectField{
			{"f1", tString, false, true},
			{"f2", TypeList{
				Name:          "", // "[nullable String]" can be calculated.
				Anon:          true,
				ValueType:     tString,
				ValueNullable: true,
			}, false, false},
			{"f3", tString, true, false},
		},
	}
	// The Universe
	example = Universe{
		tFoo.Name: tFoo,
	}
)
