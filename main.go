package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// App version.
const v = "1.0.0"

var (
	path, name, env             string
	version, buildOnly          bool
	customCmd, runCmd, buildCmd string
)

// Read command line arguments initially.
func init() {
	flag.StringVar(&path, "p", "", "")
	flag.StringVar(&name, "n", "", "")
	flag.StringVar(&env, "e", "", "")
	flag.BoolVar(&version, "v", false, "")
	flag.BoolVar(&buildOnly, "b", false, "")
	flag.StringVar(&customCmd, "cc", "", "")
	flag.StringVar(&buildCmd, "bc", "", "")
	flag.StringVar(&runCmd, "rc", "", "")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: go-autobuilder\n")
		fmt.Fprintf(os.Stderr, "options:\n")
		fmt.Fprintf(os.Stderr, "\t-p 	The directory to be watch.\n")
		fmt.Fprintf(os.Stderr, "\t-n	Project name.\n")
		fmt.Fprintf(os.Stderr, "\t-e	Environment file path.\n")
		fmt.Fprintf(os.Stderr, "\t-v	Prints the version.\n")
		fmt.Fprintf(os.Stderr, "\t-b	Build only mode.\n")
		fmt.Fprintf(os.Stderr, "\t-cc	Custom commands to run before the build.\n")
		fmt.Fprintf(os.Stderr, "\t-bc	Custom commands to run while building.\n")
		fmt.Fprintf(os.Stderr, "\t-rc	Custom commands to run while running.\n")
	}
}

func main() {
	parseFlag()
	gracefulShutdown()
}

// Read command line arguments.
func parseFlag() {
	flag.Parse()

	// Display version.
	if version {
		fmt.Printf("go-autobuilder v%s\n", v)
		os.Exit(0)
	}

	// Set default configuration values if not provided.
	setDefaults()
}

// setDefaults set default values to configuration variables if not provided.
func setDefaults() {
	if path == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("An error occurred while getting the current working directory: %v\n", err)
		}
		path, err = filepath.Abs(dir)
		if err != nil {
			log.Fatalf("An error occurred while finding an absolute working path: %v\n", err)
		}
	} else {
		dir, err := os.Stat(path)
		if err != nil {
			log.Fatalf("Given path is not valid one: %s\n", path)
		}
		if !dir.IsDir() {
			log.Fatalf("Given path is not valid: %s: The path must be a directory\n", path)
		}
	}
}

// gracefulShutdown shutdown system cleanly if any interrupts happen.
func gracefulShutdown() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	for sig := range signalChan {
		if sig == syscall.SIGINT {
			os.Exit(0)
		}
	}
}
