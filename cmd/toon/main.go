package main

import (
	"fmt"
	"os"

	"github.com/toon-format/toon-cli/internal/cli"
)

var version = "dev"

func main() {
	stdinIsTTY, err := isTerminal(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "checking stdin: %v\n", err)
		os.Exit(1)
	}

	err = cli.Run(os.Args[1:], cli.Environment{
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		StdinIsTTY: stdinIsTTY,
		Version:    version,
	})
	if err != nil {
		os.Exit(1)
	}
}

func isTerminal(file *os.File) (bool, error) {
	info, err := file.Stat()
	if err != nil {
		return false, fmt.Errorf("stat stdin: %w", err)
	}

	return info.Mode()&os.ModeCharDevice != 0, nil
}
