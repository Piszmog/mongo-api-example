package util

import (
    "go.uber.org/zap"
    "log"
)

// Create the logger
func CreateLogger() *zap.Logger {
    zapLogger, err := zap.NewProduction()
    if err != nil {
        log.Fatalf("failed to create zap logger, %v", err)
    }
    return zapLogger
}