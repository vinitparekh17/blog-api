package logger_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/jay-bhogayata/blogapi/logger"
	"github.com/stretchr/testify/assert"
)

var TestLog *slog.Logger

func TestInit(t *testing.T) {

	var buf bytes.Buffer

	logHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{AddSource: true})

	TestLog = slog.New(logHandler)

	slog.SetDefault(TestLog)

	logger.Init()

	if TestLog == nil {
		t.Errorf("Failed to initialize logger")
	}

	TestLog.Info("Test log message")
	TestLog.Error("Test Error message")

	assert.Contains(t, buf.String(), "Test log message", "log message not found in buffer")
	assert.Contains(t, buf.String(), "Test Error message", "error message not found in buffer")

	assert.Contains(t, buf.String(), "level=INFO", "error level not found in buffer")
	assert.Contains(t, buf.String(), "level=ERROR", "error level not found in buffer")

}
