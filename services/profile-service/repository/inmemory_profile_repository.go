package repository

import "context"

type InMemoryProfileRepository struct {
	items map[string]Profile
}

func NewInMemoryProfileRepository(seed map[string]Profile) *InMemoryProfileRepository {
	cp := make(map[string]Profile, len(seed))
	for k, v := range seed {
		cp[k] = v
	}
	return &InMemoryProfileRepository{items: cp}
}

func (r *InMemoryProfileRepository) GetByUUID(ctx context.Context, uuid string) (Profile, bool, error) {
	p, ok := r.items[uuid]
	return p, ok, nil
}