package model

type Profile struct {
	UUID         string `json:"uuid" bson:"uuid"`
	Name         string `json:"name" bson:"name"`
	Surname      string `json:"surname" bson:"surname"`
	Email        string `json:"email" bson:"email"`
	PersonalCode string `json:"personal_code" bson:"personal_code"`
}