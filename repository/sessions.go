package repository

import (
	"time"

	"github.com/go-redis/redis"
)

type RepoSessions struct {
	client *redis.Client
}

func NewRepoSessions(client *redis.Client) *RepoSessions {
	return &RepoSessions{client: client}
}

func (r *RepoSessions) Create(id string, duration time.Duration) error {
	return r.client.Set(id, 0, duration).Err()
}

func (r *RepoSessions) Get(id string) (*int64, error) {
	val, err := r.client.Get(id).Int64()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &val, nil
}

func (r *RepoSessions) Delete(id string) error {
	return r.client.Del(id).Err()
}
