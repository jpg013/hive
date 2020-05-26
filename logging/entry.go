package logging

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// Entry type represents a log entry
type Entry struct {
	Level Level

	// Time at which the log entry was created
	Time time.Time

	// Contains all the fields set by the user.
	Data *Fields

	// Contains log message
	Message string

	Logger *Type
}

func (entry *Entry) log() {
	// set timestamp on log
	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}

	entry.write()
}

// Log method for entry
func (entry *Entry) Log(level Level, msg string) {
	entry.Data = CopyDataFields(entry.Data, entry.Logger.Defaults...)
	entry.Message = msg
	entry.Level = level
	entry.log()
}

func (entry *Entry) write() {
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()

	entry.Data.Message = entry.Message
	serialized, err := entry.Logger.Format(entry)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return
	}

	for _, t := range entry.Logger.transports {
		if entry.IsLogLevelEnabled(t.GetLogLevel()) {
			t.Out.Write(serialized)
		}
	}

	// Release the entry once it has written
	entry.Logger.releaseEntry(entry)
}

// NewEntry returns new log entry
func NewEntry(logger *Type) *Entry {
	return &Entry{
		Logger: logger,

		Data: &Fields{},
	}
}

// WithFields takes field pairs and copies them to entry data
func (entry *Entry) WithFields(args ...*FieldPair) *Entry {
	entry.Data = CopyDataFields(entry.Data, args...)
	return entry
}

// IsLogLevelEnabled checks if the log level of the logger is greater than the level param
func (entry *Entry) IsLogLevelEnabled(level Level) bool {
	return entry.GetLogLevel() >= level
}

// GetLogLevel coerces entry level to int32 and returns value
func (entry *Entry) GetLogLevel() Level {
	return Level(atomic.LoadUint32((*uint32)(&entry.Level)))
}

// Info logs message at info level.
func (entry *Entry) Info(msg string) {
	entry.Log(InfoLevel, msg)
}

// Error logs message at error level.
func (entry *Entry) Error(msg string) {
	entry.Log(ErrorLevel, msg)
}

// Warn logs message at warn level.
func (entry *Entry) Warn(msg string) {
	entry.Log(WarnLevel, msg)
}

// Debug logs message at debug level.
func (entry *Entry) Debug(msg string) {
	entry.Log(DebugLevel, msg)
}

// Verbose logs message at verbose level.
func (entry *Entry) Verbose(msg string) {
	entry.Log(VerboseLevel, msg)
}
