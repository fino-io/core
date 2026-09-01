package logs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
)

type Service struct {
	ctx        context.Context
	logger     Logger
	callerSkip int
}

func NewService(logger Logger) *Service {
	return NewServiceWithCallerSkip(logger, 0)
}

func NewServiceWithCallerSkip(logger Logger, callerSkip int) *Service {
	if logger == nil {
		logger = newNopLogger()
	}
	return &Service{ctx: context.Background(), logger: logger, callerSkip: callerSkip}
}

// Ctx returns the default logging service bound to ctx.
func Ctx(ctx context.Context) *Service {
	return NewService(DefaultLogger()).WithContext(ctx)
}

// WithContext returns a copy of the service bound to ctx.
func (s *Service) WithContext(ctx context.Context) *Service {
	if ctx == nil {
		ctx = context.Background()
	}
	if s == nil {
		return &Service{ctx: ctx, logger: newNopLogger()}
	}
	clone := *s
	clone.ctx = ctx
	return &clone
}

func newDefaultService() *Service {
	return NewServiceWithCallerSkip(NewLoggerWith(NewDefaultConfig()), 1)
}

var (
	defaultServiceMu sync.RWMutex
	defaultService   = newDefaultService()
)

func DefaultLogger() Logger {
	return currentDefaultService().logger
}

func SetLogger(logger Logger) {
	defaultServiceMu.Lock()
	defaultService = NewServiceWithCallerSkip(logger, 1)
	defaultServiceMu.Unlock()
}

func SetLogLevel(level Level) {
	currentDefaultService().SetLogLevel(level)
}

func currentDefaultService() *Service {
	defaultServiceMu.RLock()
	defer defaultServiceMu.RUnlock()
	return defaultService
}

func (s *Service) Logger() Logger {
	if s == nil {
		return newNopLogger()
	}
	return s.logger
}

func (s *Service) SetLogger(logger Logger) {
	if s == nil {
		return
	}
	if logger == nil {
		logger = newNopLogger()
	}
	s.logger = logger
}

func (s *Service) SetLogLevel(level Level) {
	if s == nil || s.logger == nil {
		return
	}
	s.logger.SetLevel(level)
}

func (s *Service) log(level Level, msg string, fields ...Field) {
	if s == nil || s.logger == nil {
		return
	}
	ctx := s.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	s.logger.Log(ctx, Entry{
		Level:      level,
		Message:    msg,
		Fields:     mergeFields(fieldsFromContext(ctx), fields),
		CallerSkip: 1 + s.callerSkip + 2,
	})
}

func (s *Service) Debug(args ...interface{}) {
	s.log(DebugLevel, fmt.Sprint(args...))
}

func (s *Service) Info(args ...interface{}) {
	s.log(InfoLevel, fmt.Sprint(args...))
}

func (s *Service) Warn(args ...interface{}) {
	s.log(WarnLevel, fmt.Sprint(args...))
}

func (s *Service) Error(args ...interface{}) {
	s.log(ErrorLevel, fmt.Sprint(args...))
}

func (s *Service) Fatal(args ...interface{}) {
	s.log(FatalLevel, fmt.Sprint(args...))
}

func (s *Service) Debugf(template string, args ...interface{}) {
	s.log(DebugLevel, fmt.Sprintf(template, args...))
}

func (s *Service) Infof(template string, args ...interface{}) {
	s.log(InfoLevel, fmt.Sprintf(template, args...))
}

func (s *Service) Warnf(template string, args ...interface{}) {
	s.log(WarnLevel, fmt.Sprintf(template, args...))
}

func (s *Service) Errorf(template string, args ...interface{}) {
	s.log(ErrorLevel, fmt.Sprintf(template, args...))
}

func (s *Service) Fatalf(template string, args ...interface{}) {
	s.log(FatalLevel, fmt.Sprintf(template, args...))
}

func (s *Service) Debugw(msg string, keysAndValues ...interface{}) {
	s.log(DebugLevel, msg, keyValuesToFields(keysAndValues...)...)
}

func (s *Service) Infow(msg string, keysAndValues ...interface{}) {
	s.log(InfoLevel, msg, keyValuesToFields(keysAndValues...)...)
}

func (s *Service) Warnw(msg string, keysAndValues ...interface{}) {
	s.log(WarnLevel, msg, keyValuesToFields(keysAndValues...)...)
}

func (s *Service) Errorw(msg string, keysAndValues ...interface{}) {
	s.log(ErrorLevel, msg, keyValuesToFields(keysAndValues...)...)
}

func (s *Service) Fatalw(msg string, keysAndValues ...interface{}) {
	s.log(FatalLevel, msg, keyValuesToFields(keysAndValues...)...)
}

func (s *Service) NewError(args ...interface{}) error {
	msg := fmt.Sprint(args...)
	s.log(ErrorLevel, msg)
	return errors.New(msg)
}

func (s *Service) NewErrorf(template string, args ...interface{}) error {
	msg := fmt.Sprintf(template, args...)
	s.log(ErrorLevel, msg)
	return fmt.Errorf(template, args...)
}

func (s *Service) NewErrorw(msg string, keysAndValues ...interface{}) error {
	fields := keyValuesToFields(keysAndValues...)
	s.log(ErrorLevel, msg, fields...)
	return errors.New(renderErrorMessage(msg, fields))
}

func Debug(args ...interface{}) {
	currentDefaultService().Debug(args...)
}

func Info(args ...interface{}) {
	currentDefaultService().Info(args...)
}

func Warn(args ...interface{}) {
	currentDefaultService().Warn(args...)
}

func Error(args ...interface{}) {
	currentDefaultService().Error(args...)
}

func Fatal(args ...interface{}) {
	currentDefaultService().Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	currentDefaultService().Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	currentDefaultService().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	currentDefaultService().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	currentDefaultService().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	currentDefaultService().Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	currentDefaultService().Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	currentDefaultService().Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	currentDefaultService().Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	currentDefaultService().Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	currentDefaultService().Fatalw(msg, keysAndValues...)
}

func NewError(args ...interface{}) error {
	return currentDefaultService().NewError(args...)
}

func NewErrorf(template string, args ...interface{}) error {
	return currentDefaultService().NewErrorf(template, args...)
}

func NewErrorw(msg string, keysAndValues ...interface{}) error {
	return currentDefaultService().NewErrorw(msg, keysAndValues...)
}

func keyValuesToFields(keysAndValues ...interface{}) []Field {
	if len(keysAndValues) < 2 {
		return nil
	}

	fields := make([]Field, 0, len(keysAndValues)/2)
	for i := 0; i+1 < len(keysAndValues); i += 2 {
		fields = append(fields, Field{
			Key:   fmt.Sprint(keysAndValues[i]),
			Value: keysAndValues[i+1],
		})
	}
	return fields
}

func renderErrorMessage(msg string, fields []Field) string {
	if len(fields) == 0 {
		return msg
	}

	var b bytes.Buffer
	b.WriteString(msg)
	b.WriteByte(' ')
	for i, field := range fields {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("%v: %v", field.Key, field.Value))
	}
	return b.String()
}

type nopLogger struct{}

func newNopLogger() Logger { return nopLogger{} }

func (nopLogger) SetLevel(Level)             {}
func (nopLogger) GetLevel() Level            { return InfoLevel }
func (nopLogger) With(...Field) Logger       { return nopLogger{} }
func (nopLogger) Log(context.Context, Entry) {}
