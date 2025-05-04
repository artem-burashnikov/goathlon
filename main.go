package main

import (
	"bufio"
	"io"
	"os"
)

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

func run(eventsReader io.Reader, logWriter io.Writer, cfg Config) {
	eventCh := parseEvents(eventsReader, logWriter)
	competitionSummary := processEvents(logWriter, cfg, eventCh)
	generateReport(logWriter, cfg, competitionSummary)
}

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")

	cfg := must(loadConfig(cfgPath))

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	run(in, out, cfg)
}
