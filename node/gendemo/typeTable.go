package gendemo

// Type is a struct embeding a NodeStyle/Type for every Node implementation in this package.
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
	Int                     _Int__Style
	Int__Repr               _Int__ReprStyle
	Map__String__Msg3       _Map__String__Msg3__Style
	Map__String__Msg3__Repr _Map__String__Msg3__ReprStyle
	Msg3                    _Msg3__Style
	Msg3__Repr              _Msg3__ReprStyle
	String                  _String__Style
	String__Repr            _String__ReprStyle
}
