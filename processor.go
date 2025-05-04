package main

import (
	"fmt"
	"io"
	"time"
)

// Summary represents a mapping of competitor IDs to their states.
type Summary = map[int]*CompetitorState

// CompetitorState holds the state of a competitor during the competition.
type CompetitorState struct {
	ScheduledStartTime time.Time   // The scheduled start time for the competitor.
	ActualStartTime    time.Time   // The actual start time when the competitor began the race.
	Laps               []time.Time // A list of timestamps for each completed lap.
	Disqualified       bool        // Whether the competitor has been disqualified.
	CantContinue       bool        // Whether the competitor cannot continue the race.
	FinishedRace       bool        // Whether the competitor has finished the race.
}

// processEvents processes incoming events and updates the state of competitors.
func processEvents(w io.Writer, cfg Config, inCh chan Event) Summary {
	summary := make(Summary)

	for evt := range inCh {
		logEvent(w, evt)

		state := getOrCreateState(summary, evt.CompetitorID)

		// Skip processing if the competitor is disqualified, cannot continue, or has finished.
		if shouldSkip(state) {
			continue
		}

		// Update the competitor's state based on the event.
		if err := updateState(cfg, evt, state); err != nil {
			logError(w, "update failed", err)
			continue
		}

		// Generate and log any outgoing events based on the updated state.
		if outEvt, ok := maybeGenerateEvent(evt, state); ok {
			logEvent(w, outEvt)
		}
	}

	return summary
}

// getOrCreateState retrieves the state for a competitor or creates a new one if it doesn't exist.
func getOrCreateState(summary Summary, id int) *CompetitorState {
	if state, exists := summary[id]; exists {
		return state
	}
	state := &CompetitorState{}
	summary[id] = state
	return state
}

// shouldSkip determines whether a competitor's state should prevent further processing.
func shouldSkip(state *CompetitorState) bool {
	return state.Disqualified || state.CantContinue || state.FinishedRace
}

// updateState updates the state of a competitor based on an incoming event.
func updateState(cfg Config, evt Event, st *CompetitorState) error {
	switch evt.ID {
	case EventSetStartTime:
		return handleSetStartTime(evt, st)
	case EventStartedRace:
		return handleStartedRace(cfg, evt, st)
	case EventFinishedLap:
		return handleFinishedLap(cfg, evt, st)
	default:
		return nil
	}
}

// handleSetStartTime processes an EventSetStartTime event.
func handleSetStartTime(evt Event, st *CompetitorState) error {
	t, err := time.Parse(time.TimeOnly, evt.Extra[0])
	if err != nil {
		return fmt.Errorf("invalid start time: %w", err)
	}
	st.ScheduledStartTime = t
	return nil
}

// handleStartedRace processes an EventStartedRace event.
func handleStartedRace(cfg Config, evt Event, st *CompetitorState) error {
	st.ActualStartTime = evt.Timestamp

	deadline := st.ScheduledStartTime.Add(cfg.StartDelta.Duration)
	if evt.Timestamp.After(deadline) {
		st.Disqualified = true
	}
	return nil
}

// handleFinishedLap processes an EventFinishedLap event.
func handleFinishedLap(cfg Config, evt Event, st *CompetitorState) error {
	st.Laps = append(st.Laps, evt.Timestamp)
	if len(st.Laps) == cfg.Laps {
		st.FinishedRace = true
	}
	return nil
}

// maybeGenerateEvent generates an outgoing event based on the competitor's state.
func maybeGenerateEvent(incoming Event, st *CompetitorState) (Event, bool) {
	if st.Disqualified {
		return Event{
			Timestamp:    incoming.Timestamp,
			ID:           EventDisqualified,
			CompetitorID: incoming.CompetitorID,
			Extra:        incoming.Extra,
		}, true
	}
	if st.FinishedRace {
		return Event{
			Timestamp:    incoming.Timestamp,
			ID:           EventFinishedRace,
			CompetitorID: incoming.CompetitorID,
			Extra:        incoming.Extra,
		}, true
	}
	return Event{}, false
}
