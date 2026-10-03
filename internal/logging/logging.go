// Package logging configures the daemon's structured process logs.
package logging

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

// New creates a logger. Empty level and format select info and JSON.
// The caller owns the writer; logs are written synchronously without buffering.
func New(w io.Writer, level, format string) (*slog.Logger, error) {
	var severity slog.Level
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "", "info":
		severity = slog.LevelInfo
	case "debug":
		severity = slog.LevelDebug
	case "warn", "warning":
		severity = slog.LevelWarn
	case "error":
		severity = slog.LevelError
	default:
		return nil, fmt.Errorf("invalid LOG_LEVEL %q: use debug, info, warn, or error", level)
	}
	options := &slog.HandlerOptions{Level: severity}
	var handler slog.Handler
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "", "json":
		handler = slog.NewJSONHandler(w, options)
	case "text":
		handler = slog.NewTextHandler(w, options)
	default:
		return nil, fmt.Errorf("invalid LOG_FORMAT %q: use json or text", format)
	}
	return slog.New(handler).With("service", "local-cloud-daemon"), nil
}
