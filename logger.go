package main

import (
	"fmt"
	"io"
)

// logEvent writes a formatted event log message to the provided writer.
func logEvent(w io.Writer, e Event) {
	fmt.Fprintln(w, formatEventLog(e))
}

// logError writes an error message to the provided writer.
func logError(w io.Writer, kind string, err error) {
	fmt.Fprintln(w, "[ERROR]", kind, "error has occured but was ignored:", err)
}
