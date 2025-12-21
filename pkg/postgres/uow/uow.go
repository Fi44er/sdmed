package uow

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

type uowKey struct{}

var (
	ErrTxAlreadyStarted   = errors.New("transaction already started")
	ErrTxNotStarted       = errors.New("no transaction started")
	ErrRepositoryNotFound = errors.New("repository not registered")
)

type RepositoryFactory func(tx *gorm.DB) (any, error)

type Uow interface {
	RegisterRepository(name string, factory RepositoryFactory)
	GetRepository(ctx context.Context, name string) (any, error)
	Do(ctx context.Context, fn func(ctx context.Context) error) error

	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type uow struct {
	db           *gorm.DB
	repositories map[string]RepositoryFactory
	mu           sync.RWMutex
}

func New(db *gorm.DB) Uow {
	return &uow{
		db:           db,
		repositories: make(map[string]RepositoryFactory),
	}
}

func (u *uow) RegisterRepository(name string, factory RepositoryFactory) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.repositories[name] = factory
}

func (u *uow) GetRepository(ctx context.Context, name string) (any, error) {
	u.mu.RLock()
	factory, exists := u.repositories[name]
	u.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrRepositoryNotFound, name)
	}

	tx, ok := ctx.Value(uowKey{}).(*gorm.DB)
	if !ok {
		return factory(u.db.WithContext(ctx))
	}

	return factory(tx.WithContext(ctx))
}

func (u *uow) Begin(ctx context.Context) (context.Context, error) {
	if _, ok := ctx.Value(uowKey{}).(*gorm.DB); ok {
		return ctx, ErrTxAlreadyStarted
	}

	tx := u.db.Begin()
	if tx.Error != nil {
		return ctx, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	return context.WithValue(ctx, uowKey{}, tx), nil
}

func (u *uow) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(uowKey{}).(*gorm.DB)
	if !ok {
		return ErrTxNotStarted
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (u *uow) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(uowKey{}).(*gorm.DB)
	if !ok {
		return ErrTxNotStarted
	}

	if err := tx.Rollback().Error; err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}

func (u *uow) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	txCtx, err := u.Begin(ctx)
	if err != nil {
		if errors.Is(err, ErrTxAlreadyStarted) {
			return fn(ctx)
		}
		return err
	}

	var fnErr error
	defer func() {
		if p := recover(); p != nil {
			_ = u.Rollback(txCtx)
			panic(p)
		}

		if fnErr != nil {
			_ = u.Rollback(txCtx)
		}
	}()

	fnErr = fn(txCtx)

	if fnErr != nil {
		return fnErr
	}

	return u.Commit(txCtx)
}
