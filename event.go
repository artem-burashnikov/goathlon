package main

import (
	"fmt"
	"strings"
	"time"
)

// Incoming events (1-11)
// These constants represent various events that can occur during the competition.
// They are used to identify and handle specific actions or states of competitors.
const (
	_                        = iota
	EventRegistered          // A competitor has registered for the competition
	EventSetStartTime        // The start time for a competitor has been set
	EventOnStartLine         // A competitor is on the start line
	EventStartedRace         // A competitor has started the race
	EventStartedFiringRange  // A competitor has entered the firing range
	EventShotHit             // A competitor has hit a target
	EventFinishedFiringRange // A competitor has left the firing range
	EventStartedPenaltyLaps  // A competitor has started penalty laps
	EventFinishedPenaltyLaps // A competitor has finished penalty laps
	EventFinishedLap         // A competitor has completed a main lap
	EventCantContinue        // A competitor cannot continue the race
)

// Outgoing events (32-33)
// These constants represent events that are sent out as a result of certain actions or states.
const (
	EventDisqualified = iota + 32 // A competitor has been disqualified
	EventFinishedRace             // A competitor has finished the race
)

// Event represents an event that occurs during the competition.
type Event struct {
	Timestamp    time.Time // The time when the event occurred
	ID           int       // The ID of the event (e.g., EventRegistered, EventFinishedLap)
	CompetitorID int       // The ID of the competitor associated with the event
	Extra        []string  // Additional information related to the event
}

// formatEventLog formats an Event into a human-readable log message.
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
	case EventFinishedFiringRange:
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
