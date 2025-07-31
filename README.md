<div align="center">

[![Skills](https://skillicons.dev/icons?i=go,docker&perline=2)](https://skillicons.dev)

</div>


# Benchmarker

Benchmarker is a real-time monitoring tool designed to track network bandwidth and traffic metrics for Docker containers within a specific Docker Compose stack or Swarm deployment.

## Features

- **Real-time Network Monitoring**: Continuously monitors network statistics (RX/TX bytes and packets) for Docker containers
- **Docker Stack Integration**: Automatically discovers and monitors containers within a specific Docker Compose project or Docker Swarm stack
- **Configurable Monitoring**: Customizable monitoring intervals and duration
- **Structured Logging**: Comprehensive logging with logrus structured format for easy parsing and analysis 
- **Graceful Shutdown**: Supports both manual interruption (Ctrl+C) and automatic timeout-based shutdown
- **Bandwidth Calculation**: Real-time calculation of bytes/packets per second with period tracking

## Architecture

The application follows a clean architecture with separation of concerns:

```
cmd/benchmarker/          # Application entry point
├── main.go               # Main application logic

internal/
├── containers/           # Docker container discovery
│   └── containers.go     # Container listing and filtering
├── monitor/              # Core monitoring service
│   └── monitor.go        # Monitoring orchestration
├── stats/                # Network statistics collection
│   ├── bandwidth.go      # Bandwidth calculation logic
│   └── network.go        # Docker API network stats collection
└── utils/
    ├── config/           # Configuration management
    │   └── config.go     # Environment-based configuration
    └── logger/           # Logging utilities
        └── logger.go     # Structured logging setup

configs/                  # Configuration files
└── .env                  # Environment variables

logs/                     # Log output directory
├── benchmarker.log       # Main application logs
└── monitor.log           # Monitoring-specific logs
```

## Installation

### Prerequisites

- Go 1.24.1 or later
- Docker Engine with API access
- Docker Compose (if monitoring compose stacks)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/pallandos/benchmarker.git
cd benchmarker

# Build the application
make build

# Or build manually
go build -o bin/benchmarker cmd/benchmarker/main.go
```

### Dependencies

The application uses the following key dependencies:

- `github.com/docker/docker`: Docker API client
- `github.com/sirupsen/logrus`: Structured logging
- `github.com/joho/godotenv`: Environment variable management

## Configuration

Configure the application using environment variables in `configs/.env`:

```bash
# Target Docker stack/project name
STACK_NAME=your-stack-name

# Log file directory
LOG_PATH=logs

# Monitoring frequency (Go duration format)
MONITOR_INTERVAL=1s

# Total monitoring duration (Go duration format)
MONITOR_DURATION=10m
```

### Configuration Options

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `STACK_NAME` | Docker Compose project or Swarm stack name to monitor | Required | `bitcoin-network` |
| `LOG_PATH` | Directory for log files | `logs` | `./monitoring-logs` |
| `MONITOR_INTERVAL` | Frequency of network stats collection | `5s` | `1s`, `500ms`, `2m` |
| `MONITOR_DURATION` | Total monitoring time before automatic shutdown | `60s` | `10m`, `1h`, `30s` |

## Usage

### Basic Usage

```bash
# Run with default configuration
make run

# Or run the binary directly
./bin/benchmarker
```

### Example Output

```
time="2025-07-31T15:14:58+10:00" level=info msg="Bandwidth metrics" container_id=c32c1b7933ce container_name=bitcoin-network_local-1-1.1 period_ms=4000 rx_bytes_per_sec=98.49 rx_packets_per_sec=1.25 tx_bytes_per_sec=98.49 tx_packets_per_sec=1.25
time="2025-07-31T15:14:58+10:00" level=info msg="Bandwidth metrics" container_id=824faed26908 container_name=bitcoin-network_local-1-8.1 period_ms=5005 rx_bytes_per_sec=0 rx_packets_per_sec=0 tx_bytes_per_sec=0 tx_packets_per_sec=0
```

### Monitoring Workflow

1. **Discovery**: Automatically discovers containers belonging to the specified stack
2. **Parallel Monitoring**: Launches separate goroutines for each container
3. **Data Collection**: Collects network statistics via Docker API every `MONITOR_INTERVAL`
4. **Bandwidth Calculation**: Calculates differential metrics between consecutive measurements
5. **Logging**: Outputs structured logs with network performance data
6. **Graceful Shutdown**: Stops after `MONITOR_DURATION` or manual interruption

## How It Works

### Container Discovery

The application uses Docker labels to identify containers within a stack:
- `com.docker.compose.project`: For Docker Compose projects
- `com.docker.stack.namespace`: For Docker Swarm stacks

### Network Statistics Collection

- Connects to Docker Engine API
- Calls `ContainerStats` API for each monitored container
- Aggregates network statistics across all container interfaces
- Captures: RX/TX bytes, RX/TX packets, timestamps

### Bandwidth Calculation

- Maintains previous measurement state for each container
- Calculates differential metrics: `(current - previous) / time_elapsed`
- Provides per-second rates for bytes and packets
- Tracks measurement periods for accuracy validation

### Performance Considerations

- **Concurrent Monitoring**: Each container monitored in separate goroutine
- **API Optimization**: Single API call per container per interval
- **Memory Efficiency**: Only stores previous measurement for each container
- **Error Resilience**: Individual container failures don't affect others

## Log Files

### Main Application Log (`benchmarker.log`)
- Application startup/shutdown events
- Container discovery results
- Error conditions and warnings

### Monitor Log (`monitor.log`)
- Real-time bandwidth metrics
- Container-specific monitoring events
- Network performance data

## Development

### Project Structure

The codebase follows Go best practices with clear separation of concerns:

- **`cmd/`**: Application entry points
- **`internal/`**: Private application code
- **`configs/`**: Configuration files
- **`logs/`**: Runtime log files

### Key Components

1. **Monitor Service**: Orchestrates the monitoring process
2. **Docker Monitor**: Interfaces with Docker API
3. **Bandwidth Calculator**: Computes network performance metrics
4. **Container Manager**: Handles container discovery and management
5. **Configuration Manager**: Manages environment-based settings

### Adding New Metrics

To extend monitoring capabilities:

1. Add new fields to `NetworkStats` struct in `stats/network.go`
2. Extend Docker API data collection in `GetNetworkStats()`
3. Update bandwidth calculation logic in `stats/bandwidth.go`
4. Modify logging output in monitor service

## Troubleshooting

### Common Issues

**No containers found**: Verify `STACK_NAME` matches your Docker Compose project name
```bash
docker ps --filter label=com.docker.compose.project=your-stack-name
```

**Permission denied**: Ensure user has Docker API access
```bash
sudo usermod -aG docker $USER
```

**API timeout errors**: Increase monitoring interval for heavily loaded systems
```bash
MONITOR_INTERVAL=5s
```

### Debug Mode

Enable verbose logging by modifying the logger configuration in `internal/utils/logger/logger.go`.

