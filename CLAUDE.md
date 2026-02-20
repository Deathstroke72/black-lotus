# CLAUDE.md

This file provides context and guidance for AI assistants working on the **black-lotus** repository.

## Project Overview

**black-lotus** is a Go multi-agent pipeline that uses the Claude API (Anthropic SDK) to generate complete microservice codebases from a single `ServiceDefinition` struct. Given a service description, it runs four specialized AI agents in sequence — each building on the previous agent's output — and saves the generated code artifacts to disk.

## Repository Structure

```
black-lotus/
├── .gitignore                        # Go-specific ignore rules
├── README.md                         # Project description and quick-start
├── CLAUDE.md                         # This file
├── main.go                           # Entry point; wires config → pipeline → artifacts
├── config/
│   └── service_definition.go         # ServiceDefinition struct + 3 built-in examples
├── agents/
│   ├── base.go                       # BaseAgent, Agent interface, Artifact, ParseArtifacts
│   ├── api_design_agent.go           # Designs REST API contracts
│   ├── backend_db_agent.go           # Service layer + PostgreSQL schema
│   ├── messaging_agent.go            # Kafka producers/consumers + outbox pattern
│   └── testing_security_agent.go     # Tests, JWT middleware, RBAC, rate limiting
└── orchestrator/
    └── pipeline.go                   # Pipeline: runs agents sequentially, saves artifacts
```

> No `go.mod` is present yet. The module name used in import paths is `inventory-agents`.
> Run `go mod init inventory-agents && go mod tidy` to initialize before building.

## Language & Toolchain

- **Language**: Go
- **Module name**: `inventory-agents` (inferred from import paths; no `go.mod` yet)
- **Minimum Go version**: not yet specified — update once `go.mod` is present
- **External dependency**: `github.com/anthropics/anthropic-sdk-go`
- **Build tool**: standard `go` toolchain

### Initializing the module (first-time setup)

```bash
go mod init inventory-agents
go mod tidy          # downloads anthropic-sdk-go and its transitive deps
```

## Core Concepts

### ServiceDefinition (`config/service_definition.go`)

The single input to the pipeline. Fields:

| Field               | Type       | Purpose                                              |
|---------------------|------------|------------------------------------------------------|
| `Name`              | `string`   | Short service name, e.g. `"inventory"`               |
| `Description`       | `string`   | Plain-English summary of what the service does       |
| `Language`          | `string`   | Target language for generated code, e.g. `"Go"`     |
| `Entities`          | `[]string` | Core domain objects, e.g. `["Product", "StockItem"]` |
| `Operations`        | `[]string` | Key business operations                              |
| `Integrations`      | `[]string` | External services (Kafka, REST, DBs)                 |
| `ExtraRequirements` | `[]string` | Freeform additional constraints                      |

`ServiceDefinition.Prompt()` serializes the struct into a structured text prompt injected into every agent.

Three built-in examples: `config.InventoryService()`, `config.PaymentsService()`, `config.NotificationsService()`.

### Agent Interface (`agents/base.go`)

```go
type Agent interface {
    Name() string
    Description() string
    Run(ctx context.Context, svc *config.ServiceDefinition, agentContext map[string]string) (*AgentResult, error)
}
```

`agentContext` is a `map[string]string` that accumulates each agent's output under its context key, allowing downstream agents to build on prior work.

Context keys by pipeline position:

| Index | Agent                     | Context key        |
|-------|---------------------------|--------------------|
| 0     | API Design Agent          | `api_design`       |
| 1     | Backend & Database Agent  | `backend_db`       |
| 2     | Messaging & Events Agent  | `messaging`        |
| 3     | Testing & Security Agent  | `testing_security` |

`BaseAgent` provides:
- `Chat(ctx, messages)` — calls the Claude API with the agent's system prompt
- `ParseArtifacts(output)` — extracts fenced code blocks from agent responses; filenames come from `// file: <name>` or `# file: <name>` hints inside blocks

### Pipeline (`orchestrator/pipeline.go`)

`NewPipeline(cfg)` registers agent factories (constructed lazily per service so system prompts are tailored). `Pipeline.Run(ctx, svc)` executes agents sequentially, passes each agent's trimmed output (≤3000 chars) into `agentContext` for the next agent, and returns a `PipelineResult`.

`SaveArtifacts(result, outputDir)` writes all artifacts under `<outputDir>/<service-name>/<agent-name>/`, plus a top-level `README.md` summary.

### Configuration & Environment

`config.Load()` reads configuration from environment variables. **`ANTHROPIC_API_KEY` is required** — the program calls `log.Fatal` if it is missing.

The output directory defaults to `./generated` and can be overridden via the first CLI argument:

```bash
go run main.go ./my-output-dir
```

## Development Workflows

### Prerequisites

```bash
export ANTHROPIC_API_KEY=your_key_here
go mod init inventory-agents   # only needed once, if go.mod is absent
go mod tidy
```

### Running

```bash
go run main.go
# or with a custom output directory:
go run main.go ./output
```

### Building

```bash
go build ./...
```

### Running tests

```bash
go test ./...
```

With race detector (recommended before merging):

```bash
go test -race ./...
```

### Linting

```bash
go vet ./...
```

If `golangci-lint` is adopted, update this section with the config file location and command.

### Formatting

All Go code must be formatted with `gofmt` (or `goimports`). CI should reject unformatted files.

```bash
gofmt -l -w .
```

## Output Structure

```
generated/
└── <service-name>/
    ├── README.md                              # Pipeline summary
    ├── api*design*agent/
    │   ├── output.md                          # Full agent markdown output
    │   └── api_1.go, api_2.go, ...
    ├── backend*and*database*agent/
    │   ├── output.md
    │   ├── service_1.go, ...
    │   └── migration_1.sql, ...
    ├── messaging*and*events*agent/
    │   ├── output.md
    │   └── messaging_1.go, ...
    └── testing*and*security*agent/
        ├── output.md
        └── test_1.go, ...
```

> Agent directory names are lowercased with spaces replaced by `*` (see `sanitizeName` in `pipeline.go`).

## Adding a New Agent

1. Create `agents/my_agent.go` implementing the `Agent` interface:

```go
package agents

type MyAgent struct{ *BaseAgent }

// NewMyAgent creates a MyAgent whose system prompt is tailored to svc.
func NewMyAgent(cfg *config.Config, svc *config.ServiceDefinition) *MyAgent {
    return &MyAgent{
        BaseAgent: NewBaseAgentForService(cfg, "My Agent", svc, responsibilities, outputFormat),
    }
}

func (a *MyAgent) Description() string { return "Short description shown in pipeline output" }

func (a *MyAgent) Run(ctx context.Context, svc *config.ServiceDefinition, agentCtx map[string]string) (*AgentResult, error) {
    prompt := svc.Prompt()
    if v, ok := agentCtx["api_design"]; ok {
        prompt += "\n\nAPI Design:\n" + v
    }
    messages := []anthropic.MessageParam{
        anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
    }
    output, err := a.Chat(ctx, messages)
    if err != nil {
        return nil, fmt.Errorf("[%s] failed: %w", a.Name(), err)
    }
    return &AgentResult{
        AgentName: a.Name(),
        Output:    output,
        Artifacts: ParseArtifacts(output),
    }, nil
}
```

2. Register a factory and context key in `orchestrator/pipeline.go`:

```go
// In NewPipeline agentFactories slice (append):
func(svc *config.ServiceDefinition) agents.Agent { return agents.NewMyAgent(cfg, svc) },

// In contextKeys slice (same index):
var contextKeys = []string{"api_design", "backend_db", "messaging", "testing_security", "my_agent"}
```

## Branching & Git Conventions

- **Default branch**: `main`
- **Feature branches**: `claude/<short-description>-<id>` (e.g. `claude/add-claude-documentation-4394X`)
- Commit messages should be imperative, concise, and describe *why* as well as *what*.
- Never push directly to `main`. Open a pull request and request review.

## Code Conventions

Follow standard Go idioms:

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
- **No `go.mod` yet**: do not create `go.mod` unless the user explicitly requests it; confirm the intended module path before doing so (current path inferred as `inventory-agents`).
- **Update this file**: whenever significant new packages, tools, or workflows are introduced, update `CLAUDE.md` to keep it accurate.
- **Commit discipline**: commit logical units of work with clear messages; do not bundle unrelated changes.
- **No secrets in code**: never commit API keys, passwords, or credentials. Use environment variables.

## Known Unknowns (to fill in as the project matures)

- [ ] `go.mod` / minimum Go version (module name inferred as `inventory-agents` from import paths)
- [ ] Full `config.Load()` implementation — which env vars does it read beyond `ANTHROPIC_API_KEY`? Default model and `MaxTokens` values?
- [ ] CI/CD pipeline (GitHub Actions, etc.)
- [ ] Linter configuration (`golangci-lint`, `.golangci.yml`)
- [ ] Deployment target (binary, Docker image, cloud function, etc.)
- [ ] Whether a `vendor/` directory or module proxy is used
