package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldSkip(t *testing.T) {
	t.Run("active competitor", func(t *testing.T) {
		s := &CompetitorState{}
		assert.False(t, shouldSkip(s))
	})

	t.Run("disqualified", func(t *testing.T) {
		s := &CompetitorState{Status: StatusDisqualified}
		assert.True(t, shouldSkip(s))
	})
}

func TestGetOrCreateState(t *testing.T) {
	summary := make(Summary)

	state := getOrCreateState(summary, 1)
	assert.NotNil(t, state)
	assert.Equal(t, 1, state.CompetitorID)

	existingState := getOrCreateState(summary, 1)
	assert.Same(t, state, existingState)
}

func TestHandleStartedPenaltyLaps(t *testing.T) {
	s := &CompetitorState{CurrentHits: 3}
	evt := Event{Timestamp: time.Now()}

	err := handleStartedPenaltyLaps(evt, s)
	assert.Nil(t, err)

	assert.Equal(t, 2, s.TotalPenaltyLaps) // 5 targets - 3 hits
	assert.Equal(t, 0, s.CurrentHits)
	assert.False(t, s.CurrentPenalty.StartTime.IsZero())
}

func TestMaybeGenerateEvent(t *testing.T) {
	ts := time.Now()

	t.Run("disqualified", func(t *testing.T) {
		s := &CompetitorState{Status: StatusDisqualified}
		evt, ok := maybeGenerateEvent(Event{CompetitorID: 1, Timestamp: ts}, s)

		assert.True(t, ok)
		assert.Equal(t, EventDisqualified, evt.ID)
	})

	t.Run("finished race", func(t *testing.T) {
		s := &CompetitorState{Status: StatusFinished}
		evt, ok := maybeGenerateEvent(Event{CompetitorID: 2, Timestamp: ts}, s)

		assert.True(t, ok)
		assert.Equal(t, EventFinishedRace, evt.ID)
	})
}

func TestProcessEvents(t *testing.T) {
	baseTime := time.Date(0, 1, 1, 9, 00, 0, 0, time.UTC)

	cfg := Config{
		Laps:        2,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       Time{baseTime},
		StartDelta:  Duration{30 * time.Second},
	}

	t.Run("basic race completion", func(t *testing.T) {
		inCh := make(chan Event)
		go func() {
			defer close(inCh)
			inCh <- Event{ID: EventSetStartTime, CompetitorID: 1, Extra: []string{"09:00:00"}}
			inCh <- Event{ID: EventStartedRace, CompetitorID: 1, Timestamp: must(time.Parse(time.TimeOnly, "09:00:30"))}
			inCh <- Event{ID: EventFinishedLap, CompetitorID: 1, Timestamp: must(time.Parse(time.TimeOnly, "09:10:00"))}
			inCh <- Event{ID: EventFinishedLap, CompetitorID: 1, Timestamp: must(time.Parse(time.TimeOnly, "09:20:00"))}
		}()

		var buf bytes.Buffer
		summary := processEvents(&buf, cfg, inCh)

		assert.Len(t, summary, 1)
		s := summary[1]
		assert.Equal(t, 2, len(s.Laps))
		assert.Equal(t, StatusFinished, s.Status)
		assert.Equal(t, 20*time.Minute, s.TotalRaceDuration)
	})

	t.Run("penalty laps", func(t *testing.T) {
		inCh := make(chan Event)
		go func() {
			defer close(inCh)
			inCh <- Event{ID: EventStartedPenaltyLaps, CompetitorID: 1, Timestamp: must(time.Parse(time.TimeOnly, "09:00:30"))}
			inCh <- Event{ID: EventFinishedPenaltyLaps, CompetitorID: 1, Timestamp: must(time.Parse(time.TimeOnly, "09:10:30"))}
		}()

		var buf bytes.Buffer
		summary := processEvents(&buf, cfg, inCh)

		assert.Len(t, summary, 1)
		s := summary[1]
		assert.Equal(t, 5, s.TotalPenaltyLaps)
		assert.Equal(t, 0, s.CurrentHits)
		assert.Equal(t, 10*time.Minute, s.TotalPenaltyTime)
	})

	t.Run("can't continue", func(t *testing.T) {
		now := time.Now()
		later := now.Add(10 * time.Minute)
		inCh := make(chan Event)
		go func() {
			defer close(inCh)
			inCh <- Event{ID: EventCantContinue, CompetitorID: 1, Timestamp: later, Extra: []string{"Took", "wrong", "turn"}}
		}()

		var buf bytes.Buffer
		summary := processEvents(&buf, cfg, inCh)
		assert.Len(t, summary, 1)
		s := summary[1]
		assert.Equal(t, StatusCantContinue, s.Status)
		assert.Equal(t, later, s.LastSeenTime)
	})

	t.Run("impossible event", func(t *testing.T) {
		now := time.Now()
		inCh := make(chan Event)
		go func() {
			defer close(inCh)
			inCh <- Event{ID: -1, CompetitorID: 1, Timestamp: now}
		}()

		var buf bytes.Buffer
		summary := processEvents(&buf, cfg, inCh)
		assert.Len(t, summary, 1)
		assert.Contains(t, buf.String(), "IMPOSSIBLE")
	})
}

func TestUpdateState(t *testing.T) {
	baseTime := time.Date(0, 1, 1, 9, 00, 0, 0, time.UTC)

	cfg := Config{
		Laps:        3,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       Time{baseTime},
		StartDelta:  Duration{30 * time.Second},
	}

	tests := []struct {
		name     string
		evt      Event
		setup    func(*CompetitorState)
		validate func(*testing.T, *CompetitorState, error)
	}{
		{
			name: "set start time",
			evt:  Event{ID: EventSetStartTime, Extra: []string{"09:00:00"}},
			validate: func(t *testing.T, s *CompetitorState, err error) {
				require.NoError(t, err)
				assert.Equal(t, must(time.Parse(time.TimeOnly, "09:00:00")), s.ScheduledStartTime)
			},
		},
		{
			name: "start race late",
			evt:  Event{ID: EventStartedRace, Timestamp: must(time.Parse(time.TimeOnly, "09:01:30"))},
			setup: func(s *CompetitorState) {
				s.ScheduledStartTime = must(time.Parse(time.TimeOnly, "09:00:00"))
			},
			validate: func(t *testing.T, s *CompetitorState, err error) {
				require.NoError(t, err)
				assert.Equal(t, StatusDisqualified, s.Status)
			},
		},
		{
			name: "hit target",
			evt:  Event{ID: EventShotHit},
			validate: func(t *testing.T, s *CompetitorState, err error) {
				require.NoError(t, err)
				assert.Equal(t, 1, s.TotalHits)
				assert.Equal(t, 1, s.CurrentHits)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CompetitorState{}
			if tt.setup != nil {
				tt.setup(s)
			}
			err := updateState(cfg, tt.evt, s)
			tt.validate(t, s, err)
		})
	}
}

func TestHandleFinishedPenaltyLaps(t *testing.T) {
	t.Run("complete penalty laps", func(t *testing.T) {
		startTime := time.Now()
		finishTime := startTime.Add(5 * time.Minute)

		st := &CompetitorState{
			CurrentPenalty: Penalty{
				StartTime: startTime,
			},
			TotalPenaltyTime: 10 * time.Minute,
		}
		evt := Event{Timestamp: finishTime}

		err := handleFinishedPenaltyLaps(evt, st)
		assert.Nil(t, err)
		assert.Equal(t, 15*time.Minute, st.TotalPenaltyTime)
		assert.Equal(t, Penalty{}, st.CurrentPenalty)
	})

	t.Run("zero penalty start time", func(t *testing.T) {
		st := &CompetitorState{}
		evt := Event{Timestamp: time.Now()}

		err := handleFinishedPenaltyLaps(evt, st)
		assert.Error(t, err)
	})
}

func TestHandleCantContinue(t *testing.T) {
	now := time.Now()
	evt := Event{
		Timestamp:    now,
		CompetitorID: 1,
	}
	st := &CompetitorState{}

	err := handleCantContinue(evt, st)
	require.NoError(t, err)

	assert.Equal(t, StatusCantContinue, st.Status)
	assert.Equal(t, now, st.LastSeenTime)
}

func TestHandleCantContinueMultipleTimes(t *testing.T) {
	st := &CompetitorState{Status: StatusCantContinue}
	evt := Event{Timestamp: time.Now()}

	err := handleCantContinue(evt, st)
	require.NoError(t, err)
	assert.Equal(t, StatusCantContinue, st.Status)
}

func TestHandleInvalidSetStartTime(t *testing.T) {
	ev := Event{Extra: []string{"invalid"}}
	st := &CompetitorState{}

	err := handleSetStartTime(ev, st)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid")
}

func TestProcessEventsSkipProcessing(t *testing.T) {
	cfg := Config{Laps: 1}
	inCh := make(chan Event, 2)
	var logBuf bytes.Buffer

	t.Run("disqualified competitor", func(t *testing.T) {
		inCh <- Event{
			CompetitorID: 1,
			ID:           EventStartedRace,
			Timestamp:    time.Now(),
		}
		inCh <- Event{
			CompetitorID: 1,
			ID:           EventFinishedLap,
			Timestamp:    time.Now(),
		}
		close(inCh)

		result := processEvents(&logBuf, cfg, inCh)

		assert.Equal(t, StatusDisqualified, result[1].Status)
		assert.Contains(t, logBuf.String(), "disqualified")
	})
}

func TestProcessEventsLogError(t *testing.T) {
	cfg := Config{Laps: 1}
	var logBuf bytes.Buffer

	t.Run("error in updateState", func(t *testing.T) {
		inCh := make(chan Event, 1)
		inCh <- Event{
			CompetitorID: 1,
			ID:           EventSetStartTime,
			Extra:        []string{"invalid_time"},
		}
		close(inCh)

		summary := processEvents(&logBuf, cfg, inCh)

		assert.Contains(t, logBuf.String(), "update failed")
		assert.NotEmpty(t, summary[1])
		assert.NotContains(t, logBuf.String(), "disqualified")
		assert.NotContains(t, logBuf.String(), "finished")
	})
}
