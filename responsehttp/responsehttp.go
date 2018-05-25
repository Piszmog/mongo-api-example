package responsehttp

import (
    "net/http"
    "encoding/json"
)

const (
    KeyContentType       = "Content-Type"
    ValueApplicationJson = "application/json"
)

func WriteOkResponse(writer http.ResponseWriter, payload interface{}) {
    bytes, _ := json.Marshal(payload)
    WriteResponse(writer, http.StatusOK, bytes)
}

func WriteResponse(writer http.ResponseWriter, httpStatus int, bytes []byte) {
    writer.Header().Set(KeyContentType, ValueApplicationJson)
    writer.WriteHeader(httpStatus)
    if bytes != nil {
        writer.Write(bytes)
    }
}
