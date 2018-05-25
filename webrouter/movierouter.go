package webrouter

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
    "log"
    "github.com/Piszmog/mongo-api-example/responsehttp"
    "github.com/Piszmog/mongo-api-example/model"
    "encoding/json"
    "github.com/google/uuid"
    "github.com/Piszmog/mongo-api-example/db"
    "github.com/Piszmog/mongo-api-example/config"
    "os"
)

const (
    Id         = "id"
    CfServices = "VCAP_SERVICES"
)

var dbConnection db.DBConnection

func init() {
    var configReader config.Config
    cfServices := os.Getenv(CfServices)
    if len(cfServices) == 0 {
        configReader.Read()
        dbConnection.Connect(configReader.Server, configReader.Database)
    } else {
        var env model.CFEnv
        err := json.Unmarshal([]byte(cfServices), &env)
        if err != nil {
            log.Fatal("Failed to convert env map", err)
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
        log.Println(err)
        responsehttp.WriteResponse(writer, http.StatusInternalServerError, nil)
        return
    }
    responsehttp.WriteOkResponse(writer, movies)
}

func FindMovie(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
    movie, err := dbConnection.FindById(params.ByName(Id))
    if err != nil {
        log.Println(err)
        responsehttp.WriteResponse(writer, http.StatusInternalServerError, nil)
        return
    }
    responsehttp.WriteOkResponse(writer, movie)
}

func CreateMovie(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
    defer request.Body.Close()
    var movie model.Movie
    if err := json.NewDecoder(request.Body).Decode(&movie); err != nil {
        log.Println(err)
        responsehttp.WriteResponse(writer, http.StatusInternalServerError, nil)
        return
    }
    movie.Id = uuid.New().String()
    if err := dbConnection.Insert(movie); err != nil {
        log.Println(err)
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
        log.Println(err)
        responsehttp.WriteResponse(writer, http.StatusInternalServerError, nil)
        return
    }
    movie.Id = movieId
    err := dbConnection.Update(movieId, movie)
    if err != nil {
        log.Println(err)
        responsehttp.WriteResponse(writer, http.StatusInternalServerError, nil)
        return
    }
    writer.WriteHeader(http.StatusOK)
}

func DeleteMovie(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
    movieId := params.ByName(Id)
    err := dbConnection.Delete(movieId)
    if err != nil {
        log.Println(err)
        responsehttp.WriteResponse(writer, http.StatusInternalServerError, nil)
        return
    }
    responsehttp.WriteResponse(writer, http.StatusOK, nil)
}
