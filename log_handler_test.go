package nslog

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

///////////////////////////////////////////////////////////////////////////////
// Basic Use
///////////////////////////////////////////////////////////////////////////////

const DEFAULT_TIME_REGEXP = "\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2}"

func TestBasicUse1(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, nil)
	log.Info("log message", "key1", "val1")
	assert.Regexp(t, DEFAULT_TIME_REGEXP+" INFO\\. log message key1=val1", buf.String())
}

func TestBasicUse2(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, nil)
	log.Info("log message")
	assert.Regexp(t, DEFAULT_TIME_REGEXP+" INFO\\. log message", buf.String())
}

func TestBasicUse3(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, nil)
	log.Info("log message", "key1", "val1", "key2", "val2")
	assert.Regexp(t, DEFAULT_TIME_REGEXP+" INFO\\. log message key1=val1 key2=val2", buf.String())
}

///////////////////////////////////////////////////////////////////////////////
// Option: Level
///////////////////////////////////////////////////////////////////////////////

func TestDefaultLevel(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, nil)
	log.Error("log message")
	log.Warn("log message")
	log.Info("log message")
	log.Debug("log message")
	assert.Contains(t, buf.String(), "ERROR log message")
	assert.Contains(t, buf.String(), "WARN. log message")
	assert.Contains(t, buf.String(), "INFO. log message")
	assert.NotContains(t, buf.String(), "DEBUG log message")
}

func TestDefaultLevel2(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, &LogHandlerOptions{})
	log.Error("log message")
	log.Warn("log message")
	log.Info("log message")
	log.Debug("log message")
	assert.Contains(t, buf.String(), "ERROR log message")
	assert.Contains(t, buf.String(), "WARN. log message")
	assert.Contains(t, buf.String(), "INFO. log message")
	assert.NotContains(t, buf.String(), "DEBUG log message")
}

func TestErrorLevel(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, &LogHandlerOptions{Level: slog.LevelError})
	log.Error("log message")
	log.Warn("log message")
	log.Info("log message")
	log.Debug("log message")
	assert.Contains(t, buf.String(), "ERROR log message")
	assert.NotContains(t, buf.String(), "WARN. log message")
	assert.NotContains(t, buf.String(), "INFO. log message")
	assert.NotContains(t, buf.String(), "DEBUG log message")
}

func TestWarnLevel(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, &LogHandlerOptions{Level: slog.LevelWarn})
	log.Error("log message")
	log.Warn("log message")
	log.Info("log message")
	log.Debug("log message")
	assert.Contains(t, buf.String(), "ERROR log message")
	assert.Contains(t, buf.String(), "WARN. log message")
	assert.NotContains(t, buf.String(), "INFO. log message")
	assert.NotContains(t, buf.String(), "DEBUG log message")
}

func TestInfoLevel(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, &LogHandlerOptions{Level: slog.LevelInfo})
	log.Error("log message")
	log.Warn("log message")
	log.Info("log message")
	log.Debug("log message")
	assert.Contains(t, buf.String(), "ERROR log message")
	assert.Contains(t, buf.String(), "WARN. log message")
	assert.Contains(t, buf.String(), "INFO. log message")
	assert.NotContains(t, buf.String(), "DEBUG log message")
}

func TestDebugLevel(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, &LogHandlerOptions{Level: slog.LevelDebug})
	log.Error("log message")
	log.Warn("log message")
	log.Info("log message")
	log.Debug("log message")
	assert.Contains(t, buf.String(), "ERROR log message")
	assert.Contains(t, buf.String(), "WARN. log message")
	assert.Contains(t, buf.String(), "INFO. log message")
	assert.Contains(t, buf.String(), "DEBUG log message")
}

///////////////////////////////////////////////////////////////////////////////
// Groups
///////////////////////////////////////////////////////////////////////////////

func TestWithGroup(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, nil).WithGroup("Group1")
	log.Info("message")
	assert.Contains(t, buf.String(), "INFO. Group1: message")
}

func TestWithGroup2(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, nil).WithGroup("Group1").WithGroup("Group2")
	log.Info("message")
	assert.Contains(t, buf.String(), "INFO. Group1.Group2: message")
}

func TestWithAttrs(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, nil).With("pid", 0)
	log.Info("message")
	assert.Contains(t, buf.String(), "INFO. [pid=0]: message")
}

func TestWithAttrs2(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, nil).With("pid", "dead").With("tid", "beaf")
	log.Info("message")
	assert.Contains(t, buf.String(), "INFO. [pid=dead tid=beaf]: message")
}

func TestWith(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, nil).WithGroup("Group1").With("pid", 0)
	log.Info("message")
	assert.Contains(t, buf.String(), "INFO. Group1[pid=0]: message")
}

func TestWith2(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, nil).WithGroup("Group1").WithGroup("Group2").With("pid", 0)
	log.Info("message")
	assert.Contains(t, buf.String(), "INFO. Group1.Group2[pid=0]: message")

	buf2 := new(bytes.Buffer)
	log2 := NewLogger(buf2, nil).WithGroup("Group1").With("pid", 0).WithGroup("Group2")
	log2.Info("message")
	assert.Contains(t, buf2.String(), "INFO. Group1.Group2[pid=0]: message")
}

func TestWith3(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(buf, nil).WithGroup("Group1").WithGroup("Group2").With("pid", "dead").With("tid", "beaf")
	log.Info("message")
	assert.Contains(t, buf.String(), "INFO. Group1.Group2[pid=dead tid=beaf]: message")

	buf2 := new(bytes.Buffer)
	log2 := NewLogger(buf2, nil).WithGroup("Group1").With("pid", "dead").WithGroup("Group2").With("tid", "beaf")
	log2.Info("message")
	assert.Contains(t, buf2.String(), "INFO. Group1.Group2[pid=dead tid=beaf]: message")
}
