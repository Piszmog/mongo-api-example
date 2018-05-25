package webrouter

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
    "github.com/Piszmog/mongo-api-example/responsehttp"
    "github.com/Piszmog/mongo-api-example/model"
    "encoding/json"
    "github.com/google/uuid"
    "github.com/Piszmog/mongo-api-example/db"
    "os"
    "go.uber.org/zap"
    "github.com/Piszmog/mongo-api-example/util"
)

const (
    CfServices      = "VCAP_SERVICES"
    DefaultDatabase = "test"
    DefaultServer   = "localhost"
    Id              = "id"
)

var dbConnection db.Connection
var logger zap.SugaredLogger

// Creates the logger and creates the connection to the DB
func init() {
    zapLogger := util.CreateLogger()
    defer zapLogger.Sync()
    logger = *zapLogger.Sugar()
    connectToDB()
}

// Connects to the db. If "VCAP_SERVICES" is not an environment variable, it connects to the default server name "localhost"
// and the default database "test"
func connectToDB() {
    cfServices := os.Getenv(CfServices)
    if len(cfServices) == 0 {
        dbConnection.Connect(DefaultServer, DefaultDatabase)
    } else {
        var env model.CloudFoundryEnvironment
        err := json.Unmarshal([]byte(cfServices), &env)
        if err != nil {
            logger.Fatal("failed to convert env map", err)
        }
        dbConnection.ConnectWithURL(env.Mlab[0].Credentials.Uri)
    }
}

// Sets up the http routes
func SetupMovieRoutes(router *httprouter.Router) {
    router.GET("/movies", GetAllMovies)
    router.GET("/movies/:id", FindMovie)
    router.POST("/movies", CreateMovie)
    router.PUT("/movies/:id", UpdateMovie)
    router.DELETE("/movies/:id", DeleteMovie)
}

// Retrieves all the movies from the DB
func GetAllMovies(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
    movies, err := dbConnection.FindAll()
    if err != nil {
        logger.Errorf("failed to get all movies, %v", err)
        responsehttp.WriteResponse(writer, http.StatusInternalServerError, nil)
        return
    }
    responsehttp.WriteOkResponse(writer, movies)
}

// Finds the movie matching the provided id
func FindMovie(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
    movie, err := dbConnection.FindById(params.ByName(Id))
    if err != nil {
        logger.Warnf("failed to find movie, %v", err)
        responsehttp.WriteResponse(writer, http.StatusNotFound, nil)
        return
    }
    responsehttp.WriteOkResponse(writer, movie)
}

// Creates the movie provided by the request body
func CreateMovie(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
    defer request.Body.Close()
    var movie model.Movie
    if err := json.NewDecoder(request.Body).Decode(&movie); err != nil {
        logger.Errorf("failed to decode request body, %v", err)
        responsehttp.WriteResponse(writer, http.StatusInternalServerError, nil)
        return
    }
    movie.Id = uuid.New().String()
    if err := dbConnection.Insert(movie); err != nil {
        logger.Errorf("failed to create movie, %v", err)
        responsehttp.WriteResponse(writer, http.StatusInternalServerError, nil)
        return
    }
    responsehttp.WriteOkResponse(writer, model.ResponseId{Id: movie.Id})
}

// Updates the movie matching the provided id with the provided body
func UpdateMovie(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
    defer request.Body.Close()
    movieId := params.ByName(Id)
    var movie model.Movie
    if err := json.NewDecoder(request.Body).Decode(&movie); err != nil {
        logger.Errorf("failed to decode request body %v, %v", movieId, err)
        responsehttp.WriteResponse(writer, http.StatusInternalServerError, nil)
        return
    }
    movie.Id = movieId
    err := dbConnection.Update(movieId, movie)
    if err != nil {
        logger.Warnf("failed to update movie %v, %v", movieId, err)
        responsehttp.WriteResponse(writer, http.StatusNotFound, nil)
        return
    }
    writer.WriteHeader(http.StatusOK)
}

// Deletes the movie matching the provided id
func DeleteMovie(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
    movieId := params.ByName(Id)
    err := dbConnection.Delete(movieId)
    if err != nil {
        logger.Warnf("failed to delete movie %v, %v", movieId, err)
        responsehttp.WriteResponse(writer, http.StatusNotFound, nil)
        return
    }
    responsehttp.WriteResponse(writer, http.StatusOK, nil)
}
