package ctxlogger_test

import (
    "github.com/timhugh/ctxlogger"
    "github.com/stretchr/testify/assert"
    "testing"
    "strings"
    "regexp"
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

    logLevels := []ctxlogger.Level{
        ctxlogger.DebugLevel,
        ctxlogger.InfoLevel,
        ctxlogger.WarnLevel,
        ctxlogger.ErrorLevel,
    }

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
                log := ctxlogger.Logger{
                    Level: loggerLevel,
                    Out: &out,
                }
                out.Reset()
                switch msgLevel {
                case ctxlogger.DebugLevel:
                    log.Debug("test message")
                case ctxlogger.InfoLevel:
                    log.Info("test message")
                case ctxlogger.WarnLevel:
                    log.Warn("test message")
                case ctxlogger.ErrorLevel:
                    log.Error("test message")
                }
                if shouldLog {
                    assert.Regexp(t, regexp.MustCompile(timestampRegex + ` \[` + msgLevelName + `\] test message`), out.String())
                } else {
                    assert.Empty(t, out.String())
                }
            })
        }
    }
}
