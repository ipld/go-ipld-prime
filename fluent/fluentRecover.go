package fluent

import "fmt"

func Recover(fn func()) error {
	return fmt.Errorf("TODO")
	// recover, but only handle fluent.Error
	// re-raise anything else
}
