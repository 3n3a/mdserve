#!/usr/bin/env pwsh
# Download pinned third-party CSS/JS into internal/assets/ for //go:embed.
$ErrorActionPreference = "Stop"

$PicoVersion = "2"
$MermaidVersion = "11"

$PicoUrl = "https://cdn.jsdelivr.net/npm/@picocss/pico@$PicoVersion/css/pico.min.css"
$MermaidUrl = "https://cdn.jsdelivr.net/npm/mermaid@$MermaidVersion/dist/mermaid.min.js"

$Dir = Join-Path (Split-Path -Parent $PSScriptRoot) "internal\assets"
New-Item -ItemType Directory -Force -Path $Dir | Out-Null

Write-Host "Fetching Pico CSS @$PicoVersion -> $Dir\pico.min.css"
Invoke-WebRequest -Uri $PicoUrl -OutFile (Join-Path $Dir "pico.min.css") -UseBasicParsing

Write-Host "Fetching Mermaid @$MermaidVersion -> $Dir\mermaid.min.js"
Invoke-WebRequest -Uri $MermaidUrl -OutFile (Join-Path $Dir "mermaid.min.js") -UseBasicParsing

Write-Host "Done."
