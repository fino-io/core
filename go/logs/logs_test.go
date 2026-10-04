package logs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggerWithFieldsOverrideOnce(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newZapLogger(&Config{Level: "info", Encode: "json", InitFields: map[string]any{"source": "initial"}}, buf)
	child := logger.With(Field{Key: "source", Value: "base"}, Field{Key: "id", Value: 7})
	child = child.With(Field{Key: "source", Value: "child"})
	NewService(child).Infow("ready", "source", "entry")
	require.Equal(t, 1, strings.Count(buf.String(), `"source":`))
	require.Equal(t, 1, strings.Count(buf.String(), `"id":`))
	require.Contains(t, buf.String(), `"source":"entry"`)
	buf.Reset()
	NewService(logger).Info("parent")
	require.Contains(t, buf.String(), `"source":"initial"`)
	require.NotContains(t, buf.String(), `"id":`)
}

func TestLogCaller(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newZapLogger(&Config{Level: "info", Encode: "json"}, buf)
	svc := NewService(logger)
	previous := DefaultLogger()
	SetLogger(logger)
	t.Cleanup(func() { SetLogger(previous) })
	for _, write := range []func(){func() { svc.Info("caller") }, func() { Info("caller") }} {
		buf.Reset()
		write()
		var entry map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
		require.Contains(t, entry["caller"], "logs_test.go:")
		require.NotContains(t, entry["caller"], "logs.go:")
	}
	_, file, line, _ := runtime.Caller(0)
	svc.Info("exact caller")
	entries := strings.Split(strings.TrimSpace(buf.String()), "\n")
	var entry map[string]any
	require.NoError(t, json.Unmarshal([]byte(entries[len(entries)-1]), &entry))
	require.Equal(t, filepath.Join(filepath.Base(filepath.Dir(file)), filepath.Base(file))+":"+fmt.Sprint(line+1), entry["caller"])
}

func TestNewErrorf(t *testing.T) {
	err := NewErrorf("this is %s", "NewErrorf")
	require.Error(t, err)
}

func TestNewErrorfLogsWrappedError(t *testing.T) {
	buf := &bytes.Buffer{}
	svc := NewService(newZapLogger(&Config{Level: "info", Encode: "json"}, buf))
	cause := errors.New("connection failed")

	err := svc.NewErrorf("lookup: %w", cause)

	require.ErrorIs(t, err, cause)
	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	require.Equal(t, err.Error(), entry["message"])
}

func TestLoggerLevelAndFields(t *testing.T) {
	buf := &bytes.Buffer{}
	svc := NewService(newZapLogger(&Config{
		Level:  "warn",
		Encode: "json",
		Output: "console",
	}, buf))

	svc.Debug("debug")
	svc.Info("info")
	svc.Warnw("warn", "user", "alice")
	svc.Errorw("error", "count", 2)

	out := buf.String()
	require.NotContains(t, out, `"message":"debug"`)
	require.NotContains(t, out, `"message":"info"`)
	require.Contains(t, out, `"message":"warn"`)
	require.Contains(t, out, `"user":"alice"`)
	require.Contains(t, out, `"message":"error"`)
	require.Contains(t, out, `"count":2`)
}

func TestServiceHelpers(t *testing.T) {
	buf := &bytes.Buffer{}
	svc := NewService(newZapLogger(&Config{
		Level:  "info",
		Encode: "json",
		Output: "console",
	}, buf))

	require.Equal(t, InfoLevel, svc.Logger().GetLevel())
	svc.SetLogLevel(DebugLevel)
	require.Equal(t, DebugLevel, svc.Logger().GetLevel())

	err := svc.NewErrorw("failed", "kind", "network")
	require.EqualError(t, err, "failed kind: network")

	svc.Debugf("value=%s", "x")
	svc.Infow("ready", "id", 1)

	out := buf.String()
	require.Contains(t, out, `"message":"failed"`)
	require.Contains(t, out, `"kind":"network"`)
	require.Contains(t, out, `"message":"value=x"`)
	require.Contains(t, out, `"message":"ready"`)
}

func TestPackageDefaultLoggerSwap(t *testing.T) {
	prev := DefaultLogger()
	t.Cleanup(func() {
		SetLogger(prev)
	})

	buf := &bytes.Buffer{}
	SetLogger(newZapLogger(&Config{
		Level:  "info",
		Encode: "json",
		Output: "console",
	}, buf))

	Infow("hello", "id", 7)
	require.Contains(t, buf.String(), `"message":"hello"`)
	require.Contains(t, buf.String(), `"id":7`)
}

func TestLevelHandlerServeHTTP(t *testing.T) {
	logger, err := NewLoggerWith(&Config{Level: "info"})
	require.NoError(t, err)
	handler := logger.LevelHandler()

	req := httptest.NewRequest(http.MethodGet, "/level", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"level":"info"`)

	req = httptest.NewRequest(http.MethodPut, "/level", strings.NewReader(`{"level":"debug"}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, DebugLevel, logger.GetLevel())
	require.Contains(t, rec.Body.String(), `"level":"debug"`)
}

func TestFileOutput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")
	logger, err := NewLoggerWith(&Config{
		Level:  "info",
		Encode: "json",
		Output: "file",
		File: FileConfig{
			Path:   path,
			Encode: "json",
			// Encode:     "console",
			MaxSize:    1,
			MaxBackups: 1,
			MaxAge:     1,
		},
	})
	require.NoError(t, err)

	logger.Log(context.Background(), Entry{
		Level:   InfoLevel,
		Message: "persisted",
		Fields:  []Field{{Key: "user", Value: "bob"}},
	})
	require.NoError(t, logger.Sync())
	require.NoError(t, logger.Close())
	require.NoError(t, logger.Close())

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.True(t, strings.Contains(string(data), `"message":"persisted"`))
	require.True(t, strings.Contains(string(data), `"user":"bob"`))
}

func TestContextFields(t *testing.T) {
	buf := &bytes.Buffer{}
	svc := NewService(newZapLogger(&Config{
		Level:  "info",
		Encode: "json",
		Output: "console",
	}, buf))

	ctx := WithFields(context.Background(),
		Field{Key: "request_id", Value: "request-1"},
		Field{Key: "source", Value: "context"},
	)
	svc.WithContext(ctx).Infow("handled", "source", "entry")

	out := buf.String()
	require.Contains(t, out, `"request_id":"request-1"`)
	require.Contains(t, out, `"source":"entry"`)
	require.NotContains(t, out, `"source":"context"`)
}

func TestWithFieldsDoesNotMutateParent(t *testing.T) {
	parent := WithFields(context.Background(), Field{Key: "request_id", Value: "parent"})
	child := WithFields(parent, Field{Key: "request_id", Value: "child"})

	require.Equal(t, "parent", fieldsFromContext(parent)[0].Value)
	require.Equal(t, "child", fieldsFromContext(child)[0].Value)
}

func TestContextReachesLogger(t *testing.T) {
	type key struct{}
	ctx := context.WithValue(context.Background(), key{}, "value")
	logger := &contextRecordingLogger{}

	NewService(logger).WithContext(ctx).Info("handled")

	require.Equal(t, "value", logger.ctx.Value(key{}))
}

type contextRecordingLogger struct {
	nopLogger
	ctx context.Context
}

func (l *contextRecordingLogger) Log(ctx context.Context, _ Entry) {
	l.ctx = ctx
}
