package logger_test

import (
	"context"
	"encoding/json"
    "github.com/stretchr/testify/assert"
	"github.com/timhugh/digitalvenue/util/logger"
	"os"
	"strings"
	"testing"
)

func initTest(t *testing.T) (*assert.Assertions, strings.Builder, map[string]interface{}) {
    assert := assert.New(t)
	out := strings.Builder{}
	msg := make(map[string]interface{})

	err := os.Setenv("LOG_LEVEL", "debug")
	if err != nil {
		t.Fatalf("Error setting environment variable: %s", err)
	}

	return assert, out, msg
}

func TestContextLogger_BasicLogging(t *testing.T) {
	assert, out, msg := initTest(t)

	log := logger.New(&out)

	log.Debug("This is a %s level message", "debug")
	assert.NoError(json.Unmarshal([]byte(out.String()), &msg))
	assert.Equal(msg["level"], "debug")
	assert.Equal(msg["message"], "This is a debug level message")

	clearOutputs(&out, &msg)

	log.Info("This is an %s level message", "info")
	assert.NoError(json.Unmarshal([]byte(out.String()), &msg))
	assert.Equal(msg["level"], "info")
	assert.Equal(msg["message"], "This is an info level message")

	clearOutputs(&out, &msg)

	log.Warn("This is a %s level message", "warn")
	assert.NoError(json.Unmarshal([]byte(out.String()), &msg))
	assert.Equal(msg["level"], "warn")
	assert.Equal(msg["message"], "This is a warn level message")

	clearOutputs(&out, &msg)

	log.Error("This is an %s level message", "error")
	assert.NoError(json.Unmarshal([]byte(out.String()), &msg))
	assert.Equal(msg["level"], "error")
	assert.Equal(msg["message"], "This is an error level message")

	clearOutputs(&out, &msg)

	log.Fatal("This is a %s level message", "fatal")
	assert.NoError(json.Unmarshal([]byte(out.String()), &msg))
	assert.Equal(msg["level"], "fatal")
	assert.Equal(msg["message"], "This is a fatal level message")
}

func TestContextLogger_Params(t *testing.T) {
	assert, out, msg := initTest(t)

	log := logger.New(&out)

	log.AddParam("key", "value")
	log.AddParams(map[string]interface{}{"key2": 123.0, "key3": true})

	log.Debug("Message")
	assert.NoError(json.Unmarshal([]byte(out.String()), &msg))
	assert.Equal(msg["level"], "debug")
	assert.Equal(msg["message"], "Message")
	assert.Equal(msg["key"], "value")
	assert.Equal(msg["key2"], 123.0)
	assert.Equal(msg["key3"], true)
}

func TestContextLogger_NewContext(t *testing.T) {
	assert, out, msg := initTest(t)

	ctx := logger.NewContext(&out)
	_, log := logger.FromContext(ctx)

	log.Debug("Message")
	assert.NoError(json.Unmarshal([]byte(out.String()), &msg))
	assert.Equal(msg["level"], "debug")
	assert.Equal(msg["message"], "Message")
}

func TestContextLogger_ExistingContext(t *testing.T) {
	assert, out, msg := initTest(t)

	logIn := logger.New(&out)
	logIn.AddParam("key", "value")

	ctx := logger.Attach(context.Background(), logIn)

	_, logOut := logger.FromContext(ctx)
	assert.Equal(logIn, logOut)

	logOut.Debug("Message")
	assert.NoError(json.Unmarshal([]byte(out.String()), &msg))
	assert.Equal(msg["key"], "value")
}

func TestContextLogger_SubLogger(t *testing.T) {
    assert, out, msg := initTest(t)

	log := logger.New(&out)
	log.AddParam("key", "value")

	sub := log.Sub()
	sub.AddParam("key2", "value2")

	sub.Debug("Message")
	assert.NoError(json.Unmarshal([]byte(out.String()), &msg))
	assert.Equal(msg["key"], "value")
	assert.Equal(msg["key2"], "value2")

	clearOutputs(&out, &msg)

	log.Debug("Message")
	assert.NoError(json.Unmarshal([]byte(out.String()), &msg))
	assert.Equal(msg["key"], "value")
	assert.Equal(msg["key2"], nil)
}

func TestContextLogger_Chaining(t *testing.T) {
	assert, out, msg := initTest(t)

	log := logger.New(&out).AddParam("key", "value").AddParams(map[string]interface{}{"key2": 123.0, "key3": true})

	log.Debug("Message")
	assert.NoError(json.Unmarshal([]byte(out.String()), &msg))
	assert.Equal(msg["level"], "debug")
	assert.Equal(msg["message"], "Message")
	assert.Equal(msg["key"], "value")
	assert.Equal(msg["key2"], 123.0)
	assert.Equal(msg["key3"], true)
}

func clearOutputs(out *strings.Builder, msg *map[string]interface{}) {
	out.Reset()
	*msg = make(map[string]interface{})
}
