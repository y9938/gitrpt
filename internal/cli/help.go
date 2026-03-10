// Package cli provides command-line interface functionality for gitrpt.
package cli

import (
	"fmt"
	"strings"

	"github.com/y9938/gitrpt/internal/i18n"
)

// HelpFormatter formats and displays help messages.
type HelpFormatter struct {
	i18n    *i18n.I18n
	version string
}

// NewHelpFormatter creates a new help formatter with the given i18n instance and version.
func NewHelpFormatter(i18n *i18n.I18n, version string) *HelpFormatter {
	return &HelpFormatter{i18n: i18n, version: version}
}

// FormatHelp returns the formatted help message.
func (h *HelpFormatter) FormatHelp() string {
	var b strings.Builder

	// Header: tool_name version - description
	fmt.Fprintf(&b, "%s %s - %s\n", h.i18n.Name(), h.version, h.i18n.Description())

	// Empty line
	fmt.Fprintln(&b)

	// Usage line
	fmt.Fprintln(&b, h.i18n.Usage())

	// Empty line
	fmt.Fprintln(&b)

	// Options section
	fmt.Fprintln(&b, "Options:")

	// Flag descriptions with aligned formatting
	flags := []struct {
		short   string
		long    string
		desc    string
		argName string
	}{
		{"", "--lang", h.i18n.FlagLang(), "LANG"},
		{"-h", "--help", h.i18n.FlagHelp(), ""},
		{"-v", "--version", h.i18n.FlagVersion(), ""},
	}

	for _, f := range flags {
		h.formatFlag(&b, f.short, f.long, f.argName, f.desc)
	}

	// Language priority note
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, h.i18n.LanguagePriority())

	return b.String()
}

// formatFlag formats a single flag line with proper alignment.
func (h *HelpFormatter) formatFlag(b *strings.Builder, short, long, argName, desc string) {
	flagStr := "  "

	if short != "" && long != "" {
		flagStr += fmt.Sprintf("%s, %s", short, long)
	} else if short != "" {
		flagStr += short
	} else {
		flagStr += "    " + long
	}

	if argName != "" {
		flagStr += fmt.Sprintf(" %s", argName)
	}

	// Pad to align descriptions
	if len(flagStr) < 30 {
		flagStr += strings.Repeat(" ", 30-len(flagStr))
	} else {
		flagStr += "  "
	}

	fmt.Fprintf(b, "%s%s\n", flagStr, desc)
}

// PrintHelp prints the help message to stdout.
func (h *HelpFormatter) PrintHelp() {
	fmt.Print(h.FormatHelp())
}

// PrintVersion prints the version string.
func (h *HelpFormatter) PrintVersion() {
	fmt.Printf("%s %s\n", h.i18n.Name(), h.version)
}
