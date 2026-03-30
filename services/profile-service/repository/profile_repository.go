package repository

import "context"

type Profile struct {
	UUID         string
	Name         string
	Surname      string
	Email        string
	PersonalCode string
}

type ProfileRepository interface {
	GetByUUID(ctx context.Context, uuid string) (Profile, bool, error)
}