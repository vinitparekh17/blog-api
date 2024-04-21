package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

var Log *slog.Logger

func Init() {

	env := os.Getenv("ENV")

	if env == "development" {

		Log = slog.New(tint.NewHandler(os.Stdout, nil))

	} else {

		opts := &slog.HandlerOptions{
			AddSource: true,
		}

		logHandler := slog.NewTextHandler(os.Stdout, opts)

		Log = slog.New(logHandler)
	}
	slog.SetDefault(Log)

}
