package debug

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

type justLogger interface {
	print(name string, message string) error
	checkName(name string) bool
}

type namedLogger struct {
	name   string
	prefix string

	// Configurable by options
	noTime    bool
	out       io.Writer
	envSource EnvSourceFunc

	// Internals
	lk     sync.Mutex
	buf    bytes.Buffer
	parent justLogger
}

func (l *namedLogger) Print(v ...interface{}) {
	if !l.checkName(l.name) {
		return
	}
	_ = l.print(l.name, fmt.Sprint(v...))
}

func (l *namedLogger) Println(v ...interface{}) {
	if !l.checkName(l.name) {
		return
	}
	_ = l.print(l.name, fmt.Sprintln(v...))
}

func (l *namedLogger) Printf(format string, v ...interface{}) {
	if !l.checkName(l.name) {
		return
	}
	_ = l.print(l.name, fmt.Sprintf(format, v...))
}

func (l *namedLogger) Name() string { return l.name }

func (l *namedLogger) Output() io.Writer { return l.out }

func (l *namedLogger) Fork(name string, options ...Option) Logger {
	l1 := &namedLogger{name: name, parent: l}
	for _, opt := range options {
		opt.apply(l1)
	}
	return l1
}

// Internals

func (l *namedLogger) checkName(name string) bool {
	if l.parent != nil {
		return l.parent.checkName(name)
	}
	return parseDebugEnv(l.envSource()).check(name)
}

func (l *namedLogger) print(name string, s string) error {
	if l.prefix != "" {
		s = l.prefix + " " + s
	}

	if l.parent != nil {
		return l.parent.print(name, s)
	}

	// Do this early
	var prefix string
	if l.noTime {
		prefix = fmt.Sprintf("[%s]", name)
	} else {
		prefix = fmt.Sprintf("%s [%s]", time.Now().Format(time.RFC3339), name)
	}

	l.lk.Lock()
	defer l.lk.Unlock()

	l.buf.Reset()
	l.buf.WriteString(prefix)
	// Simple space check
	if s != "" && s[0] != ' ' && s[0] != '\t' {
		l.buf.WriteString(" ")
	}
	l.buf.WriteString(strings.TrimRight(s, "\r\n"))
	l.buf.WriteString("\n")
	_, err := l.buf.WriteTo(l.out)
	return err
}
