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
	sugar  *zap.SugaredLogger
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
	if strings.EqualFold(cfg.Output, "file") && cfg.File.Encode != "" {
		encode = cfg.File.Encode
	}

	encoder := newEncoder(encode)
	writeSyncer := newWriteSyncer(cfg, sink)
	core := zapcore.NewCore(encoder, writeSyncer, level)
	base := zap.New(core, zap.AddCaller())

	if len(cfg.InitFields) > 0 {
		fields := make([]zap.Field, 0, len(cfg.InitFields))
		for key, value := range cfg.InitFields {
			fields = append(fields, zap.Any(key, value))
		}
		base = base.With(fields...)
	}

	if len(cfg.LevelPattern) > 0 && cfg.LevelPort > 0 {
		if _, err := startLevelServer(&level, cfg.LevelPattern, cfg.LevelPort); err != nil {
			log.Printf("failed to start log level server: %s\n", err)
		}
	}

	return &zapLogger{
		level: level,
		sugar: base.Sugar(),
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
		return zapcore.AddSync(sink)
	}
	if strings.EqualFold(cfg.Output, "file") {
		fileCfg := cfg.File
		if fileCfg.Path == "" {
			fileCfg.Path = NewDefaultConfig().File.Path
		}
		if fileCfg.Encode == "" {
			fileCfg.Encode = "json"
		}
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   fileCfg.Path,
			MaxSize:    fileCfg.MaxSize,
			MaxBackups: fileCfg.MaxBackups,
			MaxAge:     fileCfg.MaxAge,
			Compress:   fileCfg.Compress,
		})
	}
	return zapcore.AddSync(os.Stdout)
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

func (l *zapLogger) clone() *zapLogger {
	fields := make([]Field, len(l.fields))
	copy(fields, l.fields)
	return &zapLogger{
		level:  l.level,
		sugar:  l.sugar,
		fields: fields,
	}
}

func (l *zapLogger) SetLevel(level Level) {
	l.level.SetLevel(zapcore.Level(level))
}

func (l *zapLogger) GetLevel() Level {
	return Level(l.level.Level())
}

func (l *zapLogger) With(fields ...Field) Logger {
	if len(fields) == 0 {
		return l.clone()
	}

	kv := make([]any, 0, len(fields)*2)
	for _, field := range fields {
		kv = append(kv, field.Key, field.Value)
	}

	return &zapLogger{
		level:  l.level,
		sugar:  l.sugar.With(kv...),
		fields: append(append([]Field{}, l.fields...), fields...),
	}
}

func (l *zapLogger) Log(_ context.Context, entry Entry) {
	if l == nil || l.sugar == nil {
		return
	}

	callerSkip := entry.CallerSkip
	if callerSkip <= 0 {
		callerSkip = 1
	}

	sugar := l.sugar.WithOptions(zap.AddCallerSkip(callerSkip))
	kv := fieldsToKeyValues(l.fields, entry.Fields)
	switch entry.Level {
	case DebugLevel:
		if len(kv) == 0 {
			sugar.Debug(entry.Message)
			return
		}
		sugar.Debugw(entry.Message, kv...)
	case InfoLevel:
		if len(kv) == 0 {
			sugar.Info(entry.Message)
			return
		}
		sugar.Infow(entry.Message, kv...)
	case WarnLevel:
		if len(kv) == 0 {
			sugar.Warn(entry.Message)
			return
		}
		sugar.Warnw(entry.Message, kv...)
	case ErrorLevel:
		if len(kv) == 0 {
			sugar.Error(entry.Message)
			return
		}
		sugar.Errorw(entry.Message, kv...)
	case DPanicLevel:
		if len(kv) == 0 {
			sugar.DPanic(entry.Message)
			return
		}
		sugar.DPanicw(entry.Message, kv...)
	case PanicLevel:
		if len(kv) == 0 {
			sugar.Panic(entry.Message)
			return
		}
		sugar.Panicw(entry.Message, kv...)
	case FatalLevel:
		if len(kv) == 0 {
			sugar.Fatal(entry.Message)
			return
		}
		sugar.Fatalw(entry.Message, kv...)
	default:
		if len(kv) == 0 {
			sugar.Info(entry.Message)
			return
		}
		sugar.Infow(entry.Message, kv...)
	}
}

func fieldsToKeyValues(base []Field, extra []Field) []any {
	total := len(base) + len(extra)
	if total == 0 {
		return nil
	}

	kv := make([]any, 0, total*2)
	for _, field := range base {
		kv = append(kv, field.Key, field.Value)
	}
	for _, field := range extra {
		kv = append(kv, field.Key, field.Value)
	}
	return kv
}
