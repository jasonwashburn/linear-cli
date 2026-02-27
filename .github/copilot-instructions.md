# Copilot Instructions for linear-cli

## Project Summary

`linear-cli` is a Go command-line tool for managing [Linear](https://linear.app) issues and resources directly from the terminal via the [Linear GraphQL API](https://developers.linear.app/docs). The goal is a single static binary that supports full issue CRUD, plus future phases covering teams, comments, labels, projects, cycles, search, attachments, webhooks, and an interactive TUI mode.

## Tech Stack

| Concern | Choice |
|---|---|
| Language | **Go** (single static binary, fast startup) |
| CLI framework | **[Cobra](https://github.com/spf13/cobra)** |
| Config / env vars | **[Viper](https://github.com/spf13/viper)** |
| GraphQL client | **[hasura/go-graphql-client](https://github.com/hasura/go-graphql-client)** |
| Output formatting | **[tablewriter](https://github.com/olekurowiak/tablewriter)** (or similar) |

Authentication uses a Personal API Key read from the `LINEAR_API_KEY` environment variable. Never store the key in plain text; prefer OS keychains or permissions-locked files if persistence is needed.

## Repository Layout

```
.github/
  copilot-instructions.md   # This file
README.md
PLAN.md                     # Detailed roadmap and phase breakdown
.gitignore                  # Go-specific ignores
```

Once scaffolded, the project will follow a standard Go CLI layout:

```
cmd/
  root.go           # Cobra root command, global flags, Viper setup
  issue/
    list.go         # `linear issue list`
    get.go          # `linear issue get`
    create.go       # `linear issue create`
    update.go       # `linear issue update`
    delete.go       # `linear issue delete`
internal/
  api/
    client.go       # GraphQL client wrapper (hasura/go-graphql-client)
    queries.go      # Typed GraphQL query structs
    mutations.go    # Typed GraphQL mutation structs
  auth/
    auth.go         # Read LINEAR_API_KEY via Viper
  output/
    table.go        # Table formatting helpers
    json.go         # JSON output helpers
main.go             # Entry point — calls cmd.Execute()
go.mod
go.sum
```

## Build, Test, and Lint

> The project is in early development. Once `go.mod` is present, the commands below apply.

### Bootstrap (first time)

```bash
go mod tidy          # Install / tidy dependencies
```

### Build

```bash
go build -o linear ./...
```

### Run

```bash
LINEAR_API_KEY=<your-key> ./linear issue list
```

### Test

```bash
go test ./...
```

Run with verbose output and race detector during development:

```bash
go test -v -race ./...
```

### Lint

The project uses `golangci-lint`. Install it if missing:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Then run:

```bash
golangci-lint run ./...
```

Always run `go test ./...` and `golangci-lint run ./...` before submitting changes.

## Key Conventions

- **Environment variable**: `LINEAR_API_KEY` is the only required runtime configuration. Read it via Viper in the root command's `PersistentPreRunE`.
- **Output**: every command supports a `--json` flag for machine-readable output; default is a human-readable table.
- **GraphQL queries and mutations**: define typed Go structs (not raw strings) so the `hasura/go-graphql-client` can encode them safely.
- **Error messages**: surface Linear API errors as human-readable text, not raw JSON, unless `--json` is set.
- **Stdout vs stderr**: write all command output (tables, JSON results) to **stdout**; write all diagnostic messages, log lines, and errors to **stderr**. This allows callers to pipe or redirect command output cleanly.
- **Rate limits**: Linear allows 1,500 req/hr for Personal API Keys and 500 req/hr for OAuth tokens. Expose rate-limit headers in `--verbose` mode and back off automatically.
- **Command hierarchy**: `linear <resource> <verb> [flags]` — e.g. `linear issue list`, `linear issue create`.

## Linear GraphQL API Reference

- Base URL: `https://api.linear.app/graphql`
- Docs: https://developers.linear.app/docs/graphql/working-with-the-graphql-api
- Key operations used by the MVP:

| CLI command | GraphQL operation |
|---|---|
| `issue list` | `issues` query |
| `issue get` | `issue` query |
| `issue create` | `issueCreate` mutation |
| `issue update` | `issueUpdate` mutation |
| `issue delete` | `issueDelete` mutation |

## CI / Validation

No CI pipelines are configured yet. When added, they are expected to run:

1. `go build ./...`
2. `go test -race ./...`
3. `golangci-lint run ./...`

Trust the instructions above; only search the codebase if information here is incomplete or found to be incorrect.
