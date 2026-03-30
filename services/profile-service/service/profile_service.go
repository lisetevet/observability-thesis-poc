package service

import (
	"context"

	"profile-service/repository"
)

type ProfileService struct {
	repo repository.ProfileRepository
}

func NewProfileService(repo repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) GetProfile(ctx context.Context, uuid string) (repository.Profile, bool, error) {
	return s.repo.GetByUUID(ctx, uuid)
}