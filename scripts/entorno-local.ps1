# Cargar con: . ./scripts/entorno-local.ps1
$raizPipeline = Split-Path -Parent $PSScriptRoot
$env:PATH = "$raizPipeline\.tools\node-v22.22.0-win-x64;$raizPipeline\.tools\go\bin;$raizPipeline\.tools\pnpm\node_modules\.bin;$raizPipeline\.tools\gitleaks;$raizPipeline\.tools\actionlint;$env:PATH"
$env:GOPATH = "$raizPipeline\.tools\gopath"
$env:GOCACHE = "$raizPipeline\.tools\gocache"
