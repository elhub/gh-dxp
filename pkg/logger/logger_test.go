package logger_test

import (
	"errors"
	"testing"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/logger"
)

func TestDebug(_ *testing.T) {
	log.SetLevel(log.DebugLevel)
	// Just verify it doesn't panic
	logger.Debug("debug message")
}

func TestInfo(_ *testing.T) {
	logger.Info("info message")
}

func TestWarn(_ *testing.T) {
	logger.Warn("warn message")
}

func TestError(_ *testing.T) {
	logger.Error("error message")
}

func TestDebugf(_ *testing.T) {
	log.SetLevel(log.DebugLevel)
	logger.Debugf("debug %s %d", "test", 123)
}

func TestInfof(_ *testing.T) {
	logger.Infof("info %s %d", "test", 123)
}

func TestWarnf(_ *testing.T) {
	logger.Warnf("warn %s %d", "test", 123)
}

func TestErrorf(_ *testing.T) {
	logger.Errorf("error %s %d", "test", 123)
}

func TestSetLevel(_ *testing.T) {
	logger.SetLevel(log.WarnLevel)
	logger.SetLevel(log.InfoLevel)
	logger.SetLevel(log.DebugLevel)
}

func TestPaddingFunctions(_ *testing.T) {
	logger.IncreasePadding()
	logger.DecreasePadding()
}

func TestWithError(t *testing.T) {
	testErr := errors.New("test error")
	entry := logger.WithError(testErr)

	if entry == nil {
		t.Error("WithError returned nil")
	}

	// Test that we can use the entry
	entry.Info("test message with error")
}

func TestColorFunctions(t *testing.T) {
	redText := logger.Red("test")
	expectedRed := "\033[31mtest\033[0m"
	if redText != expectedRed {
		t.Errorf("red() output mismatch: got %q, want %q", redText, expectedRed)
	}

	yellowText := logger.Yellow("test")
	expectedYellow := "\033[33mtest\033[0m"
	if yellowText != expectedYellow {
		t.Errorf("yellow() output mismatch: got %q, want %q", yellowText, expectedYellow)
	}
}
