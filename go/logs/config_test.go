package logs

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggerConfigurationValidation(t *testing.T) {
	for _, config := range []*Config{
		{Level: "typo"}, {Encode: "xml"}, {Output: "typo"},
		{File: FileConfig{Encode: "xml"}}, {File: FileConfig{MaxSize: -1}},
		{File: FileConfig{MaxBackups: -1}}, {File: FileConfig{MaxAge: -1}},
	} {
		logger, err := NewLoggerWith(config)
		require.Error(t, err)
		require.Nil(t, logger)
	}
	config := &Config{Encode: "json"}
	var output bytes.Buffer
	logger, err := NewLoggerWithWriter(config, &output)
	require.NoError(t, err)
	require.Equal(t, DebugLevel, logger.GetLevel())
	require.Empty(t, config.Level)
	logger.Log(context.Background(), Entry{Level: InfoLevel, Message: "custom output"})
	require.Contains(t, output.String(), `"message":"custom output"`)
	require.NoError(t, logger.Close())
	_, err = NewLoggerWithWriter(nil, nil)
	require.Error(t, err)
	var nilWriter *bytes.Buffer
	logger, err = NewLoggerWithWriter(nil, nilWriter)
	require.Error(t, err)
	require.Nil(t, logger)
}

func TestPublicWriterOwnership(t *testing.T) {
	sink := &borrowedSink{}
	logger, err := NewLoggerWithWriter(nil, sink)
	require.NoError(t, err)
	require.NoError(t, logger.Sync())
	require.NoError(t, logger.With(Field{Key: "child", Value: true}).Close())
	require.NoError(t, logger.Close())
	require.Zero(t, sink.closes)
}
