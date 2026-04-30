package main

import (
	"fmt"
	"os"
)

var DebugLogEnabled bool

// setupStderrFromEnv redirects stderr to a file when GOQUARIUM_STDERR is set.
// Empty env means keep default stderr.
func setupStderrFromEnv() error {
	logPath := os.Getenv("GOQUARIUM_STDERR")
	if logPath == "" {
		DebugLogEnabled = false
		return nil
	}
	DebugLogEnabled = true
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("open GOQUARIUM_STDERR file %q: %w", logPath, err)
	}
	os.Stderr = f
	return nil
}

// This is the program entry point.
// It sends command-line arguments to RunCLI.
// If something fails, it prints an error and exits.
func main() {
	if err := setupStderrFromEnv(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := RunCLI(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
