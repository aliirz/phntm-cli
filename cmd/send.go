package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aliirz/phantm-cli/internal/api"
	"github.com/aliirz/phantm-cli/internal/crypto"
	"github.com/aliirz/phantm-cli/internal/ui"
)

func runSend(args []string) {
	if len(args) == 0 {
		ui.Error("NO_FILE_SPECIFIED")
		fmt.Fprintf(os.Stderr, "Usage: phantm send <file> [--expiry 1h|6h|24h]\n")
		os.Exit(1)
	}

	filePath := args[0]
	expiryHours := 24 // default

	// Parse --expiry flag
	for i, arg := range args {
		if arg == "--expiry" && i+1 < len(args) {
			switch args[i+1] {
			case "1h":
				expiryHours = 1
			case "6h":
				expiryHours = 6
			case "24h":
				expiryHours = 24
			default:
				ui.Error(fmt.Sprintf("INVALID_EXPIRY: %s (use 1h, 6h, or 24h)", args[i+1]))
				os.Exit(1)
			}
		}
	}

	// Read file
	ui.Status("READING_FILE...")
	data, err := os.ReadFile(filePath)
	if err != nil {
		ui.Error(fmt.Sprintf("FAILED_TO_READ_FILE: %s", err))
		os.Exit(1)
	}

	fileName := filepath.Base(filePath)
	fileSize := int64(len(data))

	ui.FileInfo(fileName, formatSize(fileSize), fmt.Sprintf("%dH", expiryHours))

	// Encrypt
	ui.Progress("ENCRYPTING: AES-256-GCM...")
	key, err := crypto.GenerateKey()
	if err != nil {
		ui.Error(fmt.Sprintf("KEY_GENERATION_FAILED: %s", err))
		os.Exit(1)
	}

	ciphertext, err := crypto.EncryptFile(data, key)
	if err != nil {
		ui.Error(fmt.Sprintf("ENCRYPTION_FAILED: %s", err))
		os.Exit(1)
	}

	keyString := crypto.ExportKey(key)

	// Init upload
	ui.Progress("INITIATING_TRANSMISSION...")
	baseURL := os.Getenv("PHANTM_API_URL")
	client := api.New(baseURL)

	initResp, err := client.InitUpload(fileName, fileSize, expiryHours)
	if err != nil {
		ui.Error(fmt.Sprintf("TRANSMISSION_INIT_FAILED: %s", err))
		os.Exit(1)
	}

	// Upload to storage
	ui.Progress("TRANSMITTING: UPLOADING CIPHERTEXT...")
	if err := client.UploadToStorage(initResp.UploadURL, initResp.Token, ciphertext); err != nil {
		ui.Error(fmt.Sprintf("TRANSMISSION_FAILED: %s", err))
		os.Exit(1)
	}

	// Confirm
	if err := client.ConfirmUpload(initResp.ID, fileName, fileSize, expiryHours); err != nil {
		ui.Error(fmt.Sprintf("CONFIRMATION_FAILED: %s", err))
		os.Exit(1)
	}

	// Build share URL
	shareURL := fmt.Sprintf("%s/f/%s#%s", client.BaseURL, initResp.ID, keyString)

	ui.Success("TRANSMISSION_COMPLETE")
	ui.URL(shareURL)
}

func formatSize(bytes int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
