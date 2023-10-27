package main

import (
	"log/slog"
	"os"
)

var (
	Logr *slog.Logger
)

func SetupLogger() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	Logr = l
}
