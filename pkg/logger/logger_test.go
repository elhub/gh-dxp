package logger

import (
	"errors"
	"testing"

	"github.com/caarlos0/log"
)

func TestDebug(_ *testing.T) {
	log.SetLevel(log.DebugLevel)
	// Just verify it doesn't panic
	Debug("debug message")
}

func TestInfo(_ *testing.T) {
	Info("info message")
}

func TestWarn(_ *testing.T) {
	Warn("warn message")
}

func TestError(_ *testing.T) {
	Error("error message")
}

func TestDebugf(_ *testing.T) {
	log.SetLevel(log.DebugLevel)
	Debugf("debug %s %d", "test", 123)
}

func TestInfof(_ *testing.T) {
	Infof("info %s %d", "test", 123)
}

func TestWarnf(_ *testing.T) {
	Warnf("warn %s %d", "test", 123)
}

func TestErrorf(_ *testing.T) {
	Errorf("error %s %d", "test", 123)
}

func TestSetLevel(_ *testing.T) {
	SetLevel(log.WarnLevel)
	SetLevel(log.InfoLevel)
	SetLevel(log.DebugLevel)
}

func TestPaddingFunctions(_ *testing.T) {
	IncreasePadding()
	DecreasePadding()
}

func TestWithError(t *testing.T) {
	testErr := errors.New("test error")
	entry := WithError(testErr)

	if entry == nil {
		t.Error("WithError returned nil")
	}

	// Test that we can use the entry
	entry.Info("test message with error")
}

func TestColorFunctions(t *testing.T) {
	redText := Red("test")
	expectedRed := "\033[31mtest\033[0m"
	if redText != expectedRed {
		t.Errorf("red() output mismatch: got %q, want %q", redText, expectedRed)
	}

	yellowText := Yellow("test")
	expectedYellow := "\033[33mtest\033[0m"
	if yellowText != expectedYellow {
		t.Errorf("yellow() output mismatch: got %q, want %q", yellowText, expectedYellow)
	}
}
