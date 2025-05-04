package main

import (
	"fmt"
	"io"
	"time"
)

type Summary = map[int]*CompetitorState

type CompetitorState struct {
	ScheduledStartTime time.Time
	ActualStartTime    time.Time
	Laps               []time.Time
	Disqualified       bool
	CantContinue       bool
	FinishedRace       bool
}

func processEvents(w io.Writer, cfg Config, inCh chan Event) Summary {
	summary := make(Summary)

	for evt := range inCh {
		logEvent(w, evt)

		state := getOrCreateState(summary, evt.CompetitorID)

		if shouldSkip(state) {
			continue
		}

		if err := updateState(cfg, evt, state); err != nil {
			logError(w, "update failed", err)
			continue
		}

		if outEvt, ok := maybeGenerateEvent(evt, state); ok {
			logEvent(w, outEvt)
		}
	}

	return summary
}

func getOrCreateState(summary Summary, id int) *CompetitorState {
	if state, exists := summary[id]; exists {
		return state
	}
	state := &CompetitorState{}
	summary[id] = state
	return state
}

func shouldSkip(state *CompetitorState) bool {
	return state.Disqualified || state.CantContinue || state.FinishedRace
}

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

func handleSetStartTime(evt Event, st *CompetitorState) error {
	t, err := time.Parse(time.TimeOnly, evt.Extra[0])
	if err != nil {
		return fmt.Errorf("invalid start time: %w", err)
	}
	st.ScheduledStartTime = t
	return nil
}

func handleStartedRace(cfg Config, evt Event, st *CompetitorState) error {
	st.ActualStartTime = evt.Timestamp

	deadline := st.ScheduledStartTime.Add(cfg.StartDelta.Duration)
	if evt.Timestamp.After(deadline) {
		st.Disqualified = true
	}
	return nil
}

func handleFinishedLap(cfg Config, evt Event, st *CompetitorState) error {
	st.Laps = append(st.Laps, evt.Timestamp)
	if len(st.Laps) == cfg.Laps {
		st.FinishedRace = true
	}
	return nil
}

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
