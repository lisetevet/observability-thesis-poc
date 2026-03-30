package model

type User struct {
	Username string `json:"username" bson:"username"`
	UUID     string `json:"uuid" bson:"uuid"`
}