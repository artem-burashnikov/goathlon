package main

import (
	"fmt"
	"io"
	"time"
)

// Number of targets in the firing range.
const NumberOfTargets = 5

type CompetitorStatus int

const (
	StatusActive       CompetitorStatus = iota
	StatusDisqualified                  // Whether the competitor has been disqualified.
	StatusCantContinue                  // Whether the competitor cannot continue the race.
	StatusFinished                      // Whether the competitor has finished the race.
)

// Summary represents a mapping of competitor IDs to their states.
type Summary = map[int]*CompetitorState

type Lap struct {
	StartTime  time.Time
	FinishTime time.Time
	Duration   time.Duration
}

type Penalty struct {
	StartTime  time.Time
	FinishTime time.Time
	Duration   time.Duration
}

type CompetitorState struct {
	CompetitorID       int
	ScheduledStartTime time.Time
	ActualStartTime    time.Time
	TotalRaceDuration  time.Duration
	Laps               []Lap
	CurrentPenalty     Penalty
	TotalPenaltyTime   time.Duration
	TotalPenaltyLaps   int
	TotalHits          int
	CurrentHits        int
	Status             CompetitorStatus
	LastSeenTime       time.Time // The last time the competitor was seen.
}

// processEvents logs events, updates competitor states, and generates summary data.
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
	state := &CompetitorState{CompetitorID: id}
	summary[id] = state
	return state
}

// shouldSkip determines whether a competitor's state should prevent further processing.
func shouldSkip(state *CompetitorState) bool {
	return state.Status != StatusActive
}

// updateState updates the state of a competitor based on an incoming event.
func updateState(cfg Config, evt Event, st *CompetitorState) error {
	switch evt.ID {
	case EventSetStartTime:
		return handleSetStartTime(evt, st)

	case EventStartedRace:
		return handleStartedRace(cfg, evt, st)

	case EventShotHit:
		return handleShotHit(st)

	case EventStartedPenaltyLaps:
		return handleStartedPenaltyLaps(evt, st)

	case EventFinishedPenaltyLaps:
		return handleFinishedPenaltyLaps(evt, st)

	case EventFinishedLap:
		return handleFinishedLap(cfg, evt, st)

	case EventCantContinue:
		return handleCantContinue(evt, st)

	default:
		return nil
	}
}

// handleSetStartTime sets the scheduled start time for the competitor.
func handleSetStartTime(evt Event, st *CompetitorState) error {
	t, err := time.Parse(time.TimeOnly, evt.Extra[0])
	if err != nil {
		return fmt.Errorf("invalid start time: %w", err)
	}
	st.ScheduledStartTime = t
	return nil
}

// handleStartedPenaltyLaps starts tracking the penalty laps for the competitor.
func handleStartedPenaltyLaps(evt Event, st *CompetitorState) error {
	st.CurrentPenalty.StartTime = evt.Timestamp
	st.TotalPenaltyLaps += NumberOfTargets - st.CurrentHits
	st.CurrentHits = 0
	return nil
}

// handleFinishedPenaltyLaps stops tracking the penalty laps and updates the total penalty time.
func handleFinishedPenaltyLaps(evt Event, st *CompetitorState) error {
	if !st.CurrentPenalty.StartTime.IsZero() {
		st.CurrentPenalty.FinishTime = evt.Timestamp
		st.CurrentPenalty.Duration = evt.Timestamp.Sub(st.CurrentPenalty.StartTime)
		st.TotalPenaltyTime += st.CurrentPenalty.Duration
		st.CurrentPenalty = Penalty{}
		return nil
	}
	return fmt.Errorf("trying to finish penalty laps that werer never started")
}

// handleStartedRace sets the actual start time and initializes the first lap for the competitor.
func handleStartedRace(cfg Config, evt Event, st *CompetitorState) error {
	st.ActualStartTime = evt.Timestamp

	// Check if the competitor started within the allowed interval.
	deadline := st.ScheduledStartTime.Add(cfg.StartDelta.Duration)
	if evt.Timestamp.Before(st.ScheduledStartTime) || evt.Timestamp.After(deadline) {
		st.Status = StatusDisqualified
	}

	// Add the first lap.
	st.Laps = append(st.Laps, Lap{
		StartTime: st.ActualStartTime,
		Duration:  st.ActualStartTime.Sub(st.ScheduledStartTime),
	})

	return nil
}

// handleFinishedLap updates the finish time and duration of the current lap.
// If all laps are completed, it marks the race as finished.
func handleFinishedLap(cfg Config, evt Event, st *CompetitorState) error {
	st.Laps[len(st.Laps)-1].FinishTime = evt.Timestamp
	st.Laps[len(st.Laps)-1].Duration += evt.Timestamp.Sub(st.Laps[len(st.Laps)-1].StartTime)

	if len(st.Laps) == cfg.Laps {
		st.Status = StatusFinished
		for lap := range st.Laps {
			st.TotalRaceDuration += st.Laps[lap].Duration
		}
	} else {
		st.Laps = append(st.Laps, Lap{
			StartTime: evt.Timestamp,
		})
	}

	return nil
}

// handleShotHit increments the hit counters for the competitor.
func handleShotHit(st *CompetitorState) error {
	st.CurrentHits++
	st.TotalHits++
	return nil
}

// handleCantContinue marks the competitor as unable to continue and updates the last seen time.
func handleCantContinue(evt Event, st *CompetitorState) error {
	st.Status = StatusCantContinue
	st.LastSeenTime = evt.Timestamp
	return nil
}

// maybeGenerateEvent creates disqualification or race completion events if applicable.
func maybeGenerateEvent(incoming Event, st *CompetitorState) (Event, bool) {
	if st.Status == StatusDisqualified {
		return Event{
			Timestamp:    incoming.Timestamp,
			ID:           EventDisqualified,
			CompetitorID: incoming.CompetitorID,
			Extra:        incoming.Extra,
		}, true
	}
	if st.Status == StatusFinished {
		return Event{
			Timestamp:    incoming.Timestamp,
			ID:           EventFinishedRace,
			CompetitorID: incoming.CompetitorID,
			Extra:        incoming.Extra,
		}, true
	}
	return Event{}, false
}
