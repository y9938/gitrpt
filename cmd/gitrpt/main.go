// gitrpt - Git activity reporter
//
// A CLI tool for generating git activity reports.
// Supports multiple output formats and internationalization.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/y9938/gitrpt/internal/cli"
	"github.com/y9938/gitrpt/internal/i18n"
)

// Version is set at build time
var Version = "0.0.1"

// Config holds command-line configuration
type Config struct {
	Lang    string
	Help    bool
	Version bool
}

func main() {
	config := parseFlags()

	// Detect language (flag takes precedence over env)
	lang := i18n.DetectLanguage(config.Lang)

	// Initialize i18n
	i, err := i18n.New(lang)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Handle help and version
	if config.Help {
		help := cli.NewHelpFormatter(i, Version)
		help.PrintHelp()
		os.Exit(0)
	}

	if config.Version {
		fmt.Println(Version)
		os.Exit(0)
	}

	// TODO: Main application logic will go here
	fmt.Printf("%s %s - %s\n", i.Name(), Version, i.Description())
	fmt.Printf("Language: %s\n", i.Lang())
}

func parseFlags() Config {
	var config Config

	// Pre-parse to get lang flag early (needed for help message translation)
	flag.StringVar(&config.Lang, "lang", "", "")
	flag.BoolVar(&config.Help, "h", false, "")
	flag.BoolVar(&config.Help, "help", false, "")
	flag.BoolVar(&config.Version, "v", false, "")
	flag.BoolVar(&config.Version, "version", false, "")

	// Custom usage to suppress default help
	flag.Usage = func() {}

	flag.Parse()

	return config
}
