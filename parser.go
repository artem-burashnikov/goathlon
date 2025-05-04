package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func parseEventLine(line string) (Event, error) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return Event{}, fmt.Errorf("invalid line: %s", line)
	}

	tsStr := strings.Trim(parts[0], "[]")
	ts, err := time.Parse(time.TimeOnly, tsStr)
	if err != nil {
		return Event{}, fmt.Errorf("invalid timestamp: %w", err)
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return Event{}, fmt.Errorf("invalid event id: %w", err)
	}

	cid, err := strconv.Atoi(parts[2])
	if err != nil {
		return Event{}, fmt.Errorf("invalid competitor id: %w", err)
	}

	return Event{
		Timestamp:    ts,
		ID:           id,
		CompetitorID: cid,
		Extra:        parts[3:],
	}, nil
}

func parseEvents(r io.Reader, w io.Writer) chan Event {
	scanner := bufio.NewScanner(r)
	eventCh := make(chan Event)
	go func() {
		defer close(eventCh)
		for scanner.Scan() {
			line := scanner.Text()
			record, err := parseEventLine(line)
			if err != nil {
				logError(w, "parseEventLine", err)
				continue
			}
			eventCh <- record
		}
		if err := scanner.Err(); err != nil {
			logError(w, "scanner", err)
		}
	}()
	return eventCh
}
