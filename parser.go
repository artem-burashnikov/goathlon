package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// parseEventLine parses a single line of input into an Event object.
// The input line is expected to have the format: [timestamp] eventID competitorID [extra...]
func parseEventLine(line string) (Event, error) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return Event{}, fmt.Errorf("invalid line: %s", line)
	}

	// Parse the timestamp from the first field.
	tsStr := strings.Trim(parts[0], "[]")
	ts, err := time.Parse(time.TimeOnly, tsStr)
	if err != nil {
		return Event{}, fmt.Errorf("invalid timestamp: %w", err)
	}

	// Parse the event ID from the second field.
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return Event{}, fmt.Errorf("invalid event id: %w", err)
	}

	// Parse the competitor ID from the third field.
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

// parseEvents reads event lines from the provided reader and sends parsed Event objects to a channel.
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
