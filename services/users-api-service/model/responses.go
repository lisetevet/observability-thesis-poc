package model

type ErrorResponse struct {
	Error    string `json:"error"`
	Username string `json:"username,omitempty"`
}

type UserUUIDResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
}

type ProfileSeedResponse struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	PersonalCode string `json:"personal_code"`
}