package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseEventLine(t *testing.T) {
	baseTime := time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC)

	tests := []struct {
		name    string
		input   string
		want    Event
		wantErr bool
	}{
		{
			name:  "valid basic event",
			input: "[09:30:00] 1 123",
			want: Event{
				Timestamp:    baseTime,
				ID:           EventRegistered,
				CompetitorID: 123,
				Extra:        nil,
			},
		},
		{
			name:  "valid event with extra params",
			input: "[09:30:00.500] 5 456 1",
			want: Event{
				Timestamp:    baseTime.Add(500 * time.Millisecond),
				ID:           EventStartedFiringRange,
				CompetitorID: 456,
				Extra:        []string{"1"},
			},
		},
		{
			name:  "valid event with multiple extras",
			input: "[09:30:00] 11 789 Lost in forest",
			want: Event{
				Timestamp:    baseTime,
				ID:           EventCantContinue,
				CompetitorID: 789,
				Extra:        []string{"Lost", "in", "forest"},
			},
		},
		{
			name:    "invalid empty line",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid timestamp format",
			input:   "[09:30] 1 123",
			wantErr: true,
		},
		{
			name:    "invalid event id",
			input:   "[09:30:00] invalid 123",
			wantErr: true,
		},
		{
			name:    "invalid competitor id",
			input:   "[09:30:00] 1 invalid",
			wantErr: true,
		},
		{
			name:    "insufficient parts",
			input:   "[09:30:00] 1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			got, err := parseEventLine(tt.input)

			if !tt.wantErr {
				assert.Nil(err)
				assert.Equal(tt.want.Timestamp, got.Timestamp)
				assert.Equal(tt.want.ID, got.ID)
				assert.Equal(tt.want.CompetitorID, got.CompetitorID)
				assert.Equal(len(tt.want.Extra), len(got.Extra))
				for i := range got.Extra {
					assert.Equal(tt.want.Extra[i], got.Extra[i])
				}
			} else {
				assert.NotNil((err))
			}
		})
	}
}

func TestParseEvents(t *testing.T) {
	input := `[09:30:00] 1 123
[09:30:01] 2 456 10:00:00
[bad line
[09:30:02] 3 789
`

	tests := []struct {
		name       string
		input      string
		wantEvents []Event
		wantPanic  bool
	}{
		{
			name:  "valid input with errors",
			input: input,
			wantEvents: []Event{
				{ID: EventRegistered, CompetitorID: 123},
				{ID: EventSetStartTime, CompetitorID: 456},
				{ID: EventOnStartLine, CompetitorID: 789},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			var logBuf bytes.Buffer
			in := strings.NewReader(tt.input)

			eventCh := parseEvents(in, &logBuf)

			var received []Event
			for e := range eventCh {
				received = append(received, e)
			}

			assert.Equal(len(tt.wantEvents), len(received))
			assert.Contains(logBuf.String(), "[ERROR] parseEventLine error has occured but was ignored")

			for i := range received {
				assert.Equal(tt.wantEvents[i].ID, received[i].ID)
				assert.Equal(tt.wantEvents[i].CompetitorID, received[i].CompetitorID)
			}
		})
	}
}

func TestParseScanner(t *testing.T) {
	errorReader := &errorReader{err: io.ErrUnexpectedEOF}

	var logBuf bytes.Buffer
	eventCh := parseEvents(errorReader, &logBuf)

	if _, ok := <-eventCh; ok {
		t.Fatal("Channel should be closed")
	}

	if !strings.Contains(logBuf.String(), "scanner") {
		t.Error("Missing scanner error log")
	}
}

type errorReader struct{ err error }

func (r *errorReader) Read(p []byte) (int, error) { return 0, r.err }
