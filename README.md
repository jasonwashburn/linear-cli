# linear-cli

A command-line interface for managing [Linear](https://linear.app) issues from
your terminal, powered by the [Linear GraphQL API](https://developers.linear.app/docs).

---

## Installation

### Prerequisites

- [Go](https://go.dev/) 1.24 or later

### Build from source

```bash
git clone https://github.com/jasonwashburn/linear-cli.git
cd linear-cli
go build -o linear .
# Optionally move to your PATH:
mv linear /usr/local/bin/linear
```

---

## Authentication

`linear-cli` uses a **Personal API Key** for authentication.

1. Go to [Linear API settings](https://linear.app/settings/api) and create a Personal API key.
2. Export it as an environment variable:

```bash
export LINEAR_API_KEY=lin_api_xxxxxxxxxxxx
```

---

## Usage

### Issue commands

#### List issues

```bash
linear issue list
linear issue list --team <team-id>
linear issue list --state "In Progress"
linear issue list --assignee <user-id>
linear issue list --json
```

#### Get issue details

```bash
linear issue get ENG-123
linear issue get ENG-123 --json
```

#### Create an issue

```bash
linear issue create --title "Fix the bug" --team <team-id>
linear issue create --title "Sub-task" --team <team-id> --parent <parent-issue-id>
linear issue create \
  --title "New feature" \
  --team <team-id> \
  --description "Details here" \
  --priority 2 \
  --assignee <user-id> \
  --json
```

Priority values: `0` = No priority, `1` = Urgent, `2` = High, `3` = Medium, `4` = Low

#### Update an issue

```bash
linear issue update ENG-123 --title "New title"
linear issue update ENG-123 --state <state-id> --priority 1
linear issue update ENG-123 --assignee <user-id> --json
```

#### Delete an issue

```bash
linear issue delete ENG-123
linear issue delete ENG-123 --yes   # skip confirmation
```

---

## JSON output

Every command supports a `--json` flag for scripting and piping:

```bash
linear issue list --json | jq '.[].title'
linear issue get ENG-123 --json | jq '.url'
```

---

## Roadmap

See [PLAN.md](PLAN.md) for the full roadmap including planned features for
teams, users, comments, labels, projects, cycles, search, and more.
