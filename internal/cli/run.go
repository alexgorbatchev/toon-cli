package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	toon "github.com/toon-format/toon-go"

	"github.com/toon-format/toon-cli/internal/input"
)

const helpText = `Usage:
  toon < input.json > output.toon

Supported input formats: JSON, JSONC, NDJSON.

Flags:
  --help     Show this help message
  --version  Print the CLI version
`

type Environment struct {
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	StdinIsTTY bool
	Version    string
}

func Run(args []string, env Environment) error {
	version := env.Version
	if version == "" {
		version = "dev"
	}

	switch len(args) {
	case 0:
		return runConversion(env)
	case 1:
		switch args[0] {
		case "--help":
			_, err := io.WriteString(env.Stdout, helpText)
			return err
		case "--version":
			_, err := fmt.Fprintf(env.Stdout, "toon %s\n", version)
			return err
		default:
			err := fmt.Errorf("unknown flag: %s", args[0])
			_, _ = fmt.Fprintf(env.Stderr, "%v\n", err)
			return err
		}
	default:
		err := fmt.Errorf("unexpected arguments: %s", strings.Join(args, " "))
		_, _ = fmt.Fprintf(env.Stderr, "%v\n", err)
		return err
	}
}

func runConversion(env Environment) error {
	if env.StdinIsTTY {
		err := errors.New("no input detected on stdin")
		_, _ = fmt.Fprintf(env.Stderr, "%v\n", err)
		return err
	}

	data, err := io.ReadAll(env.Stdin)
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}

	value, err := input.Parse(data)
	if err != nil {
		return fmt.Errorf("parsing stdin: %w", err)
	}

	encoded, err := toon.MarshalString(value)
	if err != nil {
		return fmt.Errorf("encoding toon: %w", err)
	}

	if _, err := io.WriteString(env.Stdout, encoded); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	if !strings.HasSuffix(encoded, "\n") {
		if _, err := io.WriteString(env.Stdout, "\n"); err != nil {
			return fmt.Errorf("writing trailing newline: %w", err)
		}
	}

	return nil
}
