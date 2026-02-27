package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Level string

const (
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	SkipLevel  Level = "skip"
	DoneLevel  Level = "done"
	FailLevel  Level = "fail"
	DryRunLevel Level = "dry-run"
)

type Logger struct {
	out     io.Writer
	verbose bool
	json    bool
}

type LogEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

// NewLogger creates a new logger that outputs to stderr.
func NewLogger(verbose, json bool) *Logger {
	return &Logger{
		out:     os.Stderr,
		verbose: verbose,
		json:    json,
	}
}

func (l *Logger) log(level Level, msg string) {
	if l.json {
		entry := LogEntry{
			Level:   string(level),
			Message: msg,
		}
		data, _ := json.Marshal(entry)
		fmt.Fprintln(l.out, string(data))
	} else {
		prefix := fmt.Sprintf("[%s]", level)
		fmt.Fprintf(l.out, "%s %s\n", prefix, msg)
	}
}

func (l *Logger) Info(msg string) {
	if l.verbose {
		l.log(InfoLevel, msg)
	}
}

func (l *Logger) Warn(msg string) {
	l.log(WarnLevel, msg)
}

func (l *Logger) Skip(msg string) {
	l.log(SkipLevel, msg)
}

func (l *Logger) Done(msg string) {
	l.log(DoneLevel, msg)
}

func (l *Logger) Fail(msg string) {
	l.log(FailLevel, msg)
}

func (l *Logger) DryRun(msg string) {
	l.log(DryRunLevel, msg)
}

func (l *Logger) Section(title string) {
	if l.verbose {
		fmt.Fprintf(l.out, "\n=== %s ===\n", title)
	}
}
