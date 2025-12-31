# Composectl for Asustor ADM

Manage multiple Docker Compose stacks on Asustor ADM NAS with centralized configuration and automatic startup.

## Features

- Automatic stack discovery and ordering (by directory prefix)
- Centralized configuration with per-stack overrides
- Global and per-stack environment variables
- Auto-start on system boot
- Smart start/stop (uses `docker compose start` for existing containers)

## Installation

Run as root on your ADM NAS:

```sh
curl -fsSL https://raw.githubusercontent.com/kreigan/adm-composectl/main/setup.sh | sh
```

Or with wget:

```sh
wget -qO- https://raw.githubusercontent.com/kreigan/adm-composectl/main/setup.sh | sh
```

### Manual Binary Installation

To download only the latest `composectl` binary:

```sh
curl -fsSL https://github.com/kreigan/adm-composectl/releases/latest/download/composectl-linux-amd64 \
  -o /usr/local/bin/composectl && chmod +x /usr/local/bin/composectl
```

## Directory Structure

```
/volmain/.@docker_compose/
├── config.yaml          # Global configuration
├── .env                  # Global environment variables
└── stacks/
    ├── 10-traefik/
    │   ├── docker-compose.yaml
    │   ├── config.yaml   # Stack-specific config (optional)
    │   └── .env          # Stack-specific env (optional)
    ├── 20-portainer/
    │   └── docker-compose.yaml
    └── 30-nextcloud/
        └── docker-compose.yaml
```

Stacks are processed in alphabetical order by directory name. Use numeric prefixes (e.g., `10-`, `20-`) to control startup order.

## Usage

```sh
composectl <command> [stack-name]
```

| Command   | Description                                      |
|-----------|--------------------------------------------------|
| `start`   | Start all stacks (or specific stack)             |
| `stop`    | Stop containers without removing them            |
| `down`    | Stop and remove containers                       |
| `restart` | Restart stacks (stop + start)                    |
| `reload`  | Alias for restart                                |
| `list`    | Show all stacks and their status                 |

### Examples

```sh
composectl start              # Start all stacks
composectl start traefik      # Start only traefik stack
composectl stop               # Stop all stacks
composectl list               # Show stack status
```

### Flags

| Flag          | Description                              |
|---------------|------------------------------------------|
| `--base-dir`  | Base directory (default: `/volmain/.@docker_compose`) |
| `--dry-run`   | Print commands without executing         |
| `--verbose`   | Enable verbose logging                   |

## Configuration

### Global Configuration (`config.yaml`)

Located at `/volmain/.@docker_compose/config.yaml`:

```yaml
# Arguments passed to all docker compose commands
common-args: []

# Arguments for 'docker compose up'
up-args:
  - "--detach"
  - "--wait-timeout"
  - "30"
  - "--pull"
  - "always"

# Arguments for 'docker compose down'
down-args: []

# Timeout in seconds for start/stop operations
timeout: 10
```

### Per-Stack Configuration

Create `config.yaml` inside a stack directory to override global settings:

```yaml
# /volmain/.@docker_compose/stacks/10-traefik/config.yaml
up-args:
  - "--detach"
  - "--build"    # Build images for this stack

down-args:
  - "--volumes"  # Remove volumes when stopping
```

Only `up-args` and `down-args` can be overridden per-stack.

## Environment Variables

### Global Environment (`.env`)

The file `/volmain/.@docker_compose/.env` is automatically passed to all stacks via `--env-file`. Use it for shared variables:

```env
TZ=Europe/London
PUID=1000
PGID=1000
```

### Per-Stack Environment

Docker Compose automatically loads `.env` files from each stack's directory. Use these for stack-specific variables:

```env
# /volmain/.@docker_compose/stacks/20-portainer/.env
PORTAINER_PORT=9000
```

**Note:** Per-stack `.env` files are loaded by Docker Compose itself, not by composectl. They supplement (not replace) the global environment.

## Uninstall

```sh
rm -f /usr/local/etc/init.d/S99composectl.sh
rm -f /usr/local/bin/composectl
# Optionally remove data:
# rm -rf /volmain/.@docker_compose
```
