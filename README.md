# starscope-cli

Official CLI for the [Starscope](https://starscope.app) review analytics API.

## Installation

### Install script (macOS/Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/strscp/cli/main/install.sh | sh
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
starscope-cli reviews list --output table  # Default, human-readable
starscope-cli reviews list --output json   # JSON for scripting
starscope-cli reviews list --output csv    # CSV for spreadsheets
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

## API token

Generate an API token in the Starscope dashboard under [Settings > API Tokens](https://starscope.app/settings/api-tokens). API access requires a Pro plan.

## License

MIT
