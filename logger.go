package ctxlogger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var levelStrings = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
}

const timestampFormat = "2006-01-02T15:04:05Z"

var level Level = InfoLevel
var out io.Writer = os.Stdout

type contextKeyType struct{}

var contextKey contextKeyType

func SetLevel(l Level) {
	level = l
}

func SetOutput(o io.Writer) {
	out = o
}

func Debug(ctx context.Context, msg string) {
	log(ctx, DebugLevel, msg)
}

func Info(ctx context.Context, msg string) {
	log(ctx, InfoLevel, msg)
}

func Warn(ctx context.Context, msg string) {
	log(ctx, WarnLevel, msg)
}

func Error(ctx context.Context, msg string) {
	log(ctx, ErrorLevel, msg)
}

func AddParam(ctx context.Context, key, value string) context.Context {
	paramHash := ctx.Value(contextKey)
	params, ok := paramHash.(map[string]string)
	if !ok {
		params = make(map[string]string)
	}
	params[key] = value
	return context.WithValue(ctx, contextKey, params)
}

func log(ctx context.Context, l Level, msg string) {
	if l < level {
		return
	}
	_, err := fmt.Fprintf(out, "%s [%s] %s %s\n", time.Now().UTC().Format(timestampFormat), levelStrings[l], msg, stringifyParams(ctx))
	if err != nil {
		fmt.Printf("failed to write log message: %s\n", err.Error())
	}
}

func stringifyParams(ctx context.Context) string {
	paramHash := ctx.Value(contextKey)
	params, ok := paramHash.(map[string]string)
	if !ok {
		return ""
	}
	var s string
	for k, v := range params {
		s += k + "=" + v + " "
	}
	return s
}
