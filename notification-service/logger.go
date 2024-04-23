package main

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func (app *application) LoggerInit() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}
