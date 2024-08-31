package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"strings"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type contextKey int

const loggerContextKey contextKey = 1

type ContextLogger struct {
	level  Level
	params map[string]interface{}
	out    io.Writer
}

func getLevelFromEnv() Level {
	level := os.Getenv("LOG_LEVEL")
	return stringToLevel(level)
}

func stringToLevel(level string) Level {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

func levelToString(level Level) string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "info"
	}
}

func New(out io.Writer) *ContextLogger {
	level := getLevelFromEnv()
	return &ContextLogger{
		level: level,
		out:   out,
	}
}

func Default() *ContextLogger {
	return New(os.Stdout)
}

func NewContext(out io.Writer) context.Context {
	logger := New(out)
	return context.WithValue(context.Background(), loggerContextKey, logger)
}

func Attach(ctx context.Context, logger *ContextLogger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

func FromContext(ctx context.Context) (context.Context, *ContextLogger) {
	var logger *ContextLogger
	ctxValue := ctx.Value(loggerContextKey)
	if ctxValue == nil {
		logger = Default()
		logger.Debug("No logger found in context; attaching a default logger")
		ctx = Attach(ctx, logger)
	} else {
		logger = ctxValue.(*ContextLogger)
	}
	return ctx, logger
}

func (l *ContextLogger) NewContext() context.Context {
	ctx := context.Background()
	ctx = Attach(ctx, l)
	return ctx
}

func (l *ContextLogger) Sub() *ContextLogger {
	return &ContextLogger{
		params: maps.Clone(l.params),
		out:    l.out,
		level:  l.level,
	}
}

func (l *ContextLogger) AddParam(key string, value interface{}) *ContextLogger {
	if l.params == nil {
		l.params = make(map[string]interface{})
	}
	l.params[key] = value
	return l
}

func (l *ContextLogger) AddParams(params map[string]interface{}) *ContextLogger {
	if l.params == nil {
		l.params = make(map[string]interface{})
	}
	for key, value := range params {
		l.params[key] = value
	}
	return l
}

func (l *ContextLogger) Debug(msg string, params ...interface{}) {
	l.log(DebugLevel, msg, params...)
}

func (l *ContextLogger) Info(msg string, params ...interface{}) {
	l.log(InfoLevel, msg, params...)
}

func (l *ContextLogger) Warn(msg string, params ...interface{}) {
	l.log(WarnLevel, msg, params...)
}

func (l *ContextLogger) Error(msg string, params ...interface{}) {
	l.log(ErrorLevel, msg, params...)
}

func (l *ContextLogger) Fatal(msg string, params ...interface{}) {
	l.log(FatalLevel, msg, params...)
}

func (l *ContextLogger) log(level Level, msg string, params ...interface{}) {
	if level < l.level {
		return
	}

	msgJson := make(map[string]interface{})
	msgJson["level"] = levelToString(level)
	msgJson["message"] = fmt.Sprintf(msg, params...)
	for k, v := range l.params {
		msgJson[k] = v
	}

	jsonString, err := json.Marshal(msgJson)
	if err != nil {
		_, err := fmt.Fprintf(l.out, "Failed to marshal log message: %s", err)
		if err != nil {
			fmt.Printf("Failed to write log message: %s\n", err)
		}
		return
	}

	_, err = fmt.Fprintln(l.out, string(jsonString))
	if err != nil {
		fmt.Printf("Failed to write log message: %s\n", err)
	}
}
