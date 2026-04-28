package main

import (
	"fmt"
	"os"
)

// This is the program entry point.
// It sends command-line arguments to RunCLI.
// If something fails, it prints an error and exits.
func main() {
	if err := RunCLI(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
