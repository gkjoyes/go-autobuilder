package watcher

import (
	"github.com/george-kj/go-autobuilder/builder"
	"github.com/george-kj/go-autobuilder/runner"
)

// Watcher struct.
// b		: Builder.
// r		: Runner.
// dir		: Directory to be watching.
// buildOnly: Build only mode.
type Watcher struct {
	b         *builder.Builder
	r         *runner.Runner
	dir       string
	buildOnly bool
}

// New create and return new watcher object.
func New(appPath, appName string, buildOnly bool, cc, bc, rc []string) *Watcher {
	return &Watcher{
		dir:       appPath,
		buildOnly: buildOnly,
		b:         builder.New(appName, appPath, bc),
		r:         runner.New(appName, appPath, rc, cc),
	}
}
