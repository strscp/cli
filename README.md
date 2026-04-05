# starscope-cli

The official CLI for [Starscope](https://starscope.app) -- review intelligence on autopilot.

Track sentiment, compare ratings, and surface insights across Trustpilot, Google, Feedback Company, and more. From your terminal.

[Create a free account](https://starscope.app) to get started. No credit card required. Then [generate an API token](https://starscope.app/settings/api-tokens) to connect the CLI.

## Installation

### Install script (macOS/Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/strscp/cli/main/install.sh | sh
```

### npm

```bash
npm install -g @strscp/cli
```

### GitHub Releases

Download pre-built binaries from [GitHub Releases](https://github.com/strscp/cli/releases).

### Build from source

```bash
go install github.com/strscp/cli/cmd/starscope-cli@latest
```

## Quick start

```bash
# Authenticate with your API token
starscope-cli auth login

# List your reviews
starscope-cli reviews list

# Filter reviews by platform and rating
starscope-cli reviews list --platform trustpilot --rating-min 4

# Get analytics overview
starscope-cli analytics overview

# Output as JSON for scripting
starscope-cli reviews list --output json | jq '.data[].rating'

# Use aliases for speed
starscope-cli r ls --all
```

## Commands

| Command | Aliases | Description |
|---------|---------|-------------|
| `reviews list` | `r ls` | List reviews with optional filters |
| `reviews show <id>` | `r s` | Show a single review |
| `insights list` | `i ls` | List AI-generated insights |
| `insights show <id>` | `i s` | Show a single insight |
| `topics list` | `t ls` | List review topics |
| `topics show <id>` | `t s` | Show a single topic |
| `topics reviews <id>` | | List reviews for a topic |
| `connections list` | `c ls` | List review platform connections |
| `connections show <id>` | `c s` | Show a single connection |
| `analytics overview` | `a overview` | Show analytics overview |
| `analytics metrics` | `a metrics` | Show analytics metrics |
| `workspace show` | | Show workspace information |
| `auth login` | | Authenticate with the Starscope API |
| `auth logout` | | Remove stored authentication |
| `auth status` | | Show current authentication status |
| `open [target]` | | Open Starscope in the browser |
| `config set <key> <val>` | | Set a configuration value |
| `config get <key>` | | Get a configuration value |
| `config list` | | List all configuration values |
| `version` | | Show version information |
| `completion [shell]` | | Generate shell completion scripts |

Running a parent command without a subcommand defaults to `list` (e.g. `starscope-cli reviews` runs `reviews list`).

### Filtering

```bash
# Reviews
starscope-cli reviews list --platform trustpilot
starscope-cli reviews list --rating-min 4 --rating-max 5
starscope-cli reviews list --date-from 2026-01-01 --date-to 2026-03-31
starscope-cli reviews list --connection-id 5

# Insights
starscope-cli insights list --type trend --severity critical

# Analytics
starscope-cli analytics metrics --period 30 --connection-id 5
```

### Pagination

```bash
starscope-cli reviews list --per-page 50 --page 2
starscope-cli reviews list --all  # Fetch all pages (rate-limit aware)
```

### Open in browser

```bash
starscope-cli open              # Dashboard
starscope-cli open tokens       # API tokens settings
starscope-cli open connections  # Connections page
starscope-cli open settings     # Settings page
```

## Output formats

```bash
starscope-cli reviews list --output table     # Default, human-readable
starscope-cli reviews list --output json      # JSON for scripting
starscope-cli reviews list --output csv       # CSV for spreadsheets
starscope-cli reviews list --output markdown  # Markdown for GitHub/Slack
```

Table output includes color-coded ratings (green/yellow/red) and severity levels. Disable with `--no-color` or `NO_COLOR=1`.

### Watch mode

Re-run any command on an interval to monitor changes:

```bash
starscope-cli reviews list --watch 30s             # Refresh every 30 seconds
starscope-cli analytics overview --watch 1m        # Refresh every minute
starscope-cli insights list --watch 5m --quiet     # Quiet refresh
```

### Shell completions

Generate completions for your shell:

```bash
# Bash
source <(starscope-cli completion bash)

# Zsh
starscope-cli completion zsh > "${fpath[1]}/_starscope-cli"

# Fish
starscope-cli completion fish | source
```

## Automation & scripting

The CLI is designed to work well with AI agents, shell scripts, and CI pipelines.

### Auto-JSON when piped

When stdout is not a TTY (piped or redirected), the output format automatically switches to JSON. No flag needed:

```bash
starscope-cli reviews list | jq '.data[].rating'  # auto-JSON
starscope-cli reviews list --output table | cat    # explicit override
```

### Raw API responses

Use `--raw` to get the unprocessed API response -- useful for debugging or when you need fields the CLI doesn't display:

```bash
starscope-cli reviews show 42 --raw
starscope-cli reviews list --raw | jq '.meta'
```

### Quiet mode

Use `--quiet` / `-q` to suppress all non-data output (pagination footers, progress indicators):

```bash
starscope-cli reviews list --all --quiet
```

### Non-interactive auth

Authenticate without prompts by passing the token directly:

```bash
starscope-cli auth login --token "$STARSCOPE_TOKEN"
```

Or set the environment variable to skip auth entirely:

```bash
export STARSCOPE_TOKEN=your-token
starscope-cli reviews list
```

### Structured errors

When output is JSON, errors are written to stderr as structured JSON:

```json
{"error": "not_found", "message": "Resource not found", "exit_code": 3}
```

### Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Authentication error (401) |
| 2 | Forbidden (403) |
| 3 | Not found (404) |
| 4 | Validation error (422) |
| 5 | Rate limited (429) |
| 6 | Other API error |
| 7 | Configuration / usage error |

### Example: agent workflow

```bash
#!/bin/bash
set -e

export STARSCOPE_TOKEN="$1"

# Fetch all reviews as JSON, suppress progress output
reviews=$(starscope-cli reviews list --all --quiet)

# Check exit code and parse
if [ $? -eq 0 ]; then
  echo "$reviews" | jq '[.[] | select(.rating <= 2)]' > low_rated.json
fi

# Get raw API response for a specific review
starscope-cli reviews show 42 --raw | jq '.data.sentiment_score'
```

## Configuration

Configuration is stored at `~/.config/starscope-cli/config.yaml`. API tokens are stored securely in the OS keychain.

### Profiles

Use profiles to manage multiple workspaces:

```bash
starscope-cli auth login --profile staging
starscope-cli reviews list --profile staging
```

### Environment variables

| Variable | Description |
|----------|-------------|
| `STARSCOPE_TOKEN` | API token (overrides stored token) |
| `STARSCOPE_API_URL` | API base URL |
| `STARSCOPE_WORKSPACE_ID` | Workspace ID |
| `STARSCOPE_CONFIG_DIR` | Config directory override |
| `NO_COLOR` | Disable colored output |

## License

MIT
