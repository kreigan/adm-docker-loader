# Development Guide

Development documentation for the `composectl` CLI tool.

## Prerequisites

- Go 1.23+
- Docker with Docker Compose plugin
- golangci-lint v2

## Building

```sh
go build -o composectl
```

### Cross-compilation

```sh
# Linux AMD64 (for ADM NAS)
GOOS=linux GOARCH=amd64 go build -o composectl-linux-amd64

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o composectl-linux-arm64
```

## Testing

```sh
# Run all tests
go test ./...

# Verbose output
go test -v ./...

# With coverage
go test -cover ./...
```

## Linting

```sh
golangci-lint run
```

The project uses golangci-lint v2 with 40+ linters enabled. See [.golangci.yml](.golangci.yml) for configuration.

## Project Structure

```
├── cmd/                     # CLI commands (Cobra)
│   ├── root.go              # Root command, flags, config
│   ├── runner.go            # ActionRunner helper
│   ├── start.go             # start command
│   ├── stop.go              # stop command
│   ├── down.go              # down command
│   ├── restart.go           # restart command
│   ├── reload.go            # reload command (alias)
│   └── list.go              # list command
├── internal/loader/         # Core business logic
│   ├── types.go             # Stack, StackStatus, Action types
│   ├── docker.go            # DockerExecutor, ComposeClient
│   ├── repository.go        # StackRepository (discovery)
│   ├── manager.go           # StackManager (orchestration)
│   ├── config.go            # Config, StackConfig loading
│   ├── logger.go            # File/console logging
│   └── *_test.go            # Unit tests
├── main.go                  # Entry point
├── S99composectl.sh         # Init script wrapper
└── setup.sh                 # Installation script
```

## Architecture

The codebase follows these patterns:

- **Repository Pattern**: `StackRepository` handles stack discovery and retrieval
- **Dependency Injection**: `DockerExecutor` interface allows mocking Docker commands
- **Composition**: `StackManager` composes repository and compose client

### Key Types

| Type | Description |
|------|-------------|
| `Stack` | Represents a Docker Compose stack (name, dir, status) |
| `StackStatus` | Enum: `running`, `stopped`, `down` |
| `Action` | Enum: `start`, `stop`, `down`, `restart`, `reload`, `list` |
| `DockerExecutor` | Interface for running Docker commands |
| `ComposeClient` | High-level Docker Compose operations |
| `StackRepository` | Discovers and retrieves stacks |
| `StackManager` | Orchestrates actions across stacks |

## Stack Naming Convention

Stack directories must follow the pattern `NN-stack-name`:

| Directory | Project Name |
|-----------|--------------|
| `10-nginx` | `nginx` |
| `20-my-app` | `my-app` |
| `30-db_server` | `db_server` |

The numeric prefix controls startup order. Directories not matching this pattern are ignored.

## Docker Command Construction

For a stack at `stacks/10-nginx`:

**Start (new stack):**
```sh
docker compose \
  --env-file /volmain/.@docker_compose/.env \
  --project-directory /volmain/.@docker_compose/stacks/10-nginx \
  --project-name nginx \
  up --detach --wait-timeout 30 --pull always
```

**Start (existing containers):**
```sh
docker compose ... start --timeout 10
```

**Stop:**
```sh
docker compose ... stop --timeout 10
```

**Down:**
```sh
docker compose ... down
```

## Adding a New Command

1. Create `cmd/newcmd.go`:
   ```go
   package cmd

   import "github.com/spf13/cobra"

   var newCmd = &cobra.Command{
       Use:   "newcmd [stack]",
       Short: "Description",
       Args:  cobra.MaximumNArgs(1),
       RunE:  RunAction("newcmd"),
   }

   func init() {
       rootCmd.AddCommand(newCmd)
   }
   ```

2. Add action to `internal/loader/types.go`:
   ```go
   ActionNewCmd Action = "newcmd"
   ```

3. Handle in `internal/loader/manager.go` `ExecuteAction()` switch.

## Release Process

1. Update version in `cmd/root.go`
2. Build binaries for target platforms
3. Create GitHub release with binaries named:
   - `composectl-linux-amd64`
   - `composectl-linux-arm64`
