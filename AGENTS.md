# toon-cli

Go CLI that reads JSON, JSONC, or NDJSON from stdin and prints TOON.

## Commands
- Just list: `just`
- Run: `printf '{"name":"Ada","id":1}\n' | just run`
- Help: `just help`
- Version: `just version`
- Verify: `just verify`
- Test: `go test ./...`
- Test focused packages: `go test ./internal/input ./internal/cli`
- Vet: `go vet ./...`
- Lint: `golangci-lint run`
- Build: `go build ./...`
- Help: `go run ./cmd/toon --help`
- Convert sample input: `printf '{"name":"Ada","id":1}\n' | go run ./cmd/toon`

## Setup
- Use Go 1.26+.
- Local runs print `toon dev`; tagged release binaries get their version from GoReleaser via `-X main.version={{.Version}}`.

## Conventions
- Keep OS wiring in `cmd/toon/main.go`; keep flag handling and stdin/stdout orchestration in `internal/cli/`.
- Keep input-format decoding in `internal/input/`.
- Preserve input object order by building `toon.NewObject(...)`; do not decode objects into plain `map[string]any` when order matters.
- Parse integer literals as `big.Int` values instead of routing them through `json.Number`; `toon-go` rounds large `json.Number` integers.

## Gotchas
- `just run` prints CLI help when stdin is a TTY; the compiled CLI still returns `no input detected on stdin` when run without piped input.
- NDJSON is converted into one top-level array before TOON encoding.
- JSONC support is single-document only. Commented multi-document streams are out of scope unless the product requirements change.

## Boundaries
- Always: run `go test ./...`, `go vet ./...`, and `golangci-lint run` after changing Go code or GitHub workflows.
- Ask first: changes to CLI flags, the stdin-only contract, dependency additions, module path changes, or release automation semantics.
- Never: commit `.tmp/`, commit `dist/`, or replace ordered TOON objects with plain maps in the parser.

## References
- `README.md`
- `internal/input/parse.go`
- `internal/cli/run.go`
- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `.goreleaser.yml`
