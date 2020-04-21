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
