package fluent

// AllKeyStrings is a shorthand to iterate a map node and collect all the keys
// (and convert them to strings), returning them in a slice.
func AllKeyStrings(n Node) []string {
	itr := n.MapIterator()
	res := make([]string, n.Length())
	for i := 0; !itr.Done(); i++ {
		k, _ := itr.Next()
		res[i] = k.AsString()
	}
	return res
}
