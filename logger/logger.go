package logger

import (
	"log"
	"os"
	"strings"
	"time"
)

const project = "go-autobuilder"

// Bash color macros for logging messages.
const (
	err   = "\033[38;5;196m"
	info  = "\033[38;5;266m"
	warn  = "\033[38;5;214m"
	flag  = "\033[38;5;8m"
	name  = "\033[38;5;157m"
	msg   = "\033[38;5;243m"
	reset = "\033[0m"
)

// Bash color macros for command types.
const (
	modify    = "\033[38;5;10m"
	build     = "\033[38;5;7m"
	run       = "\033[38;5;7m"
	watch     = "\033[38;5;7m"
	export    = "\033[38;5;7m"
	interrupt = "\033[38;5;7m"
)

// Log struct.
type Log struct {
	lvl string
	cmd string
	msg string
}

var bl = log.New(os.Stdout, join(name, project, flag), 0)

// Info returns new logger with logging level info.
func Info() *Log {
	return &Log{lvl: "info"}
}

// Warn returns new logger with logging level warn.
func Warn() *Log {
	return &Log{lvl: "warn"}
}

// Error returns new logger with logging level error.
func Error() *Log {
	return &Log{lvl: "error"}
}

// Message set new logging message.
func (l *Log) Message(msg ...string) *Log {
	l.msg = strings.Join(msg, " ")
	return l
}

// Command set cmd for logging.
func (l *Log) Command(cmd, code string) *Log {
	switch code {
	case "m":
		l.cmd = join(modify, code, reset, " ")
	case "b":
		l.cmd = join(build, code, reset, " ")
	case "r":
		l.cmd = join(run, code, reset, " ")
	case "i":
		l.cmd = join(interrupt, code, reset, " ")
	case "w":
		l.cmd = join(watch, code, reset, " ")
	case "e":
		l.cmd = join(export, code, reset, " ")
	}
	return l
}

// Log build new logging message and display.
func (l *Log) Log() {
	out := join("[", time.Now().Format("15:04:05"), "]")

	// Logging command.
	if l.cmd != "" {
		out = join(out, l.cmd)
	}

	// Logging level.
	switch l.lvl {
	case "error":
		out = join(out, err)
	case "info":
		out = join(out, info)
	case "warn":
		out = join(out, warn)
	}

	// Logging message.
	if l.msg != "" {
		out = join(out, l.msg)
	}
	bl.Print(join(out, reset))
}

// FormattedMsg format logging message.
func FormattedMsg(message string) string {
	return join(msg, "[", message, "]")
}

func join(msg ...string) string {
	return strings.Join(msg, "")
}
