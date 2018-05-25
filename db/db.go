package db

import (
    "github.com/globalsign/mgo"
    "log"
    "github.com/globalsign/mgo/bson"
    "github.com/Piszmog/mongo-api-example/model"
)

type DBConnection struct {
    Server   string
    Database string
    db       *mgo.Database
}

const COLLECTION = "movies"

func (m *DBConnection) Connect(server string, database string) {
    session, err := mgo.Dial(m.Server)
    if err != nil {
        log.Fatal(err)
    }
    m.Server = server
    m.Database = database
    m.db = session.DB(m.Database)
}

func (m *DBConnection) ConnectWithURL(url string) {
    info, _ := mgo.ParseURL(url)
    session, err := mgo.DialWithInfo(info)
    if err != nil {
        log.Fatal(err)
    }
    m.Database = info.Database
    m.db = session.DB(m.Database)
}

func (m *DBConnection) FindAll() ([]model.Movie, error) {
    var movies []model.Movie
    err := m.db.C(COLLECTION).Find(bson.M{}).All(&movies)
    return movies, err
}

func (m *DBConnection) FindById(id string) (model.Movie, error) {
    var movie model.Movie
    err := m.db.C(COLLECTION).FindId(id).One(&movie)
    return movie, err
}

func (m *DBConnection) Insert(movie model.Movie) error {
    err := m.db.C(COLLECTION).Insert(&movie)
    return err
}

func (m *DBConnection) Delete(id string) error {
    err := m.db.C(COLLECTION).RemoveId(id)
    return err
}

func (m *DBConnection) Update(id string, movie model.Movie) error {
    err := m.db.C(COLLECTION).UpdateId(id, &movie)
    return err
}
