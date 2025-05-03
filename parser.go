package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"
)

func parseEventLine(line string) Event {
	parts := strings.Fields(line)

	tsStr := strings.Trim(parts[0], "[]")
	ts, _ := time.Parse("15:04:05.000", tsStr)

	id, _ := strconv.Atoi(parts[1])

	cid, _ := strconv.Atoi(parts[2])

	return Event{Timestamp: ts, ID: id, CompetitorID: cid, Extra: parts[3:]}
}

func parseEvents(r io.Reader, ch chan<- Event) {
	defer close(ch)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ch <- parseEventLine(scanner.Text())
	}
}
