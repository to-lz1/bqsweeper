package util

import (
	"io"
	"log/slog"
)

// CloseWithLog closes the given resource and logs an error if it occurs.
func CloseWithLog(c io.Closer) {
	if err := c.Close(); err != nil {
		slog.Error("Failed to close resource", "error", err)
	}
}
