package fluent

import (
	"github.com/ipld/go-ipld-prime/datamodel"
)

func ToInterfaceValue(node datamodel.Node) (interface{}, error) {
	switch k := node.Kind(); k {
	case datamodel.Kind_Invalid:
		panic("invalid node")
	case datamodel.Kind_Null:
		return nil, nil
	case datamodel.Kind_Bool:
		return node.AsBool()
	case datamodel.Kind_Int:
		return node.AsInt()
	case datamodel.Kind_Float:
		return node.AsFloat()
	case datamodel.Kind_String:
		return node.AsString()
	case datamodel.Kind_Bytes:
		return node.AsBytes()
	case datamodel.Kind_Link:
		return node.AsLink()
	case datamodel.Kind_Map:
		outMap := make(map[string]interface{})
		for mi := node.MapIterator(); !mi.Done(); {
			k, v, err := mi.Next()
			if err != nil {
				return nil, err
			}
			kVal, err := k.AsString()
			if err != nil {
				return nil, err
			}
			vVal, err := ToInterfaceValue(v)
			if err != nil {
				return nil, err
			}
			outMap[kVal] = vVal

			if mi.Done() {
				break
			}
		}
		return outMap, nil
	case datamodel.Kind_List:
		outList := make([]interface{}, 0, node.Length())
		for li := node.ListIterator(); !li.Done(); {
			_, v, err := li.Next()
			if err != nil {
				return nil, err
			}
			vVal, err := ToInterfaceValue(v)
			if err != nil {
				return nil, err
			}
			outList = append(outList, vVal)

			if li.Done() {
				break
			}
		}
		return outList, nil
	}
	panic("unhandled case in switch")
}
