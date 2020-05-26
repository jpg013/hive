package logging

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Level type
type Level uint32

const (
	// ErrorLevel level. Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel Level = iota
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the application.
	InfoLevel
	// VerboseLevel level. Usually only enabled when debugging. Very verbose logging.
	VerboseLevel
	// DebugLevel level. Usually only enabled when debugging. Even more verbose than verbose logging.
	DebugLevel
)

var (
	defaultLevel Level = InfoLevel
)

// Logger represents a logger interface
type Logger interface {
	WithLevel(Level) Logger
	WithDefaults(...*FieldPair) Logger
	WithFields(args ...*FieldPair) *Entry
	WithTransports(...Transport) Logger
	Info(string)
	Error(string)
	Warn(string)
	Debug(string)
	Verbose(string)
	Log(Level, string)
}

func levelToString(lvl Level) (string, error) {
	switch lvl {
	case DebugLevel:
		return "debug", nil
	case VerboseLevel:
		return "verbose", nil
	case InfoLevel:
		return "info", nil
	case WarnLevel:
		return "warn", nil
	case ErrorLevel:
		return "error", nil
	}

	return "", fmt.Errorf("not a valid log Level: %q", lvl)
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "verbose":
		return VerboseLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid log Level: %q", lvl)
}

// Type implements the Logger interface.
type Type struct {
	// List of transports.
	transports []Transport

	// The logging level the logger should log at.
	Level Level

	// Used for syncing transport writing.
	mu sync.Mutex

	// Reusable empty entry for managing multiple concurrent items.
	entryPool sync.Pool

	// Default log fields.
	Defaults []*FieldPair
}

// NewLogger returns a new logger instance
func NewLogger() Logger {
	return &Type{
		Level:    defaultLevel,
		Defaults: make([]*FieldPair, 0),
	}
}

// WithLevel Sets the log level for logger
func (t *Type) WithLevel(lvl Level) Logger {
	t.Level = lvl
	return t
}

// WithTransports takes a variadic number of transports and adds them to the logger tranpsorts
func (t *Type) WithTransports(args ...Transport) Logger {
	t.transports = append(t.transports, args...)
	return t
}

func isZeroString(s string) bool {
	return s == ""
}

func isZeroInt(i int) bool {
	return i == 0
}

func isZeroFloat(i float64) bool {
	return i == 0
}

// Info logs message at info level.
func (t *Type) Info(msg string) {
	t.Log(InfoLevel, msg)
}

// Error logs message at error level.
func (t *Type) Error(msg string) {
	t.Log(ErrorLevel, msg)
}

// Warn logs message at warn level.
func (t *Type) Warn(msg string) {
	t.Log(WarnLevel, msg)
}

// Debug logs message at debug level.
func (t *Type) Debug(msg string) {
	t.Log(DebugLevel, msg)
}

// Verbose logs message at verbose level.
func (t *Type) Verbose(msg string) {
	t.Log(VerboseLevel, msg)
}

// newEntry attempts to get an entry from the pool or creates a new one
func (t *Type) newEntry() *Entry {
	entry, ok := t.entryPool.Get().(*Entry)

	if ok {
		return entry
	}

	return NewEntry(t)
}

// Log creates a new entry and calls log with level and fields
func (t *Type) Log(level Level, msg string) {
	entry := t.newEntry()
	entry.Log(level, msg)
}

// releaseEntry adds entry back in
func (t *Type) releaseEntry(entry *Entry) {
	// Reset the data fields
	entry.Data = &Fields{}
	entry.Time = time.Time{} // Invoking an empty time.Time struct literal will return Go's zero date.
	t.entryPool.Put(entry)
}

// WithDefaults set the default fields for a log entry
func (t *Type) WithDefaults(args ...*FieldPair) Logger {
	for _, arg := range args {
		// Check if default exists
		exists := false
		for _, f := range t.Defaults {
			if arg.Name == f.Name {
				exists = true
				// Override the existing default value
				f.Value = arg.Value
				break
			}
		}

		if !exists {
			t.Defaults = append(t.Defaults, arg)
		}
	}

	return t
}

// Format transforms entry data into a byte array. Can be
// extended later to include different formatter types
func (t *Type) Format(entry *Entry) ([]byte, error) {
	entry.Data.Timestamp = entry.Time.Format(time.RFC3339)
	sev, err := levelToString(entry.Level)
	if err != nil {
		return nil, err
	}
	entry.Data.Severity = sev
	fmt.Println("what in the holiest of fucks!")
	fmt.Println(entry.Data.Message)
	return json.Marshal(entry.Data)
}

// WithFields takes field pairs and reutrns an entry
func (t *Type) WithFields(args ...*FieldPair) *Entry {
	entry := t.newEntry()
	return entry.WithFields(args...)
}
