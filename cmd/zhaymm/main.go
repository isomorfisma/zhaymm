package main

import (
	"fmt"
	"os"
)

func main() {
	/* Memanggil fungsi Execute() yang ada di root.go
	Kalo err berisi error (bukan nil), akan dikasih liat error dan langsung keluar dari appnya. */
	
	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
