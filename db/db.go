package db

import (
    "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    "github.com/Piszmog/mongo-api-example/model"
    "go.uber.org/zap"
    "log"
)

type Connection struct {
    Server   string
    Database string
    db       *mgo.Database
}

const (
    COLLECTION = "movies"
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

func (m *Connection) Connect(server string, database string) {
    session, err := mgo.Dial(m.Server)
    if err != nil {
        logger.Fatalf("failed to connect to the db, %v", err)
    }
    m.Server = server
    m.Database = database
    m.db = session.DB(m.Database)
}

func (m *Connection) ConnectWithURL(url string) {
    info, parseErr := mgo.ParseURL(url)
    if parseErr != nil {
        logger.Fatalf("failed to parse URL, %v", parseErr)
    }
    session, err := mgo.DialWithInfo(info)
    if err != nil {
        logger.Fatalf("failed to connect to the db, %v", err)
    }
    m.Database = info.Database
    m.db = session.DB(m.Database)
}

func (m *Connection) FindAll() ([]model.Movie, error) {
    var movies []model.Movie
    err := m.db.C(COLLECTION).Find(bson.M{}).All(&movies)
    return movies, err
}

func (m *Connection) FindById(id string) (model.Movie, error) {
    var movie model.Movie
    err := m.db.C(COLLECTION).FindId(id).One(&movie)
    return movie, err
}

func (m *Connection) Insert(movie model.Movie) error {
    err := m.db.C(COLLECTION).Insert(&movie)
    return err
}

func (m *Connection) Delete(id string) error {
    err := m.db.C(COLLECTION).RemoveId(id)
    return err
}

func (m *Connection) Update(id string, movie model.Movie) error {
    err := m.db.C(COLLECTION).UpdateId(id, &movie)
    return err
}
