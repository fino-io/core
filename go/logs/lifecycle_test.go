package logs

import (
	"bytes"
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

type borrowedSink struct {
	bytes.Buffer
	syncs  int
	closes int
}

func (sink *borrowedSink) Sync() error  { sink.syncs++; return nil }
func (sink *borrowedSink) Close() error { sink.closes++; return nil }

func TestLoggerSinkOwnership(t *testing.T) {
	sink := &borrowedSink{}
	logger := newZapLogger(&Config{Level: "info", Encode: "json"}, sink)
	logger.Log(context.Background(), Entry{Level: InfoLevel, Message: "written"})
	require.NoError(t, logger.Sync())
	require.Equal(t, 1, sink.syncs)
	require.NoError(t, logger.Close())
	require.NoError(t, logger.With(Field{Key: "child", Value: true}).Close())
	require.Zero(t, sink.closes)
	require.Contains(t, sink.String(), "written")
}

func TestDerivedLoggerSharesCloseAndLevel(t *testing.T) {
	logger := newZapLogger(&Config{Level: "info"}, &bytes.Buffer{})
	var closes atomic.Int32
	logger.close = sync.OnceValue(func() error { closes.Add(1); return nil })
	child := logger.With(Field{Key: "child", Value: true})
	child.SetLevel(DebugLevel)
	require.Equal(t, DebugLevel, logger.GetLevel())
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _ = child.Close(); _ = logger.Close() }()
	}
	wg.Wait()
	require.Equal(t, int32(1), closes.Load())
}

func TestServiceWithLoggerIsImmutable(t *testing.T) {
	first := newZapLogger(&Config{Level: "info"}, &bytes.Buffer{})
	second := newZapLogger(&Config{Level: "debug"}, &bytes.Buffer{})
	parent := NewService(first)
	child := parent.WithLogger(second)
	require.Same(t, first, parent.Logger())
	require.Same(t, second, child.Logger())
	require.NotSame(t, parent, child)
	parent = NewService(nil)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				parent.WithLogger(nil).Info("child")
				parent.Info("parent")
			}
		}()
	}
	wg.Wait()
}

func TestGlobalLoggerConcurrentReplacement(t *testing.T) {
	previous := DefaultLogger()
	t.Cleanup(func() { SetLogger(previous) })
	SetLogger(newNopLogger())
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			SetLogger(newNopLogger())
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			Ctx(context.Background()).Info("audit")
			SetLogLevel(InfoLevel)
		}
	}()
	wg.Wait()
}
