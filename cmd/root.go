package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aliirz/phantm-cli/internal/ui"
)

const version = "0.1.0"

func Execute() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "send":
		runSend(os.Args[2:])
	case "get":
		runGet(os.Args[2:])
	case "version", "--version", "-v":
		fmt.Printf("phantm %s\n", version)
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
  phantm send <file> [--expiry 1h|6h|24h]    Encrypt & upload a file
  phantm get <url>                            Download & decrypt a file
  phantm <file>                               Shorthand for send
  phantm version                              Print version

%sEXAMPLES:%s
  phantm send report.pdf                      Upload with 24h expiry (default)
  phantm send logs.tar.gz --expiry 1h         Upload with 1h expiry
  phantm get https://phntm.sh/f/abc123#key    Download & decrypt
  phantm report.pdf | pbcopy                  Upload & copy link to clipboard

`, "\033[38;2;0;255;209m", "\033[0m", "\033[38;2;0;255;209m", "\033[0m")
}
