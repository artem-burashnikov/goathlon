package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	Laps        int      `json:"laps"`
	LapLen      int      `json:"lapLen"`
	PenaltyLen  int      `json:"penaltyLen"`
	FiringLines int      `json:"firingLines"`
	Start       Time     `json:"start"`
	StartDelta  Duration `json:"startDelta"`
}

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	layouts := []string{
		"15:04:05.000",
		"15:04:05",
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
