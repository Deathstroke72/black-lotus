# CLAUDE.md

This file provides context and guidance for AI assistants working on the **black-lotus** repository.

## Project Overview

**black-lotus** is a Go multi-agent pipeline that uses the Claude API (Anthropic SDK) to generate complete microservice codebases from a single `ServiceDefinition` struct. Given a service description, it runs four specialized AI agents in sequence — each building on the previous agent's output — and saves the generated code artifacts to disk.

## Repository Structure

```
black-lotus/
├── .gitignore                              # Go-specific ignore rules
├── README.md                               # Project description and quick-start
├── CLAUDE.md                               # This file
├── main.go                                 # Entry point; wires config → pipeline → artifacts
└── lotus-agents/
    ├── agents/
    │   ├── base.go                         # BaseAgent, Agent interface, Artifact, ParseArtifacts, ToJSON
    │   ├── api_design_agent.go             # Designs REST API contracts
    │   ├── backend_db_agent.go             # Service layer + PostgreSQL schema
    │   ├── messaging_agent.go              # Kafka producers/consumers + outbox pattern
    │   └── testing_security_agent.go       # Tests, JWT middleware, RBAC, rate limiting
    ├── config/
    │   └── service_definition.go           # ServiceDefinition struct + 3 built-in examples
    └── orchestrator/
        └── pipeline.go                     # Pipeline: runs agents sequentially, saves artifacts
```

> No `go.mod` is present yet. The module name used in all import paths is `lotus-agents`.
> Run `go mod init lotus-agents && go mod tidy` to initialize before building.

## Language & Toolchain

- **Language**: Go
- **Module name**: `lotus-agents` (confirmed from all import paths; no `go.mod` yet)
- **Minimum Go version**: not yet specified — update once `go.mod` is present
- **External dependency**: `github.com/anthropics/anthropic-sdk-go`
- **Build tool**: standard `go` toolchain

### Initializing the module (first-time setup)

```bash
go mod init lotus-agents
go mod tidy          # downloads anthropic-sdk-go and its transitive deps
```

## Core Concepts

### ServiceDefinition (`lotus-agents/config/service_definition.go`)

The single input to the pipeline. Fields:

| Field               | Type       | Purpose                                               |
|---------------------|------------|-------------------------------------------------------|
| `Name`              | `string`   | Short service name, e.g. `"inventory"`                |
| `Description`       | `string`   | Plain-English summary of what the service does        |
| `Language`          | `string`   | Target language for generated code, e.g. `"Go"`       |
| `Entities`          | `[]string` | Core domain objects, e.g. `["Product", "StockItem"]`  |
| `Operations`        | `[]string` | Key business operations                               |
| `Integrations`      | `[]string` | External services (Kafka, REST, DBs)                  |
| `ExtraRequirements` | `[]string` | Freeform additional constraints                       |

`ServiceDefinition.Prompt()` serializes the struct into a structured text prompt injected into every agent.

Three built-in examples: `config.InventoryService()`, `config.PaymentsService()`, `config.NotificationsService()`.

### Config (`lotus-agents/config/`)

`config.Load()` reads configuration from environment variables and returns a `*config.Config`. The `Config` struct has at minimum these fields (used by `BaseAgent`):

| Field             | Source              | Notes                                               |
|-------------------|---------------------|-----------------------------------------------------|
| `AnthropicAPIKey` | `ANTHROPIC_API_KEY` | Required — `main.go` calls `log.Fatal` if absent    |
| `Model`           | env / default       | Passed to `anthropic.Model(cfg.Model)`              |
| `MaxTokens`       | env / default       | Passed as `int64(cfg.MaxTokens)` to the Claude API  |

> `config.Load()` is called in `main.go` but its implementation is not in any visible source file — see Known Unknowns.

### Agent Interface (`lotus-agents/agents/base.go`)

```go
type Agent interface {
    Name() string
    Description() string
    Run(ctx context.Context, svc *config.ServiceDefinition, agentContext map[string]string) (*AgentResult, error)
}
```

`agentContext` is a `map[string]string` that accumulates each agent's output under its context key. The pipeline seeds it with:

```go
agentContext := map[string]string{
    "project_context": svc.Prompt(),
}
```

Context keys by pipeline position:

| Index | Agent                    | Writes to key      | Reads from context          |
|-------|--------------------------|--------------------|-----------------------------|
| 0     | API Design Agent         | `api_design`       | `project_context`           |
| 1     | Backend & Database Agent | `backend_db`       | `api_design`                |
| 2     | Messaging & Events Agent | `messaging`        | `backend_db`                |
| 3     | Testing & Security Agent | `testing_security` | `api_design`, `backend_db`  |

Each agent's output is trimmed to ≤3000 characters before being stored in `agentContext` for downstream agents.

### BaseAgent (`lotus-agents/agents/base.go`)

`BaseAgent` provides shared Claude API functionality:

- `NewBaseAgent(cfg, name, systemPrompt)` — base constructor; creates an `anthropic.Client`
- `NewBaseAgentForService(cfg, name, svc, responsibilities, outputFormat)` — builds a tailored system prompt via `buildSystemPrompt` and calls `NewBaseAgent`
- `(b *BaseAgent) Name() string` — returns the agent name
- `(b *BaseAgent) Chat(ctx, messages)` — calls `Messages.New` with the agent's system prompt; concatenates all `"text"` content blocks from the response
- `(b *BaseAgent) WithSystemPrompt(prompt) *BaseAgent` — returns a shallow copy with a different system prompt
- `ParseArtifacts(output)` — package-level function; extracts fenced code blocks from markdown; detects filenames from `// file: <name>` or `# file: <name>` hints inside blocks
- `ToJSON(v any) string` — package-level helper; marshals any value to indented JSON for context passing

#### Key types

```go
type Artifact struct {
    Filename string
    Content  string
    Language string
}

type AgentResult struct {
    AgentName string
    Output    string
    Artifacts []Artifact
    Error     error
}
```

### Agent Responsibilities & Default Filenames

| Agent                    | Produces                                                          | Default filename (no hint)              |
|--------------------------|-------------------------------------------------------------------|-----------------------------------------|
| API Design Agent         | Router setup, request/response structs, OpenAPI godoc             | `api_<n>.go`                            |
| Backend & Database Agent | PostgreSQL schema, pgx repositories, service layer, migrations    | `service_<n>.go` / `migration_<n>.sql`  |
| Messaging & Events Agent | Kafka producers (outbox pattern), consumers, dead-letter handling | `messaging_<n>.go`                      |
| Testing & Security Agent | Unit/integration tests, JWT (RS256), RBAC, rate limiter, Makefile | `test_<n>.go`                           |

### Pipeline (`lotus-agents/orchestrator/pipeline.go`)

`NewPipeline(cfg)` registers agent factories (constructed lazily per service so system prompts are tailored). `Pipeline.Run(ctx, svc)` executes agents sequentially and returns a `PipelineResult`.

```go
type PipelineResult struct {
    Service   *config.ServiceDefinition
    Results   []*agents.AgentResult
    StartTime time.Time
    EndTime   time.Time
    Duration  time.Duration
}
```

`SaveArtifacts(result, outputDir)` writes all artifacts under `<outputDir>/<service-name>/<agent-name>/`, plus a top-level `README.md` summary.

#### `sanitizeName` behaviour

Agent directory names are produced by `sanitizeName`: spaces → `*`, `&` → `and`, `/` → `*`, then lowercased.

Examples:
- `"API Design Agent"` → `api*design*agent`
- `"Backend & Database Agent"` → `backend*and*database*agent`
- `"Messaging & Events Agent"` → `messaging*and*events*agent`
- `"Testing & Security Agent"` → `testing*and*security*agent`

#### `extensionForLang` mapping

| Language        | Extension |
|-----------------|-----------|
| `"go"`          | `.go`     |
| `"sql"`         | `.sql`    |
| `"yaml"`/`"yml"`| `.yaml`   |
| `"json"`        | `.json`   |
| (anything else) | `.txt`    |

## Development Workflows

### Prerequisites

```bash
export ANTHROPIC_API_KEY=your_key_here
go mod init lotus-agents   # only needed once, if go.mod is absent
go mod tidy
```

### Running

```bash
go run main.go
# or with a custom output directory:
go run main.go ./output
```

The service to generate is configured in `main.go`. Swap `config.InventoryService()` for `config.PaymentsService()`, `config.NotificationsService()`, or define a custom `ServiceDefinition` inline.

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
    ├── README.md                                  # Pipeline summary (description, duration, file list)
    ├── api*design*agent/
    │   ├── output.md                              # Full agent markdown output
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

1. Create `lotus-agents/agents/my_agent.go` implementing the `Agent` interface:

```go
package agents

import (
    "context"
    "fmt"

    "github.com/anthropics/anthropic-sdk-go"
    "lotus-agents/config"
)

const myAgentResponsibilities = `- ...`
const myAgentOutputFormat = "..."

type MyAgent struct{ *BaseAgent }

func NewMyAgent(cfg *config.Config, svc *config.ServiceDefinition) *MyAgent {
    return &MyAgent{
        BaseAgent: NewBaseAgentForService(cfg, "My Agent", svc, myAgentResponsibilities, myAgentOutputFormat),
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

2. Register a factory and context key in `lotus-agents/orchestrator/pipeline.go`:

```go
// In NewPipeline agentFactories slice (append at end):
func(svc *config.ServiceDefinition) agents.Agent { return agents.NewMyAgent(cfg, svc) },

// In contextKeys slice (same index as factory above):
var contextKeys = []string{"api_design", "backend_db", "messaging", "testing_security", "my_agent"}
```

## Branching & Git Conventions

- **Default branch**: `main`
- **Feature branches**: `claude/<short-description>-<id>` (e.g. `claude/add-claude-documentation-X14MH`)
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
- **Module name**: the confirmed module name is `lotus-agents` (from all import paths). Do not create `go.mod` unless the user explicitly requests it.
- **Update this file**: whenever significant new packages, tools, or workflows are introduced, update `CLAUDE.md` to keep it accurate.
- **Commit discipline**: commit logical units of work with clear messages; do not bundle unrelated changes.
- **No secrets in code**: never commit API keys, passwords, or credentials. Use environment variables.

## Known Unknowns (to fill in as the project matures)

- [ ] `go.mod` / minimum Go version (module name confirmed as `lotus-agents` from all import paths)
- [ ] `config.Load()` implementation — the function is called in `main.go` but its source is not present in any visible file. Known fields on the returned `*config.Config`: `AnthropicAPIKey` (string), `Model` (string), `MaxTokens` (int). Which env vars control `Model` and `MaxTokens`? What are the defaults?
- [ ] CI/CD pipeline (GitHub Actions workflows, etc.)
- [ ] Linter configuration (`golangci-lint`, `.golangci.yml`)
- [ ] Deployment target (binary, Docker image, cloud function, etc.)
- [ ] Whether a `vendor/` directory or module proxy is used
