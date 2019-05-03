package builder

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/george-kj/go-autobuilder/logger"
)

// Builder struct.
// app		: Application name.
// dir		: The directory from which build will take.
// commands : Commands to be run while building.
// lastBuild: When the last build was taken.
type Builder struct {
	app       string
	dir       string
	commands  []string
	lastBuild time.Time
}

// New create and return new builder object.
func New(appName, appPath string, commands []string) *Builder {
	return &Builder{
		app:      appName,
		dir:      appPath,
		commands: commands,
	}
}

// GetLastBuild returns last build time.
func (b *Builder) GetLastBuild() time.Time {
	return b.lastBuild
}

// SetLastBuild update lastest build time.
func (b *Builder) SetLastBuild(time time.Time) {
	b.lastBuild = time
}

// Build take build of our application.
func (b *Builder) Build() bool {

	// Create build commands with arguments.
	commands := []string{"go", "build", "-o", b.app}

	if len(b.commands) != 0 {
		commands = append(commands, b.commands...)
	}
	logger.Info().Command("Build", "b").Message(logger.FormattedMsg(strings.Join(commands, " "))).Log()

	// Execute build commands.
	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Dir = b.dir
	out, err := cmd.CombinedOutput()

	// Update last build time.
	b.SetLastBuild(time.Now())

	if err != nil {
		logger.Error().Command("Build", "b").Message(fmt.Sprintf("Failed: %v: %v\n", err.Error(), out)).Log()
		return false
	}
	return true
}
