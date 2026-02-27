# linear-cli — Plan & Roadmap

A command-line interface for interacting with [Linear](https://linear.app) via
the [Linear GraphQL API](https://developers.linear.app/docs).

---

## Authentication

Linear's API supports two authentication methods:

| Method | Use-case |
|---|---|
| **Personal API Key** | Local / personal tooling — easiest to get started |
| **OAuth 2.0** | Multi-user / shared tooling |

For the MVP we will use a **Personal API Key** passed via the environment
variable `LINEAR_API_KEY`. The key should not be stored in a plain-text
dotfile; if persistent storage is required in the future, prefer an OS
keychain or credential store (or at minimum a permissions-locked config file).

---

## MVP — Basic Issue CRUD

The MVP goal is a working CLI that lets a developer manage Linear issues from
the terminal without opening a browser.

### Commands

```
linear issue list [--team <team>] [--state <state>] [--assignee <user>]
linear issue get  <issue-id>
linear issue create --title <title> [--team <team>] [--description <desc>]
                    [--priority <0-4>] [--assignee <user>]
linear issue update <issue-id> [--title <title>] [--description <desc>]
                    [--state <state>] [--priority <0-4>] [--assignee <user>]
linear issue delete <issue-id>
```

### Underlying API operations

| Command | GraphQL operation |
|---|---|
| `issue list` | `issues` |
| `issue get` | `issue` |
| `issue create` | `issueCreate` |
| `issue update` | `issueUpdate` |
| `issue delete` | `issueDelete` |

### Suggested tech stack

| Concern | Choice | Rationale |
|---|---|---|
| Language | **Go** | Single static binary, fast startup, minimal runtime dependencies |
| CLI framework | **[Cobra](https://github.com/spf13/cobra)** | De-facto standard Go CLI framework, mature & widely used |
| Config / flags | **[Viper](https://github.com/spf13/viper)** | Pairs with Cobra; handles env vars, config files, and flag binding |
| GraphQL client | **[hasura/go-graphql-client](https://github.com/hasura/go-graphql-client)** | Actively maintained, idiomatic Go, typed struct-based queries |
| Output formatting | **[tablewriter](https://github.com/olekurowiak/tablewriter)** | Simple table rendering, no heavy deps |

### MVP deliverables checklist

- [ ] Project scaffold (Go module, `cmd/` layout, Cobra root command)
- [ ] Auth helper — read `LINEAR_API_KEY` from env via Viper
- [ ] `issue list` — paginated, filterable list of issues
- [ ] `issue get` — show full detail for one issue
- [ ] `issue create` — interactive + flag-driven creation
- [ ] `issue update` — update any writable field
- [ ] `issue delete` — with confirmation prompt
- [ ] Error handling & human-readable API error messages
- [ ] `--json` flag on every command for scripting / piping
- [ ] README with installation & usage instructions

---

## Roadmap — Future Features

Ordered roughly by user value and implementation effort.

### Phase 2 — Teams & Users

```
linear team list
linear user list
linear user me
```

- List all teams the authenticated user belongs to
- List workspace members (for `--assignee` autocomplete)
- Show information about the authenticated user (`viewer` query)

### Phase 3 — Comments

```
linear issue comment list <issue-id>
linear issue comment add  <issue-id> --body <text>
linear issue comment edit <comment-id> --body <text>
linear issue comment delete <comment-id>
```

Full CRUD on issue comments via `commentCreate`, `commentUpdate`,
`commentDelete` mutations.

### Phase 4 — Labels

```
linear label list  [--team <team>]
linear label create --name <name> --color <hex> [--team <team>]
linear issue label add    <issue-id> <label-id>
linear issue label remove <issue-id> <label-id>
```

Manage labels and attach/detach them from issues using the `labelCreate` /
`labelUpdate` / `labelDelete` mutations and `issueUpdate`.

### Phase 5 — Projects

```
linear project list [--team <team>]
linear project get  <project-id>
linear project create --name <name> --team <team> [--description <desc>]
linear project update <project-id> [--name <name>] [--state <state>]
```

Full CRUD on projects via `projectCreate`, `projectUpdate`, `projectDelete`.

### Phase 6 — Cycles (Sprints)

```
linear cycle list  [--team <team>]
linear cycle get   <cycle-id>
linear cycle create --team <team> --name <name> --start <date> --end <date>
linear cycle issue add    <cycle-id> <issue-id>
linear cycle issue remove <cycle-id> <issue-id>
```

Manage sprint-style cycles and move issues in/out via `cycleCreate` and
`cycleUpdate` mutations.

### Phase 7 — Search & Filtering

```
linear search <query>
```

Full-text search across issues (and optionally projects, comments) using
Linear's `searchIssues` query. Support advanced filter expressions passed as
JSON or a mini-DSL.

### Phase 8 — Attachments

```
linear issue attachment list   <issue-id>
linear issue attachment create <issue-id> --url <url> --title <title>
linear issue attachment delete <attachment-id>
```

Link external resources to issues via `attachmentCreate` / `attachmentDelete`.

### Phase 9 — Webhooks management

```
linear webhook list
linear webhook create --url <url> --team <team> --events <event,...>
linear webhook delete <webhook-id>
```

Register and manage Linear webhooks directly from the CLI using the
`webhookCreate` / `webhookUpdate` / `webhookDelete` mutations.

### Phase 10 — Interactive / TUI mode

```
linear tui
```

A terminal UI (using **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**) offering a keyboard-driven
dashboard: issue board, inline editing, cycle view, and notifications via
webhooks.

---

## API Rate Limits

| Auth type | Request limit |
|---|---|
| Personal API Key | 1,500 requests / hour |
| OAuth 2.0 token | 500 requests / hour |

The CLI will surface rate-limit headers in `--verbose` mode and back off
automatically when approaching limits.

---

## Reference Links

- [Linear Developer Docs](https://developers.linear.app/docs)
- [Linear GraphQL API](https://developers.linear.app/docs/graphql/working-with-the-graphql-api)
- [Personal API Key setup](https://linear.app/settings/api)
- [Cobra — Go CLI framework](https://github.com/spf13/cobra)
- [Viper — Go config management](https://github.com/spf13/viper)
- [hasura/go-graphql-client](https://github.com/hasura/go-graphql-client)
