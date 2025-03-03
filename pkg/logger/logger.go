package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Default logger instance
	logger zerolog.Logger
)

// init initializes the logger
func init() {
	// Set up console writer
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	// Configure logger
	logger = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	// Set global logger
	log.Logger = logger

	// Default log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	Info("Logger initialized")
}

// SetOutput sets the logger output
func SetOutput(w io.Writer) {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        w,
		TimeFormat: time.RFC3339,
	}
	logger = logger.Output(consoleWriter)
	log.Logger = logger
}

// SetLevel sets the logger level
func SetLevel(level string) {
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// Debug logs a debug message
func Debug(message string, args ...interface{}) {
	event := logger.Debug()
	appendArgs(event, args...)
	event.Msg(message)
}

// Info logs an info message
func Info(message string, args ...interface{}) {
	event := logger.Info()
	appendArgs(event, args...)
	event.Msg(message)
}

// Warn logs a warning message
func Warn(message string, args ...interface{}) {
	event := logger.Warn()
	appendArgs(event, args...)
	event.Msg(message)
}

// Error logs an error message
func Error(message string, args ...interface{}) {
	event := logger.Error()
	appendArgs(event, args...)
	event.Msg(message)
}

// Fatal logs a fatal message and exits
func Fatal(message string, args ...interface{}) {
	event := logger.Fatal()
	appendArgs(event, args...)
	event.Msg(message)
}

// appendArgs adds key-value pairs to the log event
func appendArgs(event *zerolog.Event, args ...interface{}) {
	if len(args)%2 != 0 {
		event.Str("error", "odd number of arguments passed as key-value pairs")
		return
	}

	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			event.Str("error", fmt.Sprintf("key %v is not a string", args[i]))
			continue
		}

		// Special handling for error objects
		if key == "error" {
			if err, ok := args[i+1].(error); ok {
				event.Err(err)
			} else {
				event.Interface(key, args[i+1])
			}
		} else {
			event.Interface(key, args[i+1])
		}
	}
}
