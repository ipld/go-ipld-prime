package gendemo

// Type is a struct embeding a NodePrototype/Type for every Node implementation in this package.
// One of its major uses is to start the construction of a value.
// You can use it like this:
//
// 		gendemo.Type.YourTypeName.NewBuilder().BeginMap() //...
//
// and:
//
// 		gendemo.Type.OtherTypeName.NewBuilder().AssignString("x") // ...
//
var Type typeSlab

type typeSlab struct {
	Int                     _Int__Prototype
	Int__Repr               _Int__ReprPrototype
	Map__String__Msg3       _Map__String__Msg3__Prototype
	Map__String__Msg3__Repr _Map__String__Msg3__ReprPrototype
	Msg3                    _Msg3__Prototype
	Msg3__Repr              _Msg3__ReprPrototype
	String                  _String__Prototype
	String__Repr            _String__ReprPrototype
}
