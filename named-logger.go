package debug

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"
)

type justLogger interface {
	print(string) error
	checkName(string) bool
}

type namedLogger struct {
	name string

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
	_ = l.print(fmt.Sprint(v...))
}

func (l *namedLogger) Println(v ...interface{}) {
	if !l.checkName(l.name) {
		return
	}
	_ = l.print(fmt.Sprintln(v...))
}

func (l *namedLogger) Printf(format string, v ...interface{}) {
	if !l.checkName(l.name) {
		return
	}
	_ = l.print(fmt.Sprintf(format, v...))
}

func (l *namedLogger) Output() io.Writer { return l.out }

func (l *namedLogger) Fork(name string) Logger {
	return &namedLogger{name: name, parent: l}
}

// Internals

func (l *namedLogger) checkName(name string) bool {
	if l.parent != nil {
		return l.parent.checkName(name)
	}
	return parseDebugEnv(l.envSource()).check(name)
}

func (l *namedLogger) print(s string) error {
	if l.parent != nil {
		return l.parent.print(s)
	}

	// Do this early
	var prefix string
	if l.noTime {
		prefix = fmt.Sprintf("[%s]", l.name)
	} else {
		prefix = fmt.Sprintf("%s [%s]", time.Now().Format(time.RFC3339), l.name)
	}

	l.lk.Lock()
	defer l.lk.Unlock()

	l.buf.Reset()
	l.buf.WriteString(prefix)
	// Simple space check
	if s != "" && s[0] != ' ' && s[0] != '\t' {
		l.buf.WriteString(" ")
	}
	l.buf.WriteString(s)
	_, err := l.buf.WriteTo(l.out)
	return err
}
