package main

import (
	"fmt"
	"io"
	"time"

	"slices"
)

type Result struct {
	*CompetitorState
	LapLen      int
	PenaltyLen  int
	FiringLines int
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}

func (r Result) String() string {
	lapsStr := "["
	for _, lap := range r.CompetitorState.Laps {
		if lap.Duration == 0 {
			lapsStr += "{, }, "
		} else {
			lapsStr += fmt.Sprintf("{%s, %.3f}, ", formatDuration(lap.Duration), calculateAverageSpeed(r.LapLen, lap.Duration))
		}
	}
	if len(r.CompetitorState.Laps) > 0 {
		// Remove trailing ", "
		lapsStr = lapsStr[:len(lapsStr)-2]
	}
	lapsStr += "]"

	var status string
	switch {
	case r.Disqualified:
		status = "NotStarted"
	case r.CantContinue:
		status = "NotFinished"
	case r.FinishedRace:
		status = formatDuration(r.TotalRaceDuration)
	}

	return fmt.Sprintf("[%s] %d %s {%s, %.3f} %d/%d",
		status,
		r.CompetitorID,
		lapsStr,
		formatDuration(r.TotalPenaltyTime),
		calculateAverageSpeed(r.PenaltyLen*r.TotalPenaltyLaps, r.TotalPenaltyTime),
		r.TotalHits,
		r.FiringLines*NumberOfTargets,
	)
}

func generateReport(w io.Writer, cfg Config, summary Summary) {
	var notStarted []Result
	var cantContinue []Result
	var finishedRace []Result

	// Sort competitors into categories based on their final state.
	for _, v := range summary {
		switch {
		case v.Disqualified:
			r := Result{
				CompetitorState: v,
				LapLen:          cfg.LapLen,
				PenaltyLen:      cfg.PenaltyLen,
				FiringLines:     cfg.FiringLines,
			}
			notStarted = append(notStarted, r)
		case v.CantContinue:
			r := Result{
				CompetitorState: v,
				LapLen:          cfg.LapLen,
				PenaltyLen:      cfg.PenaltyLen,
				FiringLines:     cfg.FiringLines,
			}
			cantContinue = append(cantContinue, r)
		case v.FinishedRace:
			r := Result{
				CompetitorState: v,
				LapLen:          cfg.LapLen,
				PenaltyLen:      cfg.PenaltyLen,
				FiringLines:     cfg.FiringLines,
			}
			finishedRace = append(finishedRace, r)
		}
	}

	// Sort competitors within each category.
	sortByScheduledStartTime(notStarted)
	sortByLastSeenTime(cantContinue)
	sortByTotalRaceDuration(finishedRace)

	for _, v := range notStarted {
		fmt.Fprintln(w, v)
	}

	for _, v := range cantContinue {
		fmt.Fprintln(w, v)
	}

	for _, v := range finishedRace {
		fmt.Fprintln(w, v)
	}
}

func calculateAverageSpeed(distance int, duration time.Duration) float64 {
	seconds := duration.Seconds()
	if seconds == 0 {
		return 0
	}
	return float64(distance) / seconds
}

func sortByScheduledStartTime(states []Result) {
	slices.SortFunc(states, func(a, b Result) int {
		if a.ScheduledStartTime.Before(b.ScheduledStartTime) {
			return -1
		}
		if a.ScheduledStartTime.After(b.ScheduledStartTime) {
			return 1
		}
		return 0
	})
}

func sortByLastSeenTime(states []Result) {
	slices.SortFunc(states, func(a, b Result) int {
		if a.LastSeenTime.Before(b.LastSeenTime) {
			return -1
		}
		if a.LastSeenTime.After(b.LastSeenTime) {
			return 1
		}
		return 0
	})
}

func sortByTotalRaceDuration(states []Result) {
	slices.SortFunc(states, func(a, b Result) int {
		if a.TotalRaceDuration < b.TotalRaceDuration {
			return -1
		}
		if a.TotalRaceDuration > b.TotalRaceDuration {
			return 1
		}
		return 0
	})
}
