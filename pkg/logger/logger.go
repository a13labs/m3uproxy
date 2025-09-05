package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// LogLevel constants
const (
	PanicLevel = "panic"
	FatalLevel = "fatal"
	ErrorLevel = "error"
	WarnLevel  = "warn"
	InfoLevel  = "info"
	DebugLevel = "debug"
	TraceLevel = "trace"
)

// Init initializes the logger with file output and configures log level from environment variables
func Init(logFile string) {

	if log != nil {
		// Logger already initialized
		return
	}

	// Check for various environment variable names
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = os.Getenv("LOGLEVEL")
	}
	if level == "" {
		level = os.Getenv("M3UPROXY_LOG_LEVEL")
	}

	if level == "" {
		level = "info"
	}

	InitWithLevel(logFile, level)
}

// InitWithLevel initializes the logger with a specific log level
func InitWithLevel(logFile string, level string) {
	if log != nil {
		// Logger already initialized
		return
	}

	log = &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     getLogLevelFromString(level),
	}
	// Configure formatter
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Configure output
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Warn("Failed to log to file, using default stderr")
		}
	}
}

// getLogLevel retrieves the log level from environment variables
func getLogLevelFromString(level string) logrus.Level {
	switch level {
	case PanicLevel:
		return logrus.PanicLevel
	case FatalLevel:
		return logrus.FatalLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case WarnLevel, "warning":
		return logrus.WarnLevel
	case InfoLevel:
		return logrus.InfoLevel
	case DebugLevel:
		return logrus.DebugLevel
	case TraceLevel:
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}

// Logging functions
func Trace(args ...interface{}) {
	log.Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	log.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}
