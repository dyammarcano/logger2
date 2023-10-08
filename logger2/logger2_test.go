package logger2

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	dir := "testlogs"
	logger, err := NewLoggerDefault(dir)

	assert.Nil(t, err, "Failed to create a Logger: %v", err)
	assert.NotNil(t, logger, "Logger is nil")

	// Test logging a message
	logger.Info("Test info log message")

	// Clean up
	if err := os.RemoveAll(logger.LogDir); err != nil {
		t.Errorf("Failed to remove test logs: %v", err)
	}
}

func TestLog(t *testing.T) {
	dir := "testlogs"
	logger, err := NewLoggerDefault(dir)

	assert.Nil(t, err, "Failed to create a Logger: %v", err)
	assert.NotNil(t, logger, "Logger is nil")

	// Test logging with each log level
	logger.Log(DebugLevel, "Testing log at level: Debug")
	logger.Log(InfoLevel, "Testing log at level: Info")
	logger.Log(WarnLevel, "Testing log at level: Warn")
	logger.Log(ErrorLevel, "Testing log at level: Error")
	logger.Log(FatalLevel, "Testing log at level: Fatal")

	// Clean up
	if err := os.RemoveAll(logger.LogDir); err != nil {
		t.Errorf("Failed to remove test logs: %v", err)
	}
}
