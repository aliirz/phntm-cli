# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`phntm-cli` is a Go CLI for [phntm.sh](https://phntm.sh) — encrypted file sharing from the terminal. Files are encrypted locally with AES-256-GCM before upload; decryption keys live in the URL fragment and never reach the server.

## Build & Run

```bash
go build -o phntm .          # Build binary
go run .                      # Run without building
go vet ./...                  # Check for issues
```

No external Go dependencies — stdlib only. Requires Go 1.24+.

## Releasing

Releases use GoReleaser via GitHub Actions (`.github/workflows/release.yml`). Push a tag (`vX.Y.Z`) to trigger a release. Version is injected at build time via ldflags into `cmd.version`.

## Architecture

Hand-rolled CLI (no cobra/urfave) — command dispatch is in `cmd/root.go` via a simple `switch` on `os.Args[1]`.

### Packages

- **`cmd/`** — CLI entry points: `send` (encrypt + upload), `get` (download + decrypt), `update` (self-update). The root also handles bare file args as implicit `send`.
- **`internal/api/`** — HTTP client for the phntm.sh API. Three-step upload flow: init (get presigned URL) → upload to storage → confirm. Download fetches metadata then blob.
- **`internal/crypto/`** — AES-256-GCM encrypt/decrypt. Wire format is `[12-byte IV][ciphertext+tag]`. Keys are base64url-encoded (no padding) for URL fragment compatibility with the web app.
- **`internal/ui/`** — Terminal output with ANSI colors. All decorative output goes to stderr; only the share URL goes to stdout, making the CLI pipe-friendly (`phntm file.pdf | pbcopy`).
- **`internal/updater/`** — Self-update via GitHub Releases. Background version check (cached daily in `~/.config/phntm/`), interactive update with binary replacement.

### Key Design Decisions

- **Pipe-friendly**: `ui.IsPiped()` detects non-TTY stdout and suppresses all decoration. Only the bare URL is printed to stdout.
- **No dependencies**: The entire CLI is stdlib-only Go — no CLI framework, no third-party HTTP client.
- **Web-compatible crypto**: The encryption format (IV prefix, base64url keys) is intentionally compatible with the phntm.sh web app's WebCrypto implementation.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PHNTM_API_URL` | `https://phntm.sh` | Override API server URL |
