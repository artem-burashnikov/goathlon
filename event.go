package main

import (
	"fmt"
	"strings"
	"time"
)

type eventSource string

// Incoming events
const (
	EventRegistered = iota + 1
	EventSetStartTime
	EventOnStartLine
	EventStartedRace
	EventStartedFiringRange
	EventShotHit
	EventFinishFiringRange
	EventStartedPenaltyLaps
	EventFinishedPenaltyLaps
	EventFinishedLap
	EventCantContinue
)

// Outgoing events
const (
	EventDisqualified = 32
	EventFinishedRace = 33
)

type Event struct {
	Timestamp    time.Time
	ID           int
	CompetitorID int
	Extra        []string
	Source       eventSource
}

func formatEventLog(e Event) string {
	ts := e.Timestamp.Format("15:04:05.000")
	switch e.ID {
	case EventRegistered:
		return fmt.Sprintf("[%s] The competitor(%d) registered", ts, e.CompetitorID)
	case EventSetStartTime:
		return fmt.Sprintf("[%s] The start time for the competitor(%d) was set by a draw to %s", ts, e.CompetitorID, e.Extra[0])
	case EventOnStartLine:
		return fmt.Sprintf("[%s] The competitor(%d) is on the start line", ts, e.CompetitorID)
	case EventStartedRace:
		return fmt.Sprintf("[%s] The competitor(%d) has started", ts, e.CompetitorID)
	case EventStartedFiringRange:
		return fmt.Sprintf("[%s] The competitor(%d) is on the firing range(%s)", ts, e.CompetitorID, e.Extra[0])
	case EventShotHit:
		return fmt.Sprintf("[%s] The target(%s) has been hit by competitor(%d)", ts, e.Extra[0], e.CompetitorID)
	case EventFinishFiringRange:
		return fmt.Sprintf("[%s] The competitor(%d) left the firing range", ts, e.CompetitorID)
	case EventStartedPenaltyLaps:
		return fmt.Sprintf("[%s] The competitor(%d) entered the penalty laps", ts, e.CompetitorID)
	case EventFinishedPenaltyLaps:
		return fmt.Sprintf("[%s] The competitor(%d) left the penalty laps", ts, e.CompetitorID)
	case EventFinishedLap:
		return fmt.Sprintf("[%s] The competitor(%d) ended the main lap", ts, e.CompetitorID)
	case EventCantContinue:
		comment := strings.Join(e.Extra, " ")
		return fmt.Sprintf("[%s] The competitor(%d) can't continue: %s", ts, e.CompetitorID, comment)
	case EventDisqualified:
		return fmt.Sprintf("[%s] The competitor(%d) is disqualified", ts, e.CompetitorID)
	case EventFinishedRace:
		return fmt.Sprintf("[%s] The competitor(%d) has finished", ts, e.CompetitorID)
	default:
		return fmt.Sprintf("[%s] IMPOSSIBLE EVENT %d for competitor(%d)", ts, e.ID, e.CompetitorID)
	}
}
