package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureOutput(fn func()) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	defer r.Close()

	old := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = old }()

	fn()
	w.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	return buf.String(), err
}

func captureStderr(fn func()) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	defer r.Close()

	old := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = old }()

	fn()
	w.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	return buf.String(), err
}

func parseEntry(t *testing.T, line string) Entry {
	t.Helper()
	var entry Entry
	err := json.Unmarshal([]byte(strings.TrimSpace(line)), &entry)
	require.NoError(t, err)
	return entry
}

func TestInfo(t *testing.T) {
	output, err := captureOutput(func() {
		Info("hello %s", "world")
	})
	require.NoError(t, err)

	entry := parseEntry(t, output)
	assert.Equal(t, LevelInfo, entry.Level)
	assert.Equal(t, "hello world", entry.Message)
	assert.NotEmpty(t, entry.Timestamp)
}

func TestWarn(t *testing.T) {
	output, err := captureOutput(func() {
		Warn("warning: %d", 42)
	})
	require.NoError(t, err)

	entry := parseEntry(t, output)
	assert.Equal(t, LevelWarn, entry.Level)
	assert.Equal(t, "warning: 42", entry.Message)
}

func TestError(t *testing.T) {
	output, err := captureStderr(func() {
		Error("something went wrong: %v", "timeout")
	})
	require.NoError(t, err)

	entry := parseEntry(t, output)
	assert.Equal(t, LevelError, entry.Level)
	assert.Equal(t, "something went wrong: timeout", entry.Message)
}

func TestDebug(t *testing.T) {
	output, err := captureOutput(func() {
		Debug("debug msg %s", "test")
	})
	require.NoError(t, err)

	entry := parseEntry(t, output)
	assert.Equal(t, LevelDebug, entry.Level)
	assert.Equal(t, "debug msg test", entry.Message)
}

func TestFatal(t *testing.T) {
	if os.Getenv("TEST_FATAL") == "1" {
		Fatal("fatal error")
		return
	}

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	oldStderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = oldStderr }()

	osExit = func(code int) {
		// restore so test can continue
		os.Stderr = oldStderr
		w.Close()
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		entry := parseEntry(t, buf.String())
		assert.Equal(t, LevelFatal, entry.Level)
		assert.Equal(t, "fatal error", entry.Message)
		assert.Equal(t, 1, code)
	}
	defer func() { osExit = os.Exit }()

	Fatal("fatal error")
}

func TestLogLevelOutputToStdout(t *testing.T) {
	tests := []struct {
		name  string
		level Level
		fn    func(string, ...interface{})
		pipe  func() (*os.File, *os.File, error)
	}{
		{"info goes to stdout", LevelInfo, func(msg string, args ...interface{}) { Info(msg, args...) }, func() (*os.File, *os.File, error) { return os.Pipe() }},
		{"warn goes to stdout", LevelWarn, func(msg string, args ...interface{}) { Warn(msg, args...) }, func() (*os.File, *os.File, error) { return os.Pipe() }},
		{"debug goes to stdout", LevelDebug, func(msg string, args ...interface{}) { Debug(msg, args...) }, func() (*os.File, *os.File, error) { return os.Pipe() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, err := tt.pipe()
			require.NoError(t, err)
			defer r.Close()

			old := os.Stdout
			os.Stdout = w
			defer func() { os.Stdout = old }()

			tt.fn("test")
			w.Close()

			var buf bytes.Buffer
			_, err = buf.ReadFrom(r)
			require.NoError(t, err)
			assert.NotEmpty(t, buf.String())
		})
	}
}

func TestLogLevelOutputToStderr(t *testing.T) {
	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	old := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = old }()

	Error("stderr test")
	w.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)
	entry := parseEntry(t, buf.String())
	assert.Equal(t, LevelError, entry.Level)
	assert.Equal(t, "stderr test", entry.Message)
}

func TestLogWithNoArgs(t *testing.T) {
	output, err := captureOutput(func() {
		Info("simple message")
	})
	require.NoError(t, err)
	entry := parseEntry(t, output)
	assert.Equal(t, "simple message", entry.Message)
}

func TestLogJSONFormat(t *testing.T) {
	output, err := captureOutput(func() {
		Info("json check")
	})
	require.NoError(t, err)
	assert.True(t, json.Valid([]byte(strings.TrimSpace(output))))

	entry := parseEntry(t, output)
	assert.Equal(t, LevelInfo, entry.Level)
	assert.Equal(t, "json check", entry.Message)
	assert.NotEmpty(t, entry.Timestamp)

	var raw map[string]interface{}
	err = json.Unmarshal([]byte(strings.TrimSpace(output)), &raw)
	require.NoError(t, err)
	assert.Contains(t, raw, "timestamp")
	assert.Contains(t, raw, "level")
	assert.Contains(t, raw, "message")
}


