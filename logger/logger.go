package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init() {

	opts := &slog.HandlerOptions{
		AddSource: true,
	}

	logHandler := slog.NewTextHandler(os.Stdout, opts)

	Log = slog.New(logHandler)

	slog.SetDefault(Log)

}
