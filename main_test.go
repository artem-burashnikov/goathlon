package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	assert.Panics(t, func() { must(loadConfig("")) })
}

func TestRunSingle(t *testing.T) {
	assert := assert.New(t)

	events, err := os.Open("examples/single/events")
	assert.Nil(err)
	defer events.Close()

	cfg, err := loadConfig("examples/single/config.json")
	assert.Nil(err)

	want, err := os.ReadFile("examples/single/output")
	assert.Nil(err)

	var out bytes.Buffer
	run(events, &out, cfg)

	assert.Equal(string(want), out.String())
}

func TestRunMultiple(t *testing.T) {
	assert := assert.New(t)

	events, err := os.Open("examples/multiple/events")
	assert.Nil(err)
	defer events.Close()

	cfg, err := loadConfig("examples/multiple/config.json")
	assert.Nil(err)

	want, err := os.ReadFile("examples/multiple/output")
	assert.Nil(err)

	var out bytes.Buffer
	run(events, &out, cfg)

	assert.Equal(string(want), out.String())
}
