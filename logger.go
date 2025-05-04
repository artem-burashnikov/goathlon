package main

import (
	"fmt"
	"io"
)

func logEvent(w io.Writer, e Event) {
	fmt.Fprintln(w, formatEventLog(e))
}

func logError(w io.Writer, kind string, err error) {
	fmt.Fprintln(w, "[ERROR]", kind, "error has occured but was ignored:", err)
}
