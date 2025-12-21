package sessionstore

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Fi44er/sdmed/pkg/session"
	"github.com/redis/go-redis/v9"
)

type RedisSessionStore struct {
	client *redis.Client
}

func NewRedisSessionStore(client *redis.Client) *RedisSessionStore {
	return &RedisSessionStore{
		client: client,
	}
}

func (s *RedisSessionStore) Read(id string) (*session.Session, error) {
	ctx := context.Background()
	data, err := s.client.Get(ctx, "session:"+id).Bytes()
	if err == redis.Nil {
		return nil, nil // Session not found
	} else if err != nil {
		return nil, err
	}

	var session session.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *RedisSessionStore) Write(session *session.Session) error {
	ctx := context.Background()
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = s.client.Set(ctx, "session:"+session.ID, data, 24*time.Hour).Err()
	return err
}

func (s *RedisSessionStore) Destroy(id string) error {
	ctx := context.Background()
	return s.client.Del(ctx, "session:"+id).Err()
}

func (s *RedisSessionStore) Gc(idleExpiration, absoluteExpiration time.Duration) error {
	ctx := context.Background()
	iter := s.client.Scan(ctx, 0, "session:*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		data, err := s.client.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var session session.Session
		if err := json.Unmarshal(data, &session); err != nil {
			continue
		}

		if time.Since(session.LastActivityAt) > idleExpiration ||
			time.Since(session.CreatedAt) > absoluteExpiration {
			s.client.Del(ctx, key)
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}

func (s *RedisSessionStore) Close() error {
	return s.client.Close()
}
