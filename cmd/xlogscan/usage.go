// File: "usage.go"

package main

import (
	"fmt"
	"os"
)

// Print usage and exit
func printUsage() {
	fmt.Println(APP_DESCRIPTION + " v" + Version)

	if BuildTime != "" {
		fmt.Println("Build time: " + BuildTime)
	}

	fmt.Print(`
Usage: ` + APP_NAME + ` [options] [command]

Options:
  -h|--h               - Show short help about options and exit
  -help|--help|help    - Show full help and exit
  -v|--version|version - Show version and exit

  -file <log-file>     - Input log file (use stdin by default)
  -chain               - Use SumChain option
  -log-*               - Logger options

Commands:
  scan - default command
  test - generate test JSON log file

Keys (signals):
  Ctrl+C (SIGINT)  - terminate application
  Ctrl+\ (SIGQUIT) - abort application

Environment variables:
  LOG_* - Logger options
`)
	os.Exit(0)
}

// Print version and exit
func printVersion() {
	fmt.Println(APP_NAME + " version " + Version)
	os.Exit(0)
}

// EOF: "usage.go"
