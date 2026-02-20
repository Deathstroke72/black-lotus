# CLAUDE.md

This file provides context and guidance for AI assistants working on the **black-lotus** repository.

## Project Overview

**black-lotus** is a Go project currently in its initial setup phase. The repository has a single commit establishing the project scaffold (`.gitignore` and `README.md`). As development progresses, this document should be updated to reflect the evolving codebase.

## Repository Structure

```
black-lotus/
├── .gitignore     # Go-specific ignore rules (binaries, test artifacts, coverage output)
├── README.md      # Project description
└── CLAUDE.md      # This file
```

> As source files, packages, and tooling are added, update the directory tree above.

## Language & Toolchain

- **Language**: Go
- **Minimum Go version**: _not yet specified — update once `go.mod` is present_
- **Build tool**: standard `go` toolchain

### Expected Go module setup

When the Go module is initialized, the project root will contain:

```
go.mod    # module path and dependency declarations
go.sum    # cryptographic checksums for dependencies
```

## Development Workflows

### Building

```bash
go build ./...
```

### Running tests

```bash
go test ./...
```

Run with race detector (recommended before merging):

```bash
go test -race ./...
```

### Linting

```bash
go vet ./...
```

If a linter such as `golangci-lint` is adopted, update this section with the specific command and configuration file location.

### Formatting

All Go code must be formatted with `gofmt` (or `goimports`). CI should reject unformatted files.

```bash
gofmt -l -w .
```

## Branching & Git Conventions

- **Default branch**: `main`
- **Feature branches**: `claude/<short-description>-<id>` (e.g., `claude/add-claude-documentation-4394X`)
- Commit messages should be imperative, concise, and describe *why* as well as *what*.
- Never push directly to `main`. Open a pull request and request review.

## Code Conventions

Because this is a Go project, follow standard Go idioms:

1. **Package names**: lowercase, single word, no underscores.
2. **Error handling**: always check and propagate errors; do not swallow them silently.
3. **Exported identifiers**: must have a doc comment.
4. **Tests**: live alongside source in `_test.go` files; use table-driven tests where appropriate.
5. **No globals**: prefer dependency injection over package-level variables.
6. **Context propagation**: pass `context.Context` as the first argument to functions that perform I/O or long-running work.
7. **Avoid `init()`**: side-effects in `init()` make code hard to test and reason about.

## AI Assistant Guidelines

- **Read before editing**: always read a file in full before making changes.
- **Minimal diffs**: only change what is necessary; do not reformat unrelated code.
- **Do not invent structure**: if module layout, CLI framework, or database layer are not yet decided, ask rather than assume.
- **Update this file**: whenever significant new packages, tools, or workflows are introduced, update `CLAUDE.md` to keep it accurate.
- **Commit discipline**: commit logical units of work with clear messages; do not bundle unrelated changes.
- **No secrets in code**: never commit API keys, passwords, or credentials. Use environment variables or a secrets manager.

## Known Unknowns (to fill in as the project matures)

- [ ] Go module path (`module` directive in `go.mod`)
- [ ] Minimum Go version
- [ ] CLI framework or HTTP framework in use (e.g., `cobra`, `net/http`, `gin`, `chi`)
- [ ] Database / storage layer
- [ ] CI/CD pipeline (GitHub Actions, etc.)
- [ ] Linter configuration (`golangci-lint`, `.golangci.yml`)
- [ ] Deployment target (binary, Docker image, cloud function, etc.)
- [ ] Environment variable names and configuration approach
