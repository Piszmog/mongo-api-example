package model

type Movie struct {
    Id          string `bson:"_id" json:"_id"`
    Name        string `bson:"name" json:"name"`
    Description string `bson:"description" json:"description"`
}

type ResponseId struct {
    Id string `json:"id"`
}

type CloudFoundryEnvironment struct {
    Mlab []MLab `json:"mlab"`
}
type MLab struct {
    Credentials      Credentials `json:"credentials"`

}
type Credentials struct {
    Uri string `json:"uri"`
}