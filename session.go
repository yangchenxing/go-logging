package logging

import (
	"fmt"
	"time"
)

// Session hold special field of user defined session.
type Session map[string]interface{}

var (
	defaultSession = Session(nil)
)

// NewSession create a new session.
func NewSession() Session {
	return make(map[string]interface{})
}

// Set set a key-value pair to session
func (s Session) Set(key string, value interface{}) {
	s[key] = value
}

// LogWithSkip log message with specified skip.
func (s Session) LogWithSkip(skip int, level string, format string, args ...interface{}) {
	targetHandlers := handlers[level]
	if len(targetHandlers) == 0 {
		return
	}
	localStaticContext := getLocalStaticContext(skip + 1)
	dynamicContext := map[string]interface{}{
		"level":   level,
		"time":    time.Now().Format(timeFormat),
		"message": fmt.Sprintf(format, args...),
	}
	if s == nil {
		for _, handler := range targetHandlers {
			handler.write(globalStaticContext, localStaticContext, dynamicContext)
		}
	} else {
		for _, handler := range targetHandlers {
			handler.write(globalStaticContext, localStaticContext, dynamicContext, s)
		}
	}
}

// Log log message with any level
func (s Session) Log(level string, format string, args ...interface{}) {
	s.LogWithSkip(2, level, format, args...)
}

// Debug log message with "debug" level
func (s Session) Debug(format string, args ...interface{}) {
	s.LogWithSkip(2, "debug", format, args...)
}

// Info log message with "info" level
func (s Session) Info(format string, args ...interface{}) {
	s.LogWithSkip(2, "info", format, args...)
}

// Warn log message with "warn" level
func (s Session) Warn(format string, args ...interface{}) {
	s.LogWithSkip(2, "warn", format, args...)
}

// Error log message with "error" level
func (s Session) Error(format string, args ...interface{}) {
	s.LogWithSkip(2, "error", format, args...)
}

// Fatal log message with "fatal" level
func (s Session) Fatal(format string, args ...interface{}) {
	s.LogWithSkip(2, "fatal", format, args...)
}
