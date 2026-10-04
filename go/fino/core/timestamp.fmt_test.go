package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp_Format(t *testing.T) {
	ts := &Timestamp{Seconds: 1759754416, Nanoseconds: 995717000}
	str := ts.Format()

	assert.NotEmpty(t, str)
	assert.Equal(t, "2025-10-06T12:40:16.995717Z", str)
}
