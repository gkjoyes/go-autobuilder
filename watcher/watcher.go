package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/george-kj/go-autobuilder/builder"
	"github.com/george-kj/go-autobuilder/logger"
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

// Watch monitors the given path for changes.
// An initial build and run were performed before watching begins.
// Return error if any.
func (w *Watcher) Watch() error {
	logger.Info().Command("Watching", "W").Message(logger.FormattedMsg(w.dir)).Log()

	// Do first build and run.
	w.doBuildRun()

	// Check file modification in every consecutive interval.
	stopWatch := make(chan error)
	go func() {
		ticker := time.NewTicker(400 * time.Millisecond)
		for {

			// Go through all files in the directory and check status.
			err := filepath.Walk(w.dir, w.watchFunc)
			if err != nil && err != filepath.SkipDir {
				stopWatch <- err
				break
			}
			<-ticker.C
		}
	}()
	return <-stopWatch
}

// Called for each file in the directory.
// If a change is found, new build going to perform and run service over new build.
// Directories and file staring with dot are going to skip and return error if any.
func (w *Watcher) watchFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// Skip all directories and files start with dot.
	if strings.HasPrefix(filepath.Base(path), ".") {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	// Track all go files.
	if filepath.Ext(path) == ".go" {

		// If any file gets modified after the last build, the build and run again.
		if info.ModTime().After(w.b.GetLastBuild()) {
			p, err := filepath.Rel(w.dir, path)
			if err != nil {
				return err
			}
			logger.Info().Command("Modified", "M").Message(logger.FormattedMsg(p)).Log()

			w.doBuildRun()
		}
	}
	return nil
}

// doBuildRun builds and runs watching project every time files get modified.
func (w *Watcher) doBuildRun() {

	// Execute custom commands.
	err := w.r.Custom()
	if err != nil {
		logger.Error().Message(err.Error()).Log()
	}

	// Building our application.
	ok := w.b.Build()
	if ok {
		if !w.buildOnly {

			// Running our application.
			if err := w.r.Run(); err != nil {
				logger.Error().Message(err.Error()).Log()
			}
		}
	}
}
