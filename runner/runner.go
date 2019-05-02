package runner

import (
	"os/exec"
	"path/filepath"
)

// Runner struct.
// cmd	   		 : cmd represents an external command being prepared or run.
// process 		 : Absolute path to process for running.
// commands		 : External commands to be included while running.
// customCommands: Custom commands for other purposes like formatting, linting etc.
type Runner struct {
	cmd            *exec.Cmd
	process        string
	commands       []string
	customCommands []string
}

// New create and return new runner object.
func New(appName, appPath string, rc, cc []string) *Runner {
	return &Runner{
		process:        filepath.Join(appName, appPath),
		commands:       rc,
		customCommands: cc,
	}
}
