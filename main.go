package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gkjoyes/go-autobuilder/logger"
	"github.com/gkjoyes/go-autobuilder/watcher"
)

// App version.
const v = "1.0.0"

var (
	appPath, appName, env, customCmd, runCmd, buildCmd string
	version, buildOnly                                 bool
)

// Read command line arguments initially.
func init() {
	flag.StringVar(&appPath, "p", "", "")
	flag.StringVar(&appName, "n", "", "")
	flag.StringVar(&customCmd, "cc", "", "")
	flag.StringVar(&buildCmd, "bc", "", "")
	flag.StringVar(&runCmd, "rc", "", "")
	flag.StringVar(&env, "e", "", "")
	flag.BoolVar(&version, "v", false, "")
	flag.BoolVar(&buildOnly, "b", false, "")
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
	go gracefulShutdown()
	parseFlag()

	// Set env variables, If env path was provided.
	if env != "" {
		setEnv(env)
	}

	// Create new watcher object.
	w := watcher.New(
		appPath,
		appName,
		buildOnly,
		prepareCommands(customCmd),
		prepareCommands(buildCmd),
		prepareCommands(runCmd),
	)

	// Watching given path for changes.
	if err := w.Watch(); err != nil {
		logger.Error().Message(err.Error()).Log()
		os.Exit(1)
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
}

// setDefaults set default values to configuration variables if not provided.
func setDefaults() {
	if appPath == "" {
		dir, err := os.Getwd()
		if err != nil {
			logger.Error().Message(fmt.Sprintf("An error occurred while getting the current working directory: %v\n", err)).Log()
			os.Exit(1)
		}
		appPath, err = filepath.Abs(dir)
		if err != nil {
			logger.Error().Message(fmt.Sprintf("An error occurred while finding an absolute working path: %v\n", err)).Log()
			os.Exit(1)
		}
	} else {
		dir, err := os.Stat(appPath)
		if err != nil {
			logger.Error().Message(fmt.Sprintf("Given path is not valid one: %s\n", appPath)).Log()
			os.Exit(1)
		}
		if !dir.IsDir() {
			logger.Error().Message(fmt.Sprintf("Given path is not valid: %s: The path must be a directory\n", appPath)).Log()
			os.Exit(1)
		}
	}

	if appName == "" {
		appName = filepath.Base(appPath)
	}
}

// prepareCommands split commands and remove unwanted space between these.
func prepareCommands(c string) []string {

	cmd := strings.Split(strings.TrimSpace(c), " ")
	final := make([]string, 0, len(cmd))
	keys := make(map[string]bool, 0)

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
			logger.Warn().Command("Interrupt", "I").Message("Exiting...").Log()
			os.Exit(0)
		}
	}
}

// setEnv setup configuration variables from given env file.
func setEnv(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Warn().Message(fmt.Sprintf("An error occurred while reading env file: %v\n", err)).Log()
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") {
			params := strings.SplitN(line, "=", 2)
			if len(params) == 2 {
				os.Setenv(strings.TrimSpace(params[0]), strings.TrimSpace(params[1]))
			}
		}
	}
}
