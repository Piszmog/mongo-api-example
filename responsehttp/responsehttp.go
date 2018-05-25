package responsehttp

import (
    "net/http"
    "encoding/json"
    "go.uber.org/zap"
    "github.com/Piszmog/mongo-api-example/util"
)

const (
    KeyContentType       = "Content-Type"
    ValueApplicationJson = "application/json"
)

var logger zap.SugaredLogger

// Create the logger for the package
func init() {
    zapLogger := util.CreateLogger()
    defer zapLogger.Sync()
    logger = *zapLogger.Sugar()
}

// Writes an OK (200) response to the client
func WriteOkResponse(writer http.ResponseWriter, payload interface{}) {
    bytes, err := json.Marshal(payload)
    if err != nil {
        logger.Errorf("failed to marshal %v", payload)
    }
    WriteResponse(writer, http.StatusOK, bytes)
}

// Writes a response to the client with the provided status and body
func WriteResponse(writer http.ResponseWriter, httpStatus int, bytes []byte) {
    writer.Header().Set(KeyContentType, ValueApplicationJson)
    writer.WriteHeader(httpStatus)
    if bytes != nil {
        writer.Write(bytes)
    }
}
