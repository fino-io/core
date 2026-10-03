package logs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	SetLevel(Level)
	GetLevel() Level
	With(...Field) Logger
	Log(context.Context, Entry)
}

type Field struct {
	Key   string
	Value any
}

type Entry struct {
	Level      Level
	Message    string
	Fields     []Field
	CallerSkip int
}

type zapLogger struct {
	level  zap.AtomicLevel
	logger *zap.Logger
	fields []Field
}

func NewLoggerWith(cfg *Config) Logger {
	if cfg == nil {
		cfg = NewDefaultConfig()
	}
	return newZapLogger(cfg, nil)
}

func newZapLogger(cfg *Config, sink io.Writer) *zapLogger {
	if cfg == nil {
		cfg = NewDefaultConfig()
	}

	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level.SetLevel(zapcore.InfoLevel)
	}

	encode := cfg.Encode
	if strings.EqualFold(cfg.Output, "file") {
		encode = cfg.File.Encode
		if encode == "" {
			encode = "json"
		}
	}

	encoder := newEncoder(encode)
	writeSyncer := newWriteSyncer(cfg, sink)
	core := zapcore.NewCore(encoder, writeSyncer, level)
	base := zap.New(core, zap.AddCaller())

	keys := make([]string, 0, len(cfg.InitFields))
	for key := range cfg.InitFields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	fields := make([]Field, 0, len(keys))
	for _, key := range keys {
		fields = append(fields, Field{Key: key, Value: cfg.InitFields[key]})
	}

	if len(cfg.LevelPattern) > 0 && cfg.LevelPort > 0 {
		if _, err := startLevelServer(&level, cfg.LevelPattern, cfg.LevelPort); err != nil {
			log.Printf("failed to start log level server: %s\n", err)
		}
	}

	return &zapLogger{
		level:  level,
		logger: base,
		fields: fields,
	}
}

func newEncoder(encode string) zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:             "time",
		LevelKey:            "level",
		NameKey:             "logger",
		CallerKey:           "caller",
		MessageKey:          "message",
		StacktraceKey:       "stacktrace",
		LineEnding:          zapcore.DefaultLineEnding,
		EncodeLevel:         zapcore.LowercaseLevelEncoder,
		EncodeTime:          zapcore.ISO8601TimeEncoder,
		EncodeDuration:      zapcore.SecondsDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		NewReflectedEncoder: jsoniterReflectedEncoder,
	}

	if strings.EqualFold(encode, "json") {
		return zapcore.NewJSONEncoder(cfg)
	}
	return zapcore.NewConsoleEncoder(cfg)
}

func newWriteSyncer(cfg *Config, sink io.Writer) zapcore.WriteSyncer {
	if sink != nil {
		return zapcore.Lock(zapcore.AddSync(sink))
	}
	if strings.EqualFold(cfg.Output, "file") {
		fileCfg := cfg.File
		if fileCfg.Path == "" {
			fileCfg.Path = NewDefaultConfig().File.Path
		}
		return zapcore.Lock(zapcore.AddSync(&lumberjack.Logger{
			Filename:   fileCfg.Path,
			MaxSize:    fileCfg.MaxSize,
			MaxBackups: fileCfg.MaxBackups,
			MaxAge:     fileCfg.MaxAge,
			Compress:   fileCfg.Compress,
		}))
	}
	return zapcore.Lock(zapcore.AddSync(os.Stdout))
}

func levelHandler(level *zap.AtomicLevel, pattern string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(pattern, level)
	return mux
}

func startLevelServer(level *zap.AtomicLevel, pattern string, port int) (*http.Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	return startLevelServerWithListener(level, pattern, ln), nil
}

func startLevelServerWithListener(level *zap.AtomicLevel, pattern string, ln net.Listener) *http.Server {
	svc := &http.Server{
		Addr:         ln.Addr().String(),
		Handler:      levelHandler(level, pattern),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		addr := ln.Addr().String()
		fmt.Printf(
			"level serve on addr:%s\nusage: [GET] curl http://%s%s\nusage: [PUT] curl -XPUT --data '{\"level\":\"debug\"}' http://%s%s\n",
			addr, addr, pattern, addr, pattern,
		)
		if err := svc.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to serve log level: %s\n", err)
		}
	}()

	return svc
}

func jsoniterReflectedEncoder(w io.Writer) zapcore.ReflectedEncoder {
	enc := jsoniter.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc
}

func (l *zapLogger) SetLevel(level Level) {
	l.level.SetLevel(zapcore.Level(level))
}

func (l *zapLogger) GetLevel() Level {
	return Level(l.level.Level())
}

func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{
		level:  l.level,
		logger: l.logger,
		fields: mergeFields(l.fields, fields),
	}
}

func (l *zapLogger) Log(_ context.Context, entry Entry) {
	if l == nil || l.logger == nil {
		return
	}

	callerSkip := entry.CallerSkip
	if callerSkip <= 0 {
		callerSkip = 1
	}
	level := entry.Level
	if level < DebugLevel || level > FatalLevel {
		level = InfoLevel
	}
	fields := mergeFields(l.fields, entry.Fields)
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	l.logger.WithOptions(zap.AddCallerSkip(callerSkip)).Log(zapcore.Level(level), entry.Message, zapFields...)
}
