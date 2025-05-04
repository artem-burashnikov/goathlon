package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Config represents the configuration for the biathlon competition.
type Config struct {
	Laps        int      `json:"laps"`
	LapLen      int      `json:"lapLen"`
	PenaltyLen  int      `json:"penaltyLen"`
	FiringLines int      `json:"firingLines"`
	Start       Time     `json:"start"`
	StartDelta  Duration `json:"startDelta"`
}

// Time is a custom type that embeds time.Time and provides custom JSON unmarshaling.
// It is used to parse the 'Start' field in the configuration.
type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	layouts := []string{
		"15:04:05.000", // Format with milliseconds
		"15:04:05",     // Format without milliseconds
	}
	var err error
	for _, layout := range layouts {
		var parsed time.Time
		parsed, err = time.Parse(layout, s)
		if err == nil {
			t.Time = parsed
			return nil
		}
	}
	return fmt.Errorf("invalid 'start' format: %s", s)
}

// Duration is a custom type that embeds time.Duration and provides custom JSON unmarshaling.
// It is used to parse the 'StartDelta' field in the configuration.
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("15:04:05", s)
	if err == nil {
		d.Duration = time.Duration(t.Hour())*time.Hour +
			time.Duration(t.Minute())*time.Minute +
			time.Duration(t.Second())*time.Second
		return nil
	}
	return fmt.Errorf("invalid 'startDelta' format: %s", s)
}

// loadConfig reads and parses the configuration file from the given path.
// It returns a Config object or an error if the file cannot be read or parsed.
func loadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}
