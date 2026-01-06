#!/bin/bash
set -e

PLIST="$HOME/Library/LaunchAgents/com.ukpos.print-agent.plist"
INSTALL_DIR="$HOME/ukpos-agent-runtime"

echo "🧹 Uninstalling UKPOS Print Agent..."

launchctl unload "$PLIST" 2>/dev/null || true
rm -f "$PLIST"
rm -rf "$INSTALL_DIR"

echo "✅ UKPOS Agent removed"

