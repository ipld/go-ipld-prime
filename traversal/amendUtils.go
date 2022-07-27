package traversal

import (
	"log"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func nodeForType(base interface{}) datamodel.Node {
	if base == nil {
		return nil
	}
	switch typ := base.(type) {
	case Amender:
		return base.(Amender).Build()
	case datamodel.Node:
		return base.(datamodel.Node)
	case datamodel.Link:
		return basicnode.NewLink(base.(datamodel.Link))
	case bool:
		return basicnode.NewBool(base.(bool))
	case int8:
		return basicnode.NewInt(int64(base.(int8)))
	case int16:
		return basicnode.NewInt(int64(base.(int16)))
	case int32:
		return basicnode.NewInt(int64(base.(int32)))
	case int64:
		return basicnode.NewInt(base.(int64))
	case int:
		return basicnode.NewInt(int64(base.(int)))
	case uint8:
		return basicnode.NewUint(uint64(base.(uint8)))
	case uint16:
		return basicnode.NewUint(uint64(base.(uint16)))
	case uint32:
		return basicnode.NewUint(uint64(base.(uint32)))
	case uint64:
		return basicnode.NewUint(base.(uint64))
	case uint:
		return basicnode.NewUint(uint64(base.(uint)))
	case float32:
		return basicnode.NewFloat(float64(base.(float32)))
	case float64:
		return basicnode.NewFloat(base.(float64))
	case string:
		return basicnode.NewString(base.(string))
	case []byte: // Special handling for array of bytes
		{
			return basicnode.NewBytes(base.([]byte))
		}
	case map[string]interface{}:
		{
			a := AmendOptions{}.newMapAmender(nil, nil, false)
			for k, v := range base.(map[string]interface{}) {
				a.(datamodel.Map).Put(k, v)
			}
			return a.Build()
		}
	case []interface{}:
		{
			a := AmendOptions{}.newListAmender(nil, nil, false)
			a.(datamodel.List).Append(base.([]interface{}))
			return a.Build()
		}
	default:
		log.Printf("invalid type: %s", typ)
		panic("unreachable")
	}
}
