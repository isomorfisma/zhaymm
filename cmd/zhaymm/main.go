package main

import (
	"fmt"
	"os"
)

func main() {
	/* Calls Execute() func from root.go
	Exits on error */
	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
