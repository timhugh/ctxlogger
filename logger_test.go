package ctxlogger_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/timhugh/ctxlogger"
	"regexp"
	"strings"
	"testing"
)

var levelStrings = map[ctxlogger.Level]string{
	ctxlogger.DebugLevel: "DEBUG",
	ctxlogger.InfoLevel:  "INFO",
	ctxlogger.WarnLevel:  "WARN",
	ctxlogger.ErrorLevel: "ERROR",
}

const timestampRegex = `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`

func Test_BasicLogging(t *testing.T) {
	var out strings.Builder
	ctxlogger.SetOutput(&out)

	logLevels := []ctxlogger.Level{
		ctxlogger.DebugLevel,
		ctxlogger.InfoLevel,
		ctxlogger.WarnLevel,
		ctxlogger.ErrorLevel,
	}

	ctx := context.Background()

	for _, loggerLevel := range logLevels {
		for _, msgLevel := range logLevels {
			out.Reset()
			shouldLog := msgLevel >= loggerLevel
			loggerLevelName := levelStrings[loggerLevel]
			msgLevelName := levelStrings[msgLevel]

			var testName string
			if shouldLog {
				testName = "should log " + loggerLevelName + " when logger level is " + msgLevelName
			} else {
				testName = "should not log " + loggerLevelName + " when logger level is " + msgLevelName
			}
			t.Run(testName, func(t *testing.T) {
				out.Reset()
				ctxlogger.SetLevel(loggerLevel)

				switch msgLevel {
				case ctxlogger.DebugLevel:
					ctxlogger.Debug(ctx, "test message")
				case ctxlogger.InfoLevel:
					ctxlogger.Info(ctx, "test message")
				case ctxlogger.WarnLevel:
					ctxlogger.Warn(ctx, "test message")
				case ctxlogger.ErrorLevel:
					ctxlogger.Error(ctx, "test message")
				}
				if shouldLog {
					assert.Regexp(t, regexp.MustCompile(timestampRegex+` \[`+msgLevelName+`\] test message`), out.String())
				} else {
					assert.Empty(t, out.String())
				}
			})
		}
	}
}

func Test_LoggingContextParams(t *testing.T) {
	var out strings.Builder
	ctxlogger.SetOutput(&out)
	ctxlogger.SetLevel(ctxlogger.InfoLevel)

	t.Run("should log context params", func(t *testing.T) {
		ctx := context.Background()
		ctx = ctxlogger.AddParam(ctx, "key", "value")

		ctxlogger.Info(ctx, "test message")

		assert.Regexp(t, regexp.MustCompile(timestampRegex+` \[INFO\] test message key=value`), out.String())
	})
}
