package debug

import (
	"io"
	"os"
)

// Option allows to configure logger. Use them in NewLogger as arguments.
type Option interface {
	apply(*namedLogger)
}

// Options is just a list of Option. It implements the Option interface by itself.
type Options []Option

func (o Options) apply(l *namedLogger) {
	for _, opt := range o {
		opt.apply(l)
	}
}

type optionFn func(*namedLogger)

func (f optionFn) apply(l *namedLogger) { f(l) }

// WithOutput sets the output writer for the logger.
func WithOutput(out io.Writer) Option { return optionFn(func(l *namedLogger) { l.out = out }) }

// WithoutTime excludes the timestamps from log output.
func WithoutTime() Option { return optionFn(func(l *namedLogger) { l.noTime = true }) }

// UseEnvSource allows to use custom source of DEBUG value. Use UseEnvVar to obtain DEBUG from
// the environment variable.
func UseEnvSource(src EnvSourceFunc) Option {
	return optionFn(func(l *namedLogger) { l.envSource = src })
}

// UseEnvVar allow logger to use the given environment variable as DEBUG value.
func UseEnvVar(name string) Option {
	env := os.Getenv(name)
	return UseEnvSource(func() string { return env })
}

// EnvSourceFunc is a type of function that returns a DEBUG value. Use it with UseEnvSource option
// if you don't want to use standard UseEnvVar option.
type EnvSourceFunc func() string
