package logger

import (
	"log/slog"
	"os"
)

// global logger
var Log *slog.Logger

// parseLevel parses log level in string
func parseLevel(s string) (slog.Level, error) {
	var lvl slog.Level

	if err := lvl.UnmarshalText([]byte(s)); err != nil {
		return 0, err
	}

	return lvl.Level(), nil
}

// Initialize initiates global slog logger with log level
// log level text: DEBUG, INFO, WARN and ERROR
func Initialize(level string) error {
	// parse level
	lvl, err := parseLevel(level)
	if err != nil {
		return err
	}

	// create new logger
	Log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})).With("server", "gophkeeper")

	return nil
}
