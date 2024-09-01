package ctxlogger

import (
    "fmt"
    "io"
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

type Logger struct {
    Level Level
    Out   io.Writer
}

func (l *Logger) Debug(msg string) {
    l.log(DebugLevel, msg)
}

func (l *Logger) Info(msg string) {
    l.log(InfoLevel, msg)
}

func (l *Logger) Warn(msg string) {
    l.log(WarnLevel, msg)
}

func (l *Logger) Error(msg string) {
    l.log(ErrorLevel, msg)
}

func (l *Logger) log(level Level, msg string) {
    if level >= l.Level {
        now := time.Now().UTC()
        _, err := l.Out.Write([]byte(now.Format(timestampFormat) + " [" + levelStrings[level] + "] " + msg))
        if err != nil {
            fmt.Printf("failed to write log message: %s\n", err.Error())
        }
    }
}

