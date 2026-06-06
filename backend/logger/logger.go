package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelDebug Level = "debug"
	LevelFatal Level = "fatal"
)

type Entry struct {
	Timestamp string `json:"timestamp"`
	Level     Level  `json:"level"`
	Message   string `json:"message"`
}

func log(level Level, msg string, args ...interface{}) {
	entry := Entry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   fmt.Sprintf(msg, args...),
	}
	data, _ := json.Marshal(entry)
	out := os.Stdout
	if level == LevelError || level == LevelFatal {
		out = os.Stderr
	}
	fmt.Fprintln(out, string(data))
}

func Info(msg string, args ...interface{})  { log(LevelInfo, msg, args...) }
func Warn(msg string, args ...interface{})  { log(LevelWarn, msg, args...) }
func Error(msg string, args ...interface{}) { log(LevelError, msg, args...) }
func Debug(msg string, args ...interface{}) { log(LevelDebug, msg, args...) }

func Fatal(msg string, args ...interface{}) {
	log(LevelFatal, msg, args...)
	os.Exit(1)
}
