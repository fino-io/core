package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp_Format(t *testing.T) {
	now := Now()
	s := now.String()
	f := now.Format()
	fmt.Println("s: ", s)
	fmt.Println("f: ", f)

	ts := &Timestamp{Seconds: 1759754416, Nanoseconds: 995717000}
	str := ts.Format()

	assert.NotEmpty(t, str)
	assert.Equal(t, "2025-10-06T20:40:16.995+08:00", str)
}
