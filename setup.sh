#!/bin/sh

# Setup script for composectl on Asustor ADM NAS
# Downloads and installs composectl CLI tool with default configuration

set -e

REPO="kreigan/adm-composectl"
BASE_DIR="${COMPOSECTL_DIR:-/volmain/.@docker_compose}"
BIN_DIR="/usr/local/bin"
INIT_DIR="/usr/local/etc/init.d"

echo "==> Installing composectl to $BASE_DIR"

# Check if already installed
if [ -d "$BASE_DIR" ]; then
    echo "Warning: Directory '$BASE_DIR' already exists."
    printf "Reinstall? [y/N] "
    read -r answer
    case "$answer" in
        [yY]*) ;;
        *) echo "Aborted."; exit 1 ;;
    esac
fi

# Create directory structure
echo "==> Creating directory structure..."
mkdir -p "$BASE_DIR/stacks"

# Create default config.yaml
echo "==> Creating default configuration..."
cat > "$BASE_DIR/config.yaml" << 'EOF'
# Docker Loader Configuration
# See: https://github.com/kreigan/adm-composectl

# Arguments passed to all docker compose commands
# common-args: []

# Arguments for 'docker compose up' (default shown below)
up-args:
  - "--detach"
  - "--wait-timeout"
  - "30"
  - "--pull"
  - "always"

# Arguments for 'docker compose down'
# down-args: []

# Timeout in seconds for start/stop operations
timeout: 10
EOF

# Create empty .env file if not exists
if [ ! -f "$BASE_DIR/.env" ]; then
    echo "==> Creating empty .env file..."
    touch "$BASE_DIR/.env"
fi

# Download latest composectl binary
echo "==> Downloading composectl binary..."
LATEST_URL="https://github.com/$REPO/releases/latest/download/composectl-linux-amd64"
if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$LATEST_URL" -o "$BIN_DIR/composectl"
elif command -v wget >/dev/null 2>&1; then
    wget -q "$LATEST_URL" -O "$BIN_DIR/composectl"
else
    echo "Error: curl or wget required"
    exit 1
fi
chmod +x "$BIN_DIR/composectl"

# Download init script
echo "==> Installing init script..."
SCRIPT_URL="https://raw.githubusercontent.com/$REPO/main/S99composectl.sh"
if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$SCRIPT_URL" -o "$BASE_DIR/S99composectl.sh"
else
    wget -q "$SCRIPT_URL" -O "$BASE_DIR/S99composectl.sh"
fi
chmod +x "$BASE_DIR/S99composectl.sh"

# Create symlink for auto-start
echo "==> Enabling auto-start..."
ln -sf "$BASE_DIR/S99composectl.sh" "$INIT_DIR/S99composectl.sh"

echo ""
echo "==> Installation complete!"
echo ""
echo "Next steps:"
echo "  1. Edit $BASE_DIR/.env with global environment variables"
echo "  2. Create stack directories in $BASE_DIR/stacks/"
echo "     Example: mkdir $BASE_DIR/stacks/10-myapp"
echo "  3. Add docker-compose.yaml to each stack directory"
echo "  4. Run: composectl start"
echo ""
