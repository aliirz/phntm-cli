# phntm

Encrypted file sharing from the terminal. CLI for [phntm.sh](https://phntm.sh).

Zero-knowledge, end-to-end encrypted. Your files are encrypted locally before they ever leave your machine. The server only sees ciphertext. Decryption keys live in the URL fragment — never sent to any server.

## Install

```sh
# Quick install (macOS, Linux, Windows)
curl -sL https://phntm.sh/install | sh

# Or via Homebrew
brew tap aliirz/phntm
brew install phntm

# Or via npm
npm install -g @aliirz/phntm
```

Pre-built binaries for macOS (Intel & Apple Silicon), Linux (amd64 & arm64), and Windows.

## Usage

```sh
# Upload a file (24h expiry by default)
phntm send report.pdf

# Upload with custom expiry
phntm send secrets.tar.gz --expiry 1h

# Share a text note
phntm note "API key: sk-abc123"
echo "password" | phntm note

# Download & decrypt
phntm get https://phntm.sh/f/abc123#key

# Shorthand — just pass the file
phntm report.pdf

# Pipe-friendly — outputs just the URL
phntm report.pdf | pbcopy
```

### Notes

The `phntm note` command encrypts and shares text (passwords, credentials, snippets):

```sh
# Direct input
phntm note "WiFi password: Guest123" --expiry 1h

# Pipe from stdin
cat secrets.txt | phntm note

# Retrieve and view inline
phntm get https://phntm.sh/f/abc#key
# Notes render inline in terminal + save to disk

# Pipe to clipboard
phntm get https://phntm.sh/f/abc#key | pbcopy
```

Notes are stored as `note.txt` with a 10KB limit. Same encryption as files.

## How it works

1. **Encrypt** — AES-256-GCM encryption happens locally on your machine. Large files use streaming encryption (64KB chunks) for memory efficiency.
2. **Upload** — Only ciphertext is transmitted to the server
3. **Share** — You get a URL with the decryption key in the `#fragment` (never sent to server)
4. **Download** — Recipient downloads ciphertext, decrypts locally
5. **Expire** — Files self-destruct after 1, 6, or 24 hours

The encryption key only exists in the URL fragment. Browsers and HTTP clients never send fragments to servers. Not even phntm.sh can read your files.

### Wire Format (Streaming)

Large files use Rogaway's STREAM construction:

```
[PHNT magic][version][chunk_size][total_chunks][base_nonce]
[nonce_0][ciphertext_0][tag_0]...[nonce_n][ciphertext_n][tag_n]
```

Each chunk has a unique nonce derived from a base nonce using a counter + last-block flag. This provides truncation resistance and proven security (nOBE).

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PHNTM_API_URL` | `https://phntm.sh` | Custom API server URL |

## Building from source

```sh
go build -o phntm .
```

Requires Go 1.21+.

## License

MIT
