package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		configJSON string
		wantConfig Config
		wantErr    bool
	}{
		{
			name: "valid config without milliseconds",
			configJSON: `{
			    "laps": 2,
				"lapLen": 3651,
				"penaltyLen": 50,
				"firingLines": 1,
				"start": "09:30:00",
				"startDelta": "00:00:30"
			}`,
			wantConfig: Config{
				Laps:        2,
				LapLen:      3651,
				PenaltyLen:  50,
				FiringLines: 1,
				Start:       Time{parseTime(t, time.TimeOnly, "09:30:00")},
				StartDelta:  Duration{30 * time.Second},
			},
			wantErr: false,
		},
		{
			name: "valid config with milliseconds",
			configJSON: `{
				"laps": 3,
				"lapLen": 1000,
				"penaltyLen": 25,
				"firingLines": 2,
				"start": "12:34:56.789",
				"startDelta": "00:02:03.456"
			}`,
			wantConfig: Config{
				Laps:        3,
				LapLen:      1000,
				PenaltyLen:  25,
				FiringLines: 2,
				Start:       Time{parseTime(t, "15:04:05.000", "12:34:56.789")},
				StartDelta:  Duration{2*time.Minute + 3*time.Second},
			},
			wantErr: false,
		},
		{
			name: "invalid start time format",
			configJSON: `{
				"start": "invalid_time"
			}`,
			wantErr: true,
		},
		{
			name: "invalid startDelta format",
			configJSON: `{
				"startDelta": "invalid_delta"
			}`,
			wantErr: true,
		},
		{
			name: "invalid JSON syntax",
			configJSON: `{
				"laps": 2,
				"lapLen="invalid";
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "config*.json")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.WriteString(tt.configJSON); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			got, err := loadConfig(tmpfile.Name())

			assert := assert.New(t)

			if strings.HasPrefix(tt.name, "valid") {
				assert.Nil(err)
				assert.Equal(tt.wantConfig.Laps, got.Laps)
				assert.Equal(tt.wantConfig.LapLen, got.LapLen)
				assert.Equal(tt.wantConfig.PenaltyLen, got.PenaltyLen)
				assert.Equal(tt.wantConfig.FiringLines, got.FiringLines)
				assert.Equal(tt.wantConfig.Start.Time, got.Start.Time)
				assert.Equal(tt.wantConfig.StartDelta.Duration, got.StartDelta.Duration)
			} else {
				assert.NotNil(err)
			}
		})
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := loadConfig("non_existent_file.json")
	assert.NotNil(t, err)
}

func parseTime(t *testing.T, format string, s string) time.Time {
	tt, err := time.Parse(format, s)
	if err != nil {
		t.Fatalf("Error parsing test time: %v", err)
	}
	return tt
}
