package main

import (
	"fmt"
	"os"

	"github.com/go-nunu/nunu/cmd/nunu"
)

func main() {
	err := nunu.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "execute error:", err.Error())
		os.Exit(1)
	}
}
