package model

type User struct {
	Username     string `json:"username" bson:"username"`
	UUID         string `json:"uuid" bson:"uuid"`
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Surname      string `json:"surname,omitempty" bson:"surname,omitempty"`
	Email        string `json:"email,omitempty" bson:"email,omitempty"`
	PersonalCode string `json:"personal_code,omitempty" bson:"personal_code,omitempty"`
}
