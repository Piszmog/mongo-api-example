package main

import (
    "github.com/julienschmidt/httprouter"
    "net/http"
    "log"
    "os"
    "github.com/Piszmog/mongo-api-example/webrouter"
)

const (
    Port        = "PORT"
    DefaultPort = "8080"
)

func getPort() string {
    var port string
    if port = os.Getenv(Port); len(port) == 0 {
        port = DefaultPort
    }
    return port
}

func main() {
    port := getPort()
    router := httprouter.New()
    webrouter.SetupMovieRoutes(router)
    log.Fatal(http.ListenAndServe(":"+port, router))
}
