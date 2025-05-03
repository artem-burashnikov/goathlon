package main

import (
	"bufio"
	"io"
	"os"
)

const buffSize = 100

func run(eventsReader io.Reader, logWriter io.Writer, cfg Config) {
	eventCh := make(chan Event, buffSize)
	go parseEvents(eventsReader, eventCh)
	processEvents(logWriter, eventCh, cfg)
}

func main() {
	cfg := Must(loadConfig("config.conf"))

	in, _ := os.Open("example")
	defer in.Close()

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	run(in, out, cfg)
	generateReport(out, cfg)
}
