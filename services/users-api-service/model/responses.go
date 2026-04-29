package model

type ErrorResponse struct {
	Error    string `json:"error"`
	Username string `json:"username,omitempty"`
}

type UserUUIDResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
}
