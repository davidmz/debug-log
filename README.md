# debug-log

[![GoDoc](https://godoc.org/github.com/davidmz/debug-log?status.svg)](https://godoc.org/github.com/davidmz/debug-log)
[![GitHub license](https://img.shields.io/github/license/davidmz/debug-log)](https://github.com/davidmz/debug-log/blob/master/LICENSE.txt)

The `debug-log` Go package (imported as `debug`) provides a thin logger inspired by the NPM's
[debug](https://www.npmjs.com/package/debug) package. It outputs the log messages depending on the
DEBUG environment variable.

Example of usage:
```go
package main

import (
	"os"

	"github.com/davidmz/debug-log"
)

func main() {
	_ = os.Setenv("DEBUG", "myproject:*")
	testLogger := debug.NewLogger("myproject:test")

	testLogger.Println("log line")
	// Output: TIMESTAMP [myproject:test] log line
}
```
