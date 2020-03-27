package corpus

import (
	"bytes"
)

func ents(n int, segFn func(i int) string) string {
	if n <= 0 {
		return ""
	}
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		buf.WriteString(segFn(i))
		buf.WriteString(`,`)
	}
	buf.Truncate(buf.Len() - 1)
	return buf.String()
}
