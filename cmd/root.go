package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aliirz/phntm-cli/internal/ui"
	"github.com/aliirz/phntm-cli/internal/updater"
)

var version = "0.1.1"

func Execute() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	command := os.Args[1]

	// Background update check for non-update commands
	if command != "update" && command != "version" && command != "--version" {
		updater.CheckForUpdateQuietly(version)
		// Give the goroutine a moment to print if cached
		time.Sleep(10 * time.Millisecond)
	}

	switch command {
	case "send":
		runSend(os.Args[2:])
	case "get":
		runGet(os.Args[2:])
	case "update":
		updater.RunUpdate(version)
	case "version", "--version", "-v":
		fmt.Printf("phntm %s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		// If it looks like a file path, treat it as `send <file>`
		if !strings.HasPrefix(command, "-") {
			runSend(os.Args[1:])
		} else {
			ui.Error(fmt.Sprintf("UNKNOWN_COMMAND: %s", command))
			printUsage()
			os.Exit(1)
		}
	}
}

func printUsage() {
	ui.Banner()
	fmt.Fprintf(os.Stderr, `%sUSAGE:%s
  phntm send <file> [--expiry 1h|6h|24h]    Encrypt & upload a file
  phntm get <url>                            Download & decrypt a file
  phntm update                               Check for updates and self-update
  phntm <file>                               Shorthand for send
  phntm version                              Print version

%sEXAMPLES:%s
  phntm send report.pdf                      Upload with 24h expiry (default)
  phntm send logs.tar.gz --expiry 1h         Upload with 1h expiry
  phntm get https://phntm.sh/f/abc123#key    Download & decrypt
  phntm report.pdf | pbcopy                  Upload & copy link to clipboard

%sCONTACT:%s
  %sali@aliirz.com                              https://phntm.sh%s

`, "\033[38;2;0;255;209m", "\033[0m", "\033[38;2;0;255;209m", "\033[0m", "\033[38;2;0;255;209m", "\033[0m", "\033[38;2;85;85;85m", "\033[0m")
}
