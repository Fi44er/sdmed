package cart_adapters

import (
	"context"
)

type ICartUsecase interface {
	MoveByUserID(ctx context.Context, fromUserID string, toUserID string) error
}

type AuthMigrator struct {
	cartUsecase ICartUsecase
}

func NewAuthMigrator(cartUsecase ICartUsecase) *AuthMigrator {
	return &AuthMigrator{
		cartUsecase: cartUsecase,
	}
}

func (a *AuthMigrator) Migrate(ctx context.Context, fromUserID, toUserID string) error {
	return a.cartUsecase.MoveByUserID(ctx, fromUserID, toUserID)
}
