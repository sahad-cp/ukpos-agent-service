#!/bin/bash
set -e

# --------------------------------------------------
# Resolve directory where this script lives
# --------------------------------------------------
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
USER_NAME="$(whoami)"

echo "🔧 Installing UKPOS Print Agent..."
echo "📂 Installer directory: $SCRIPT_DIR"
echo "👤 User: $USER_NAME"

INSTALL_DIR="$HOME/ukpos-agent-runtime"
PLIST_DIR="$HOME/Library/LaunchAgents"

mkdir -p "$INSTALL_DIR/logs"
mkdir -p "$PLIST_DIR"

# --------------------------------------------------
# Copy files from installer directory
# --------------------------------------------------
cp "$SCRIPT_DIR/ukpos-agent" "$INSTALL_DIR/"
cp "$SCRIPT_DIR/config.json" "$INSTALL_DIR/"
cp "$SCRIPT_DIR/com.ukpos.print-agent.plist" "$PLIST_DIR/"

# --------------------------------------------------
# Replace username placeholder in plist
# --------------------------------------------------
sed -i '' "s|__USER__|$USER_NAME|g" "$PLIST_DIR/com.ukpos.print-agent.plist"

chmod +x "$INSTALL_DIR/ukpos-agent"

# --------------------------------------------------
# Load launchd service
# --------------------------------------------------
launchctl unload "$PLIST_DIR/com.ukpos.print-agent.plist" 2>/dev/null || true
launchctl load "$PLIST_DIR/com.ukpos.print-agent.plist"

echo "✅ UKPOS Agent installed and running"
echo "📄 Logs: $INSTALL_DIR/logs/agent.log"
echo "🛑 Stop: launchctl unload $PLIST_DIR/com.ukpos.print-agent.plist"