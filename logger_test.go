package ctxlogger_test

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/timhugh/ctxlogger"
	"strings"
	"testing"
)

var levelStrings = map[ctxlogger.Level]string{
	ctxlogger.DebugLevel: "DEBUG",
	ctxlogger.InfoLevel:  "INFO",
	ctxlogger.WarnLevel:  "WARN",
	ctxlogger.ErrorLevel: "ERROR",
}

var logLevels = []ctxlogger.Level{
	ctxlogger.DebugLevel,
	ctxlogger.InfoLevel,
	ctxlogger.WarnLevel,
	ctxlogger.ErrorLevel,
}

const timestampRegex = `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`

func Test_BasicLogging(t *testing.T) {
	var out strings.Builder
	ctxlogger.SetOutput(&out)

	ctx := context.Background()

	for _, loggerLevel := range logLevels {
		for _, msgLevel := range logLevels {
			out.Reset()
			shouldLog := msgLevel >= loggerLevel
			loggerLevelName := levelStrings[loggerLevel]
			msgLevelName := levelStrings[msgLevel]

			var testName string
			if shouldLog {
				testName = "should log " + msgLevelName + " when logger level is " + loggerLevelName
			} else {
				testName = "should not log " + msgLevelName + " when logger level is " + loggerLevelName
			}
			t.Run(testName, func(t *testing.T) {
				out.Reset()
				ctxlogger.SetLevel(loggerLevel)

				logAtLevel(msgLevel, ctx, "test message")
				if shouldLog {
					assertJSON(t, map[string]interface{}{
						"level":   msgLevelName,
						"message": "test message",
					}, out.String())
				} else {
					assert.Empty(t, out.String())
				}
			})
		}
	}
}

func Test_LoggingWithFormat(t *testing.T) {
	var out strings.Builder
	ctxlogger.SetOutput(&out)
	ctxlogger.SetLevel(ctxlogger.DebugLevel)

	for _, msgLevel := range logLevels {
		msgLevelName := levelStrings[msgLevel]
		t.Run(msgLevelName+" should log with format", func(t *testing.T) {
			out.Reset()
			logAtLevel(msgLevel, context.Background(), "test message %s", "with format")
			assertJSON(t, map[string]interface{}{
				"level":   msgLevelName,
				"message": "test message with format",
			}, out.String())
		})
	}
}

func Test_LoggingContextParams(t *testing.T) {
	var out strings.Builder
	ctxlogger.SetOutput(&out)
	ctxlogger.SetLevel(ctxlogger.DebugLevel)

	for _, msgLevel := range logLevels {
		msgLevelName := levelStrings[msgLevel]
		t.Run(msgLevelName+" should log context params", func(t *testing.T) {
			out.Reset()
			ctx := context.Background()
			ctx = ctxlogger.AddParam(ctx, "key", "value")

			logAtLevel(msgLevel, ctx, "test message")

			assertJSON(t, map[string]interface{}{
				"level":   msgLevelName,
				"message": "test message",
				"key":     "value",
			}, out.String())
		})
	}
}

func assertJSON(t *testing.T, expected map[string]interface{}, actual string) {
	var actualValues map[string]interface{}
	err := json.Unmarshal([]byte(actual), &actualValues)
	assert.NoError(t, err)
	assert.EqualValues(t, expected, actualValues)
}

func logAtLevel(level ctxlogger.Level, ctx context.Context, msg string, params ...interface{}) {
	switch level {
	case ctxlogger.DebugLevel:
		ctxlogger.Debug(ctx, msg, params...)
	case ctxlogger.InfoLevel:
		ctxlogger.Info(ctx, msg, params...)
	case ctxlogger.WarnLevel:
		ctxlogger.Warn(ctx, msg, params...)
	case ctxlogger.ErrorLevel:
		ctxlogger.Error(ctx, msg, params...)
	}
}
