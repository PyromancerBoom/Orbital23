package main

// Reusable logger for the Hertz server
/*
Methods for loggging :
 - Debug(msg string, fields ...Field): Logs a debug message.
 - Info(msg string, fields ...Field): Logs an info message.
 - Warn(msg string, fields ...Field): Logs a warning message.
 - Error(msg string, fields ...Field): Logs an error message.
 - DPanic(msg string, fields ...Field): Logs a message at panic level and then panics.
 - Panic(msg string, fields ...Field): Logs a message at panic level and then panics.
 - Fatal(msg string, fields ...Field): Logs a message at fatal level and then calls os.Exit(1).

Import "go.uber.org/zap" to use.

Use by calling the methods on zap.L()
*/

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func initLogger() {
	logger, _ = zap.NewProduction()
	defer logger.Sync()

	// Replace the global logger with our own.
	zap.ReplaceGlobals(logger)

	// Use the logger in other files.
}
