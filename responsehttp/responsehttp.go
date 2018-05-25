package responsehttp

import (
    "net/http"
    "encoding/json"
    "go.uber.org/zap"
    "log"
)

const (
    KeyContentType       = "Content-Type"
    ValueApplicationJson = "application/json"
)

var logger zap.SugaredLogger

func init() {
    log1, err := zap.NewProduction()
    if err != nil {
        log.Fatalf("failed to create zap logger, %v", err)
    }
    defer log1.Sync()
    logger = *log1.Sugar()
}

func WriteOkResponse(writer http.ResponseWriter, payload interface{}) {
    bytes, err := json.Marshal(payload)
    if err != nil {
        logger.Errorf("failed to marshal %v", payload)
    }
    WriteResponse(writer, http.StatusOK, bytes)
}

func WriteResponse(writer http.ResponseWriter, httpStatus int, bytes []byte) {
    writer.Header().Set(KeyContentType, ValueApplicationJson)
    writer.WriteHeader(httpStatus)
    if bytes != nil {
        writer.Write(bytes)
    }
}
