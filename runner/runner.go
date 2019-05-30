package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/george-kj/go-autobuilder/logger"
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
		process:        filepath.Join(appPath, appName),
		commands:       rc,
		customCommands: cc,
	}
}

// Run our application.
func (r *Runner) Run() error {

	// Terminate running process.
	err := r.terminate()
	if err != nil {
		return err
	}
	logger.Info().Command("Running", "R").Message(logger.FormattedMsg(filepath.Base(r.process) + " " + strings.Join(r.commands, " "))).Log()

	// Run app.
	r.cmd = exec.Command(r.process, r.commands...)
	r.cmd.Stderr = os.Stderr
	r.cmd.Stdout = os.Stdout
	return r.cmd.Start()
}

// Custom will execute extra command we provided.
func (r *Runner) Custom() error {
	if len(r.customCommands) == 0 {
		return nil
	}
	logger.Info().Command("Running", "R").Message(logger.FormattedMsg(strings.Join(r.customCommands, " "))).Log()

	var cmd *exec.Cmd
	if len(r.customCommands) == 1 {
		cmd = exec.Command(r.customCommands[0])
	} else {
		cmd = exec.Command(r.customCommands[0], r.customCommands[1:]...)
	}

	// Execute commands.
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error().Command("Running", "R").Message(fmt.Sprintf("%v: %v\n", err.Error(), out)).Log()
		return err
	}
	return nil
}

// terminate stop running process if anyone has a pending status.
func (r *Runner) terminate() error {
	if r.cmd == nil {
		return nil
	}

	timer := time.NewTicker(time.Second)
	done := make(chan error, 1)

	// Try to stop the running processes.
	go func() {
		done <- r.cmd.Wait()
	}()

	// Wait for almost one second for the process to stop.
	// If the process didn't exit within 1-second kill it.
	select {
	case <-timer.C:
		r.cmd.Process.Kill()
		<-done
	case err := <-done:
		if err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				return err
			}
		}
	}
	r.cmd = nil
	return nil
}
