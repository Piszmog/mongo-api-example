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
    "log"
)

const (
    CfServices      = "VCAP_SERVICES"
    DefaultDatabase = "test"
    DefaultServer   = "localhost"
    Id              = "id"
)

var dbConnection db.Connection
var logger zap.SugaredLogger

func init() {
    log1, err := zap.NewProduction()
    if err != nil {
        log.Fatalf("failed to create zap logger, %v", err)
    }
    defer log1.Sync()
    logger = *log1.Sugar()
    cfServices := os.Getenv(CfServices)
    if len(cfServices) == 0 {
        dbConnection.Connect(DefaultServer, DefaultDatabase)
    } else {
        var env model.CFEnv
        err := json.Unmarshal([]byte(cfServices), &env)
        if err != nil {
            logger.Fatal("failed to convert env map", err)
        }
        dbConnection.ConnectWithURL(env.Mlab[0].Credentials.Uri)
    }
}

func SetupMovieRoutes(router *httprouter.Router) {
    router.GET("/movies", GetAllMovies)
    router.GET("/movies/:id", FindMovie)
    router.POST("/movies", CreateMovie)
    router.PUT("/movies/:id", UpdateMovie)
    router.DELETE("/movies/:id", DeleteMovie)
}

func GetAllMovies(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
    movies, err := dbConnection.FindAll()
    if err != nil {
        logger.Errorf("failed to get all movies, %v", err)
        responsehttp.WriteResponse(writer, http.StatusInternalServerError, nil)
        return
    }
    responsehttp.WriteOkResponse(writer, movies)
}

func FindMovie(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
    movie, err := dbConnection.FindById(params.ByName(Id))
    if err != nil {
        logger.Warnf("failed to find movie, %v", err)
        responsehttp.WriteResponse(writer, http.StatusNotFound, nil)
        return
    }
    responsehttp.WriteOkResponse(writer, movie)
}

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
        logger.Errorf("failed to update movie %v, %v", movieId, err)
        responsehttp.WriteResponse(writer, http.StatusInternalServerError, nil)
        return
    }
    writer.WriteHeader(http.StatusOK)
}

func DeleteMovie(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
    movieId := params.ByName(Id)
    err := dbConnection.Delete(movieId)
    if err != nil {
        logger.Errorf("failed to delete movie %v, %v", movieId, err)
        responsehttp.WriteResponse(writer, http.StatusInternalServerError, nil)
        return
    }
    responsehttp.WriteResponse(writer, http.StatusOK, nil)
}
