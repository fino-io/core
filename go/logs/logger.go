package logs

import (
	"context"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

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
	LevelHandler() http.Handler
	Sync() error
	Close() error
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
	close  func() error
}

func NewLoggerWith(cfg *Config) Logger {
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
	writeSyncer, closeSink := newWriteSyncer(cfg, sink)
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

	return &zapLogger{
		level:  level,
		logger: base,
		fields: fields,
		close:  closeSink,
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

func newWriteSyncer(cfg *Config, sink io.Writer) (zapcore.WriteSyncer, func() error) {
	if sink != nil {
		return zapcore.Lock(zapcore.AddSync(sink)), func() error { return nil }
	}
	if strings.EqualFold(cfg.Output, "file") {
		fileCfg := cfg.File
		if fileCfg.Path == "" {
			fileCfg.Path = NewDefaultConfig().File.Path
		}
		file := &lumberjack.Logger{
			Filename:   fileCfg.Path,
			MaxSize:    fileCfg.MaxSize,
			MaxBackups: fileCfg.MaxBackups,
			MaxAge:     fileCfg.MaxAge,
			Compress:   fileCfg.Compress,
		}
		return zapcore.Lock(zapcore.AddSync(file)), sync.OnceValue(file.Close)
	}
	return zapcore.Lock(zapcore.AddSync(os.Stdout)), func() error { return nil }
}

// LevelHandler can be mounted on the application's HTTP server.
func (l *zapLogger) LevelHandler() http.Handler {
	return l.level
}

func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

// Close releases the owned file sink. Derived loggers share its lifetime.
// Console output and writers supplied by callers are not owned by the logger.
func (l *zapLogger) Close() error {
	return l.close()
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
	clone := *l
	clone.fields = mergeFields(l.fields, fields)
	return &clone
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
