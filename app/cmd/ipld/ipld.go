package main

import (
	"os"

	"github.com/ipld/go-ipld-prime/app"
)

func main() {
	code, _ := app.Main(os.Args, os.Stdin, os.Stdout, os.Stderr)
	os.Exit(code)
}
