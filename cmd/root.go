package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aliirz/phntm-cli/internal/ui"
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
  phntm <file>                               Shorthand for send
  phntm version                              Print version

%sEXAMPLES:%s
  phntm send report.pdf                      Upload with 24h expiry (default)
  phntm send logs.tar.gz --expiry 1h         Upload with 1h expiry
  phntm get https://phntm.sh/f/abc123#key    Download & decrypt
  phntm report.pdf | pbcopy                  Upload & copy link to clipboard

`, "\033[38;2;0;255;209m", "\033[0m", "\033[38;2;0;255;209m", "\033[0m")
}
