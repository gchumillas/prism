package main

import (
	"fmt"
	"os"
)

func main() {
	if err := platformInit(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
