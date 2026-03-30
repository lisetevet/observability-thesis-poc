package repository

import "context"

type InMemoryUserRepository struct {
	users map[string]string
}

func NewInMemoryUserRepository(seed map[string]string) *InMemoryUserRepository {
	// shallow copy to avoid accidental outside mutation
	cp := make(map[string]string, len(seed))
	for k, v := range seed {
		cp[k] = v
	}
	return &InMemoryUserRepository{users: cp}
}

func (r *InMemoryUserRepository) GetUUIDByUsername(ctx context.Context, username string) (string, bool, error) {
	uuid, ok := r.users[username]
	return uuid, ok, nil
}