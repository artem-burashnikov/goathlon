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

func (r Result) String() string {
	var lapsStr string
	for _, lap := range r.Laps {
		if lap.Duration == 0 {
			lapsStr += "{,}, "
		} else {
			lapsStr += fmt.Sprintf("{%s, %.3f}, ", formatDuration(lap.Duration), calculateAverageSpeed(r.LapLen, lap.Duration))
		}
	}
	if len(r.Laps) > 0 {
		// Remove trailing ", "
		lapsStr = lapsStr[:len(lapsStr)-2]
	}

	var penaltyStr string
	if r.TotalPenaltyTime == 0 {
		penaltyStr = "{,}"
	} else {
		penaltyStr += "{"
		penaltyStr += formatDuration(r.TotalPenaltyTime)
		penaltyStr += fmt.Sprintf(", %.3f", calculateAverageSpeed(r.PenaltyLen*r.TotalPenaltyLaps, r.TotalPenaltyTime))
		penaltyStr += "}"
	}

	var status string
	switch r.Status {
	case StatusDisqualified:
		status = "NotStarted"
	case StatusCantContinue:
		status = "NotFinished"
	case StatusFinished:
		status = formatDuration(r.TotalRaceDuration)
	}

	return fmt.Sprintf("[%s] %d [%s] %s %d/%d",
		status,
		r.CompetitorID,
		lapsStr,
		penaltyStr,
		r.TotalHits,
		r.FiringLines*NumberOfTargets,
	)
}

func formatDuration(d time.Duration) string {
	hours := d / time.Hour
	d -= hours * time.Hour

	minutes := d / time.Minute
	d -= minutes * time.Minute

	seconds := d / time.Second
	d -= seconds * time.Second

	milliseconds := d / time.Millisecond

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}

func generateReport(w io.Writer, cfg Config, summary Summary) {
	var notStarted []Result
	var cantContinue []Result
	var finishedRace []Result

	// Sort competitors into categories based on their final state.
	for _, competitorState := range summary {
		competitorResult := Result{
			CompetitorState: competitorState,
			LapLen:          cfg.LapLen,
			PenaltyLen:      cfg.PenaltyLen,
			FiringLines:     cfg.FiringLines,
		}
		switch competitorState.Status {
		case StatusDisqualified:
			notStarted = append(notStarted, competitorResult)
		case StatusCantContinue:
			cantContinue = append(cantContinue, competitorResult)
		case StatusFinished:
			finishedRace = append(finishedRace, competitorResult)
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
