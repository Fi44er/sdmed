package session

import (
	"time"
)

type Session struct {
	ID             string
	Data           map[string]any
	CreatedAt      time.Time
	LastActivityAt time.Time
}

func newSession() *Session {
	return &Session{
		ID:             generateSessionId(),
		Data:           make(map[string]any),
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
	}
}

func (s *Session) Get(key string) any {
	return s.Data[key]
}

func (s *Session) Put(key string, value any) {
	s.Data[key] = value
}

func (s *Session) Delete(key string) {
	delete(s.Data, key)
}
