# `toon`

`toon` reads JSON, JSONC, or NDJSON from stdin and writes TOON to stdout using [`toon-go`](https://github.com/toon-format/toon-go).

## Requirements

- Go 1.26 or newer to build from source

## Quick start

Run `toon` directly from a checkout:

```sh
printf '{"name":"Ada","id":1}\n' | go run ./cmd/toon
```

Output:

```text
name: Ada
id: 1
```

Or use the checked-in `justfile`:

```sh
printf '{"name":"Ada","id":1}\n' | just run
```

## Usage

`toon` is stdin-first. It accepts exactly two flags: `--help` and `--version`.

```sh
toon < input.json > output.toon
```

From a local checkout:

```sh
go run ./cmd/toon < input.json > output.toon
```

### JSONC

`toon` accepts JSONC as a single top-level document.

```sh
printf '{\n  // preserve field order\n  "name": "Ada",\n  "languages": ["go", "toon"],\n}\n' | go run ./cmd/toon
```

### NDJSON

NDJSON input is decoded as one top-level array before TOON encoding.

```sh
printf '{"id":1,"name":"Ada"}\n{"id":2,"name":"Bob"}\n' | go run ./cmd/toon
```

Output:

```text
[2]{id,name}:
  1,Ada
  2,Bob
```

## Version output

Local builds report `toon dev`:

```sh
go run ./cmd/toon --version
```

Tagged release builds inject the runtime version through GoReleaser, so `--version` matches the Git tag.

## Numeric behavior

`toon` preserves integer precision for JSON-style inputs. Integer literals beyond the IEEE-754 safe range are emitted as quoted decimal strings so their digits are not rounded during TOON encoding.

## Development

Common local commands are available through `just`:

- Run CLI: `just run` (prints help when no stdin is piped)
- Help: `just help`
- Version: `just version`
- Test: `just test`
- Focused tests: `just test-focused`
- Vet: `just vet`
- Lint: `just lint`
- Build: `just build`
- Full verification: `just verify`

Equivalent Go commands:

- Test: `go test ./...`
- Vet: `go vet ./...`
- Lint: `golangci-lint run`
- Build: `go build ./...`

## Releases

GitHub Actions runs CI on pushes and pull requests to `main`.

Pushing a tag that matches `v*` runs the release workflow, which:

- verifies the module with `go mod tidy -diff`, `go build ./...`, `go vet ./...`, `go test -race ./...`, and `golangci-lint run`
- builds release archives with GoReleaser
- injects the runtime version into `main.version` for `toon --version`
