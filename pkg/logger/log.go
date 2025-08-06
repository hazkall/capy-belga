package logger

import (
	"log/slog"
	"os"
)

func Start(debug string) {

	l := new(slog.HandlerOptions)

	if debug == "true" {
		l.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, l))

	slog.SetDefault(logger)
}
