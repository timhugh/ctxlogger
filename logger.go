package ctxlogger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
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

var level = InfoLevel
var out io.Writer = os.Stdout

type contextKeyType struct{}

var contextKey contextKeyType

func SetLevel(l Level) {
	level = l
}

func SetOutput(o io.Writer) {
	out = o
}

func Debug(ctx context.Context, msg string, params ...interface{}) {
	log(ctx, DebugLevel, msg, params...)
}

func Info(ctx context.Context, msg string, params ...interface{}) {
	log(ctx, InfoLevel, msg, params...)
}

func Warn(ctx context.Context, msg string, params ...interface{}) {
	log(ctx, WarnLevel, msg, params...)
}

func Error(ctx context.Context, msg string, params ...interface{}) {
	log(ctx, ErrorLevel, msg, params...)
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

func log(ctx context.Context, l Level, msgFormat string, msgParams ...interface{}) {
	if l < level {
		return
	}
	fullMsg := map[string]string{
		"message": fmt.Sprintf(msgFormat, msgParams...),
		"level":   levelStrings[l],
	}
	paramHash := ctx.Value(contextKey)
	logParams, ok := paramHash.(map[string]string)
	if ok {
		for k, v := range logParams {
			fullMsg[k] = v
		}
	}
	msg, err := json.Marshal(fullMsg)
	if err != nil {
		fmt.Printf("failed to marshal log message: %s\n", err.Error())
	}
	_, err = fmt.Fprintln(out, string(msg))
	if err != nil {
		fmt.Printf("failed to write log message: %s\n", err.Error())
	}
}
