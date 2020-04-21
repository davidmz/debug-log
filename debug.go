package debug

import (
	"io"
	"os"
)

// Logger is a minimal logger interface.
type Logger interface {
	// Print writes a message via the fmt.Sprint
	Print(v ...interface{})
	// Print writes a message via the fmt.Sprintf
	Printf(format string, v ...interface{})
	// Print writes a message via the fmt.Sprintln
	Println(v ...interface{})
	// Output returns the underlying writer of the logger
	Output() io.Writer
	// Fork creates a new Logger with the same options but with a new name
	Fork(name string) Logger
}

// NewLogger creates a new named logger with the given options. If no options are provided,
// the following defaults are used: WithOutput(os.Stdout) and UseEnvVar("DEBUG").
func NewLogger(name string, options ...Option) Logger {
	l := &namedLogger{name: name}

	options = append(
		// Defaults
		Options{WithOutput(os.Stdout), UseEnvVar("DEBUG")},
		// Custom options
		options...,
	)

	for _, opt := range options {
		opt.apply(l)
	}
	return l
}
