# epgstationctl

A command-line interface for [EPGStation](https://github.com/l3tnun/EPGStation), allowing you to manage channels, programs, recordings, and reservations from the terminal.

## Features

- **Channel Management**: List and view channel information
- **Program Guide**: Browse program schedules, current broadcasts, and search programs
- **Recording Management**: Monitor active recordings and system status
- **Multiple Output Formats**: Table (default) and JSON output for scripting
- **Simple Configuration**: All configuration via command-line flags

## Installation

### Building from Source

```bash
git clone https://github.com/miscord-dev/epgstationctl.git
cd epgstationctl
go build -o bin/epgstationctl ./cmd/epgstationctl
```

### Using Go Install

```bash
go install github.com/miscord-dev/epgstationctl/cmd/epgstationctl@latest
```

## Configuration

Configuration is handled through command-line flags and environment variables. No configuration files are used.

### Environment Variables

All CLI flags can be set via environment variables with the `EPGSTATIONCTL_` prefix:

```bash
# Set server URL
export EPGSTATIONCTL_SERVER="http://192.168.1.100:8888"

# Set default output format
export EPGSTATIONCTL_OUTPUT="json"

# Set timeout
export EPGSTATIONCTL_TIMEOUT="60"

# Enable verbose output
export EPGSTATIONCTL_VERBOSE="true"

# Hide table headers
export EPGSTATIONCTL_NO_HEADER="true"
```

### Configuration Precedence

1. Command-line flags (highest priority)
2. Environment variables
3. Default values (lowest priority)

## Usage

### Global Flags

- `--server string`: EPGStation server URL (default: `http://localhost:8888`)
- `--timeout int`: Request timeout in seconds (default: 30)
- `--output, -o string`: Output format (`table` or `json`) (default: `table`)
- `--no-header`: Hide table headers
- `--verbose, -v`: Verbose output
- `--help, -h`: Show help

### Commands

#### Channels

```bash
# List all channels
epgstationctl channels list

# Show specific channel details
epgstationctl channels show 123

# JSON output
epgstationctl channels list --output=json

# Connect to different server
epgstationctl channels list --server=http://192.168.1.100:8888
```

#### Programs

```bash
# List today's programs
epgstationctl programs list

# List programs for specific date
epgstationctl programs list --date=2024-06-27

# List programs for specific channel
epgstationctl programs list --channel=123

# Show currently broadcasting programs
epgstationctl programs current

# Search programs by keyword
epgstationctl programs search "ニュース"

# Search with result limit
epgstationctl programs search "ドラマ" --limit=10
```

#### Recordings

```bash
# List current recordings
epgstationctl recordings list

# Show recording system status
epgstationctl recordings status

# List with pagination
epgstationctl recordings list --offset=10 --limit=20
```

### Advanced Usage Examples

#### Using Environment Variables

```bash
# Set up environment for JSON output and custom server
export EPGSTATIONCTL_OUTPUT="json"
export EPGSTATIONCTL_SERVER="http://192.168.1.100:8888"

# Now all commands use these defaults
epgstationctl channels list
epgstationctl programs current

# Override environment with CLI flags
epgstationctl channels list --output=table
```

#### Scripting with JSON Output

```bash
# Get channel count
epgstationctl channels list --output=json | jq length

# Find specific channel by name
epgstationctl channels list --output=json | jq '.[] | select(.name | contains("NHK"))'

# Check if any recordings are active
epgstationctl recordings status --output=json | jq '.active_count > 0'
```

#### Monitoring Scripts

```bash
#!/bin/bash
# Check recording status and send notification
STATUS=$(epgstationctl recordings status --output=json)
ACTIVE=$(echo "$STATUS" | jq '.active_count')

if [ "$ACTIVE" -gt 0 ]; then
    echo "Currently recording $ACTIVE programs"
else
    echo "No active recordings"
fi
```

## Configuration Reference

| Flag | Environment Variable | Type | Default | Description |
|------|---------------------|------|---------|-------------|
| `--server` | `EPGSTATIONCTL_SERVER` | string | `http://localhost:8888` | EPGStation server URL |
| `--timeout` | `EPGSTATIONCTL_TIMEOUT` | int | `30` | Request timeout in seconds |
| `--output, -o` | `EPGSTATIONCTL_OUTPUT` | string | `table` | Output format (`table` or `json`) |
| `--no-header` | `EPGSTATIONCTL_NO_HEADER` | bool | `false` | Hide table headers |
| `--verbose, -v` | `EPGSTATIONCTL_VERBOSE` | bool | `false` | Enable verbose output |

## API Compatibility

This tool is built using EPGStation's OpenAPI specification and supports:

- EPGStation v2.4.2 and later
- All documented API endpoints for channels, schedules, and recordings
- Proper error handling and HTTP status codes

## Development

### Prerequisites

- Go 1.24 or later
- EPGStation server for testing

### Building

```bash
# Install dependencies
go mod tidy

# Generate API client (if needed)
cd api && go generate .

# Build
go build -o bin/epgstationctl ./cmd/epgstationctl

# Run tests
go test ./...
```

### Project Structure

```
epgstationctl/
├── cmd/epgstationctl/        # Main CLI entry point
├── internal/
│   ├── commands/             # CLI command implementations
│   │   ├── channels/         # Channel commands
│   │   ├── programs/         # Program commands
│   │   └── recordings/       # Recording commands
│   ├── client/               # EPGStation API client wrapper
│   ├── config/               # Configuration management
│   ├── epgstation/           # Generated API client
│   └── output/               # Output formatting utilities
├── api/                      # OpenAPI schema and generation
└── .tmp/                     # Development documentation
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [EPGStation](https://github.com/l3tnun/EPGStation) - The excellent TV recording system this tool interfaces with
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) - OpenAPI client generation

## Related Projects

- [annict-epgstation-connector](https://github.com/miscord-dev/annict-epgstation-connector) - Sync EPGStation recordings with Annict
- [EPGStation](https://github.com/l3tnun/EPGStation) - The main EPGStation project