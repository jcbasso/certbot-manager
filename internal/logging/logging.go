package logging

import (
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
)

// Setup initializes the global logrus logger with the specified level.
func Setup(levelStr string) error {
	logrus.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.00",
		LogFormat:       "%time% [%lvl%] %msg%\n",
	})

	logrus.SetOutput(os.Stderr)

	// Parse the level string
	level, err := logrus.ParseLevel(strings.ToLower(levelStr))
	if err != nil {
		log.Printf("Warning: Invalid log level '%s' provided: %v. Defaulting to 'info'.", levelStr, err)
		level = logrus.InfoLevel
	}

	logrus.SetLevel(level)

	return nil
}

// logrusWriter adapts logrus entry to io.Writer for standard logger
type logrusWriter struct {
	entry *logrus.Entry
	level logrus.Level
}

// Write implements io.Writer, logging messages at the specified level
func (w *logrusWriter) Write(p []byte) (n int, err error) {
	// Trim trailing newline standard logger often adds, as logrus adds its own
	msg := strings.TrimSuffix(string(p), "\n")
	switch w.level {
	case logrus.PanicLevel:
		w.entry.Panic(msg)
	case logrus.FatalLevel:
		w.entry.Fatal(msg)
	case logrus.ErrorLevel:
		w.entry.Error(msg)
	case logrus.WarnLevel:
		w.entry.Warn(msg)
	case logrus.InfoLevel:
		w.entry.Info(msg)
	case logrus.DebugLevel:
		w.entry.Debug(msg)
	case logrus.TraceLevel:
		w.entry.Trace(msg)
	default:
		w.entry.Print(msg) // Default fallback
	}
	return len(p), nil
}

// NewLogrusStandardLogger creates a standard library logger that writes to logrus
// at the specified level. Prefix is handled by logrus, flags usually 0.
func NewLogrusStandardLogger(level logrus.Level, component string) *log.Logger {
	entry := logrus.WithField("source", component) // Add a field to distinguish source
	writer := &logrusWriter{entry: entry, level: level}
	return log.New(writer, "", 0) // No prefix or flags needed from stdlib logger
}
