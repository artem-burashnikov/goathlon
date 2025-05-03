package main

import (
	"fmt"
	"io"
)

func logEvent(w io.Writer, e Event) {
	fmt.Fprintln(w, formatEventLog(e))
}
