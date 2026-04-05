# starscope-cli

Official CLI for the [Starscope](https://starscope.app) review analytics API.

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

# Get analytics overview
starscope-cli analytics overview

# Output as JSON
starscope-cli reviews list --output json

# Filter reviews
starscope-cli reviews list --platform trustpilot --rating-min 4 --date-from 2026-01-01
```

## Commands

| Command | Description |
|---------|-------------|
| `auth login` | Authenticate with the Starscope API |
| `auth logout` | Remove stored authentication |
| `auth status` | Show current authentication status |
| `reviews list` | List reviews with optional filters |
| `reviews show <id>` | Show a single review |
| `insights list` | List AI-generated insights |
| `insights show <id>` | Show a single insight |
| `topics list` | List review topics |
| `topics show <id>` | Show a single topic |
| `topics reviews <id>` | List reviews for a topic |
| `connections list` | List review platform connections |
| `connections show <id>` | Show a single connection |
| `analytics overview` | Show analytics overview |
| `analytics metrics` | Show analytics metrics |
| `workspace show` | Show workspace information |
| `config set <key> <val>` | Set a configuration value |
| `config get <key>` | Get a configuration value |
| `config list` | List all configuration values |
| `version` | Show version information |

## Configuration

Configuration is stored at `~/.config/starscope-cli/config.yaml`. API tokens are stored securely in the OS keychain.

### Profiles

Use profiles to manage multiple workspaces:

```bash
# Login to a specific profile
starscope-cli auth login --profile staging

# Use a profile for a command
starscope-cli reviews list --profile staging
```

### Environment variables

| Variable | Description |
|----------|-------------|
| `STARSCOPE_TOKEN` | API token (overrides stored token) |
| `STARSCOPE_API_URL` | API base URL |
| `STARSCOPE_WORKSPACE_ID` | Workspace ID |
| `STARSCOPE_PROFILE` | Active profile name |
| `STARSCOPE_CONFIG_DIR` | Config directory override |
| `NO_COLOR` | Disable colored output |

### Output formats

```bash
starscope-cli reviews list --output table  # Default, human-readable
starscope-cli reviews list --output json   # JSON for scripting
starscope-cli reviews list --output csv    # CSV for spreadsheets
```

## API token

Generate an API token in the Starscope dashboard under **Settings > API Tokens**. API access requires a Pro plan.

## License

MIT
