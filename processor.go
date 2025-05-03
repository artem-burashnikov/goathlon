package main

import (
	"bufio"
	"io"
)

func processEvents(w io.Writer, ch <-chan Event, _ Config) {
	buf := bufio.NewWriter(w)
	defer buf.Flush()
	for e := range ch {
		logEvent(buf, e)
	}
}
