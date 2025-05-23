package logger

import (
	"1337bo4rd/internal/adapter/config"
	"log/slog"
	"os"
)

var logger *slog.Logger

func Set(config *config.App) {
	logger = slog.New(
		slog.NewTextHandler(os.Stderr, nil),
	)

	slog.SetDefault(logger)
}
