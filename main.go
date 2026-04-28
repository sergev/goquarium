package main

import (
	"fmt"
	"os"
)

func main() {
	if err := RunCLI(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
