package main

import (
	"fmt"
	"os"

	"github.com/vak/goquarium/internal/aquarium"
)

func main() {
	if err := aquarium.RunCLI(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
