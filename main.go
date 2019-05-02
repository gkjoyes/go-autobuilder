package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/george-kj/go-autobuilder/watcher"
)

// App version.
const v = "1.0.0"

var (
	appPath, appName, env       string
	version, buildOnly          bool
	customCmd, runCmd, buildCmd string
)

// Read command line arguments initially.
func init() {
	flag.StringVar(&appPath, "p", "", "")
	flag.StringVar(&appName, "n", "", "")
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
	go gracefulShutdown()

	// Set env variables, If env path was provided.
	if env != "" {
		setEnv(env)
	}
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

	// Create new watcher object.
	w := watcher.New(
		appPath,
		appName,
		buildOnly,
		prepareCommands(customCmd),
		prepareCommands(buildCmd),
		prepareCommands(runCmd),
	)

	_ = w
}

// setDefaults set default values to configuration variables if not provided.
func setDefaults() {
	if appPath == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("An error occurred while getting the current working directory: %v\n", err)
		}
		appPath, err = filepath.Abs(dir)
		if err != nil {
			log.Fatalf("An error occurred while finding an absolute working path: %v\n", err)
		}
	} else {
		dir, err := os.Stat(appPath)
		if err != nil {
			log.Fatalf("Given path is not valid one: %s\n", appPath)
		}
		if !dir.IsDir() {
			log.Fatalf("Given path is not valid: %s: The path must be a directory\n", appPath)
		}
	}

	if appName == "" {
		appName = filepath.Base(appPath)
	}
}

// prepareCommands split commands and remove unwanted space between these.
func prepareCommands(c string) []string {

	var (
		cmd   = strings.Split(strings.TrimSpace(c), " ")
		final = make([]string, 0, len(cmd))
		keys  = make(map[string]bool, 0)
	)

	// Eliminate duplicate commands and empty spaces.
	for _, v := range cmd {
		if got := keys[v]; !got && v != "" {
			keys[v] = true
			final = append(final, v)
		}
	}

	return final
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

// setEnv setup configuration variables from given env file.
func setEnv(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("An error occurred while reading env file: %v\n", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") {
			params := strings.Split(line, "=")
			if len(params) == 2 {
				os.Setenv(strings.TrimSpace(params[0]), strings.TrimSpace(params[1]))
			}
		}
	}
}
