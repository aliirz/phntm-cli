package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliirz/phntm-cli/internal/api"
	"github.com/aliirz/phntm-cli/internal/crypto"
	"github.com/aliirz/phntm-cli/internal/ui"
)

func runGet(args []string) {
	if len(args) == 0 {
		ui.Error("NO_URL_SPECIFIED")
		fmt.Fprintf(os.Stderr, "Usage: phntm get <url>\n")
		os.Exit(1)
	}

	rawURL := args[0]

	// Parse the URL to extract file ID and key
	parsed, err := url.Parse(rawURL)
	if err != nil {
		ui.Error(fmt.Sprintf("INVALID_URL: %s", err))
		os.Exit(1)
	}

	// Extract file ID from path: /f/{id}
	pathParts := strings.Split(strings.TrimPrefix(parsed.Path, "/"), "/")
	if len(pathParts) < 2 || pathParts[0] != "f" {
		ui.Error("INVALID_URL: expected format https://phntm.sh/f/{id}#{key}")
		os.Exit(1)
	}
	fileID := pathParts[1]

	// Extract key from fragment
	keyString := parsed.Fragment
	if keyString == "" {
		ui.Error("MISSING_DECRYPTION_KEY: URL must contain #key fragment")
		os.Exit(1)
	}

	// Determine API base URL from the share URL
	baseURL := fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
	envURL := os.Getenv("PHNTM_API_URL")
	if envURL != "" {
		baseURL = envURL
	}

	client := api.New(baseURL)

	// Fetch metadata
	ui.Status("LOCATING_TRANSMISSION...")
	meta, err := client.GetFileMetadata(fileID)
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	ui.FileInfo(meta.FileName, formatSize(meta.FileSize), meta.ExpiresAt)

	// Download encrypted blob
	ui.Progress("DOWNLOADING_CIPHERTEXT...")
	ciphertext, err := client.DownloadFile(fileID)
	if err != nil {
		ui.Error(fmt.Sprintf("DOWNLOAD_FAILED: %s", err))
		os.Exit(1)
	}

	// Decrypt
	ui.Progress("DECRYPTING: AES-256-GCM...")
	key, err := crypto.ImportKey(keyString)
	if err != nil {
		ui.Error(fmt.Sprintf("INVALID_KEY: %s", err))
		os.Exit(1)
	}

	plaintext, err := crypto.DecryptFile(ciphertext, key)
	if err != nil {
		ui.Error(fmt.Sprintf("DECRYPTION_FAILED: %s", err))
		os.Exit(1)
	}

	// Determine output path
	outputPath := meta.FileName
	// If file already exists, add a suffix
	if _, err := os.Stat(outputPath); err == nil {
		ext := filepath.Ext(outputPath)
		base := strings.TrimSuffix(outputPath, ext)
		outputPath = fmt.Sprintf("%s_phntm%s", base, ext)
	}

	// Write to disk
	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		ui.Error(fmt.Sprintf("WRITE_FAILED: %s", err))
		os.Exit(1)
	}

	ui.Success(fmt.Sprintf("DECRYPTION_COMPLETE: %s", outputPath))
}
