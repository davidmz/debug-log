package debug_test

import (
	"os"

	"github.com/davidmz/debug-log"
)

func ExampleNewLogger() {
	_ = os.Setenv("DEBUG", "test")
	logger := debug.NewLogger("test", debug.WithoutTime())

	logger.Println("log line")
	// Output: [test] log line
}

func ExampleNamedLogger_Fork() {
	_ = os.Setenv("DEBUG", "test*")
	logger := debug.NewLogger("test", debug.WithoutTime())
	logger = logger.Fork("test:fork")

	logger.Println("log line")
	// Output: [test:fork] log line
}
