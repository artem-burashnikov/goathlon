package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventString(t *testing.T) {
	fixedTime := time.Date(2025, 10, 1, 9, 30, 0, 0, time.UTC)
	tests := []struct {
		name     string
		event    Event
		expected string
	}{
		{
			name: "EventRegistered",
			event: Event{
				Timestamp:    fixedTime,
				ID:           EventRegistered,
				CompetitorID: 1,
			},
			expected: "[09:30:00.000] The competitor(1) registered",
		},
		{
			name: "EventSetStartTime",
			event: Event{
				Timestamp:    fixedTime,
				ID:           EventSetStartTime,
				CompetitorID: 2,
				Extra:        []string{"10:15:00.000"},
			},
			expected: "[09:30:00.000] The start time for the competitor(2) was set by a draw to 10:15:00.000",
		},
		{
			name: "EventOnStartLine",
			event: Event{
				Timestamp:    fixedTime,
				ID:           EventOnStartLine,
				CompetitorID: 3,
			},
			expected: "[09:30:00.000] The competitor(3) is on the start line",
		},
		{
			name: "EventStartedRace",
			event: Event{
				Timestamp:    fixedTime,
				ID:           EventStartedRace,
				CompetitorID: 4,
			},
			expected: "[09:30:00.000] The competitor(4) has started",
		},
		{
			name: "EventStartedFiringRange",
			event: Event{
				Timestamp:    fixedTime,
				ID:           EventStartedFiringRange,
				CompetitorID: 5,
				Extra:        []string{"2"},
			},
			expected: "[09:30:00.000] The competitor(5) is on the firing range(2)",
		},
		{
			name: "EventShotHit",
			event: Event{
				Timestamp:    fixedTime,
				ID:           EventShotHit,
				CompetitorID: 6,
				Extra:        []string{"3"},
			},
			expected: "[09:30:00.000] The target(3) has been hit by competitor(6)",
		},
		{
			name: "EventFinishedFiringRange",
			event: Event{
				Timestamp:    fixedTime,
				ID:           EventFinishedFiringRange,
				CompetitorID: 7,
			},
			expected: "[09:30:00.000] The competitor(7) left the firing range",
		},
		{
			name: "EventStartedPenaltyLaps",
			event: Event{
				Timestamp:    fixedTime,
				ID:           EventStartedPenaltyLaps,
				CompetitorID: 8,
			},
			expected: "[09:30:00.000] The competitor(8) entered the penalty laps",
		},
		{
			name: "EventFinishedPenaltyLaps",
			event: Event{
				Timestamp:    fixedTime,
				ID:           EventFinishedPenaltyLaps,
				CompetitorID: 9,
			},
			expected: "[09:30:00.000] The competitor(9) left the penalty laps",
		},
		{
			name: "EventFinishedLap",
			event: Event{
				Timestamp:    fixedTime,
				ID:           EventFinishedLap,
				CompetitorID: 10,
			},
			expected: "[09:30:00.000] The competitor(10) ended the main lap",
		},
		{
			name: "EventCantContinue",
			event: Event{
				Timestamp:    fixedTime,
				ID:           EventCantContinue,
				CompetitorID: 11,
				Extra:        []string{"Lost", "equipment"},
			},
			expected: "[09:30:00.000] The competitor(11) can't continue: Lost equipment",
		},
		{
			name: "EventDisqualified",
			event: Event{
				Timestamp:    fixedTime,
				ID:           EventDisqualified,
				CompetitorID: 12,
			},
			expected: "[09:30:00.000] The competitor(12) is disqualified",
		},
		{
			name: "EventFinishedRace",
			event: Event{
				Timestamp:    fixedTime,
				ID:           EventFinishedRace,
				CompetitorID: 13,
			},
			expected: "[09:30:00.000] The competitor(13) has finished",
		},
		{
			name: "UnknownEvent",
			event: Event{
				Timestamp:    fixedTime,
				ID:           99,
				CompetitorID: 14,
			},
			expected: "[09:30:00.000] IMPOSSIBLE EVENT 99 for competitor(14)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.event.String()
			assert.Equal(t, result, tt.expected)
		})
	}
}
