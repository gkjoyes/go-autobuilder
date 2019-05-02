package builder

import "time"

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
