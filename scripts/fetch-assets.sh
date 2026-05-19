#!/usr/bin/env bash
# Download pinned third-party CSS/JS into internal/assets/ for //go:embed.
set -euo pipefail

PICO_VERSION="2"
MERMAID_VERSION="11"

PICO_URL="https://cdn.jsdelivr.net/npm/@picocss/pico@${PICO_VERSION}/css/pico.min.css"
MERMAID_URL="https://cdn.jsdelivr.net/npm/mermaid@${MERMAID_VERSION}/dist/mermaid.min.js"

DIR="$(cd "$(dirname "$0")/.." && pwd)/internal/assets"
mkdir -p "$DIR"

echo "Fetching Pico CSS @${PICO_VERSION} -> $DIR/pico.min.css"
curl -fsSL "$PICO_URL" -o "$DIR/pico.min.css"

echo "Fetching Mermaid @${MERMAID_VERSION} -> $DIR/mermaid.min.js"
curl -fsSL "$MERMAID_URL" -o "$DIR/mermaid.min.js"

echo "Done."
