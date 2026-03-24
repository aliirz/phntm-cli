package ui

import (
	"fmt"
	"os"
	"strings"
)

// Box renders lines inside a Unicode box-drawing frame.
// Box-drawing characters (─│┌┐└┘) are part of the "Box Drawing" Unicode
// block (U+2500–U+257F). They're supported by virtually every modern terminal.
//
//	┌──────────────────────────┐
//	│  FILE   report.pdf       │
//	│  SIZE   3.2 MB           │
//	└──────────────────────────┘
func Box(lines []string) {
	if IsPiped() {
		return
	}

	// Find the longest line so we can pad all others to the same width.
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}
	width := maxLen + 4 // 2 chars padding on each side

	// ─ is U+2500 (horizontal line), repeated to fill the width.
	top := fmt.Sprintf("  %s┌%s┐%s", muted, strings.Repeat("─", width), reset)
	bot := fmt.Sprintf("  %s└%s┘%s", muted, strings.Repeat("─", width), reset)

	fmt.Fprintf(os.Stderr, "%s\n", top)
	for _, line := range lines {
		// Right-pad each line so the closing │ aligns.
		padded := line + strings.Repeat(" ", maxLen-len(line))
		fmt.Fprintf(os.Stderr, "  %s│%s  %s  %s│%s\n", muted, reset, padded, muted, reset)
	}
	fmt.Fprintf(os.Stderr, "%s\n", bot)
}

// URLBox renders the share URL in an emphasized box — this is the
// "hero" moment of the CLI. The URL is what the user came for.
func URLBox(url string, expiry string) {
	if IsPiped() {
		// When piped, only the bare URL goes to stdout (already handled by URL()).
		return
	}

	lines := []string{
		fmt.Sprintf("%s%s%s", bold+cyan, url, reset),
		"",
		fmt.Sprintf("%sexpires in %s%s", dim, expiry, reset),
	}

	fmt.Fprintln(os.Stderr)
	Box(lines)
	fmt.Fprintln(os.Stderr)
}

// FileInfoBox renders file metadata in a structured box.
func FileInfoBox(name string, size string, expiry string) {
	if IsPiped() {
		return
	}

	lines := []string{
		fmt.Sprintf("%sFILE%s   %s", muted, reset, name),
		fmt.Sprintf("%sSIZE%s   %s", muted, reset, size),
		fmt.Sprintf("%sEXPIRY%s %s", muted, reset, expiry),
	}

	fmt.Fprintln(os.Stderr)
	Box(lines)
}
