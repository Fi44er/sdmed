package session

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

const SessionContextKey = "session_id"

type IStore interface {
	Read(id string) (*Session, error)
	Write(session *Session) error
	Destroy(id string) error
	Gc(idleExpiration, absoluteExpiration time.Duration) error
}

type HttpContext interface {
	Context() context.Context
	SetContext(ctx context.Context)
	GetCookie(name string) string
	SetCookie(name, value string, maxAge int)
	SetHeader(key, value string)
}

type SessionManager struct {
	store              IStore
	idleExpiration     time.Duration
	absoluteExpiration time.Duration
	cookieName         string
}

func NewSessionManager(
	store IStore,
	gcInterval,
	idleExpiration,
	absoluteExpiration time.Duration,
	cookieName string,
) *SessionManager {

	m := &SessionManager{
		store:              store,
		idleExpiration:     idleExpiration,
		absoluteExpiration: absoluteExpiration,
		cookieName:         cookieName,
	}

	go m.gc(gcInterval)

	return m
}

func (m *SessionManager) gc(d time.Duration) {
	ticker := time.NewTicker(d)

	for range ticker.C {
		m.store.Gc(m.idleExpiration, m.absoluteExpiration)
	}
}

func (m *SessionManager) validate(session *Session) bool {
	if time.Since(session.CreatedAt) > m.absoluteExpiration ||
		time.Since(session.LastActivityAt) > m.idleExpiration {
		err := m.store.Destroy(session.ID)
		if err != nil {
			panic(err)
		}
		return false
	}

	return true
}

func (m *SessionManager) Start(ctx HttpContext) (context.Context, *Session) {
	var session *Session
	var err error

	cookie := ctx.GetCookie(m.cookieName)

	session, err = m.store.Read(cookie)
	if err != nil {
		log.Printf("Failed to read session from store: %v", err)
	}

	if session == nil || !m.validate(session) {
		session = newSession()
	}

	return context.WithValue(ctx.Context(), SessionContextKey, session), session
}

func (m *SessionManager) Save(ctx HttpContext, session *Session) error {
	session.LastActivityAt = time.Now()

	err := m.store.Write(session)
	if err != nil {
		return err
	}

	ctx.SetCookie(m.cookieName, session.ID, int(m.absoluteExpiration.Seconds()))

	return nil
}

func (m *SessionManager) migrate(session *Session) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	err := m.store.Destroy(session.ID)
	if err != nil {
		return err
	}

	session.ID = generateSessionId()

	return nil
}

func FromContext(ctx context.Context) *Session {
	if s, ok := ctx.Value(SessionContextKey).(*Session); ok {
		return s
	}
	return nil
}

func FromFiberContext(ctx *fiber.Ctx) *Session {
	if s, ok := ctx.UserContext().Value(SessionContextKey).(*Session); ok {
		return s
	}
	return nil
}
