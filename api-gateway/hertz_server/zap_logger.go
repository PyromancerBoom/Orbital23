package main

/*
Reusable logger for the Hertz server

Methods for loggging :
 - Debug(msg string, fields ...Field): Logs a debug message
 - Info(msg string, fields ...Field): Logs an info message
 - Warn(msg string, fields ...Field): Logs a warning message
 - Error(msg string, fields ...Field): Logs an error message
 - Fatal(msg string, fields ...Field): Logs a message at fatal level and then calls os.Exit(1)

Import "go.uber.org/zap" to use.

 -----> Use by calling the methods on zap.L()
 		For example, zap.L().Info(...)
*/

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Replace the global logger with our own.
	zap.ReplaceGlobals(logger)

	// Use the logger in other files.
}
