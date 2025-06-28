# epgstationctl

A comprehensive command-line interface for [EPGStation](https://github.com/l3tnun/EPGStation), providing full access to TV channel management, program scheduling, and recording operations from the terminal. Built with type-safe API integration and designed for both interactive use and automation scripts.

## Features

- **Complete Channel Management**: List all channels and view detailed channel information
- **Program Guide Operations**: 
  - Browse program schedules with date and channel filtering
  - View currently broadcasting programs
  - Search programs by keyword across title and description
- **Recording System Monitoring**: 
  - List active recordings with pagination support
  - Monitor recording system status and statistics
- **Reservation Management**:
  - List, create, update, and delete reservations with flexible column display
  - Handle reservation conflicts, skips, and overlaps
  - View reservation status summaries and counts
  - Customizable output with default, verbose, full, and custom column modes
- **Encoding Job Management**:
  - List running and queued encoding jobs with progress and flexible column display
  - Add new manual encoding jobs for recorded content
  - Cancel active or queued encoding jobs
  - Customizable output with default, verbose, full, and custom column modes
- **Multiple Output Formats**: 
  - Clean table format for human reading (default)
  - JSON output for scripting and automation
  - Configurable table headers
- **Flexible Configuration**: 
  - CLI flags for all options
  - Environment variable support with `EPGSTATIONCTL_` prefix
  - No configuration files required
- **Type-Safe API Integration**: Generated client from EPGStation's OpenAPI specification
- **Production Ready**: Comprehensive error handling and timeout management

## Installation

### Requirements

- Go 1.21 or later
- Access to an EPGStation server

### Building from Source

```bash
git clone https://github.com/miscord-dev/epgstationctl.git
cd epgstationctl
go build -o bin/epgstationctl ./cmd/epgstationctl

# Optionally install to GOPATH/bin
go install ./cmd/epgstationctl
```

### Using Go Install

```bash
go install github.com/miscord-dev/epgstationctl/cmd/epgstationctl@latest
```

### Verify Installation

```bash
epgstationctl --help
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
# List today's programs (all channels)
epgstationctl programs list

# List programs for specific date
epgstationctl programs list --date=2024-06-27

# List programs for specific channel and multiple days
epgstationctl programs list --channel=123 --days=3

# Show currently broadcasting programs
epgstationctl programs current

# Search programs by keyword (searches title and description)
epgstationctl programs search "ニュース"

# Search with result limit and half-width characters
epgstationctl programs search "ドラマ" --limit=20 --half-width=false
```

#### Recordings

```bash
# List current recordings
epgstationctl recordings list

# Show recording system status and statistics
epgstationctl recordings status

# List with pagination and custom limit
epgstationctl recordings list --offset=0 --limit=50

# Monitor recordings with verbose output
epgstationctl recordings list --verbose --half-width=true
```

#### Reservations

```bash
# List all reservations (default: essential columns only)
epgstationctl reserves list

# Show more detailed columns
epgstationctl reserves list --verbose

# Show all available columns
epgstationctl reserves list --full

# Custom column selection
epgstationctl reserves list --columns=Id,Name,StartAt,ChannelId

# Filter reservations by type and rule
epgstationctl reserves list --type=normal --rule-id=123

# Show detailed reservation information
epgstationctl reserves show 456

# Create a new reservation
epgstationctl reserves create --program-id=789 --encode-mode1=H.264

# Update reservation settings
epgstationctl reserves update 456 --encode-mode1=H.265 --parent-dir=/recordings

# Delete a reservation
epgstationctl reserves delete 456

# Remove skip status from reservation
epgstationctl reserves unskip 456

# Show reservation status summary
epgstationctl reserves status

# Trigger reservation system update
epgstationctl reserves update-system
```

#### Encoding Jobs

```bash
# List all encoding jobs (default: essential columns only)
epgstationctl encodes list

# Show more detailed columns including logs
epgstationctl encodes list --verbose

# Custom column selection
epgstationctl encodes list --columns=ID,Mode,Status,Name

# Add a new encoding job
epgstationctl encodes add --recorded-id=123 --mode=H.264 --parent-dir=/encoded

# Cancel an encoding job
epgstationctl encodes cancel 789

# Show encoding system status
epgstationctl encodes status

# Add encoding job with options
epgstationctl encodes add --recorded-id=123 --mode=H.265 --remove-original --save-same-dir
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

This tool is built using EPGStation's OpenAPI specification and provides:

- **Full API Coverage**: Channels, schedules, recordings, and program search
- **Type Safety**: Generated client ensures API compatibility 
- **Version Support**: EPGStation v2.4.2 and later
- **Error Handling**: Proper HTTP status codes and user-friendly error messages
- **Network Resilience**: Configurable timeouts and connection handling

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
│   │   ├── recordings/       # Recording commands
│   │   ├── reserves/         # Reservation commands
│   │   └── encodes/          # Encoding commands
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