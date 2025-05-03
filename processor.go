package main

import (
	"io"
)

func processEvents(w io.Writer, ch <-chan Event, _ Config) {
	for e := range ch {
		logEvent(w, e)
	}
}
