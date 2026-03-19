package repository

import (
	"context"

	"social-app/internal/model"
)

type LikeRepository interface {
	Create(ctx context.Context, like *model.Like) error
	Delete(ctx context.Context, contentID, userID uint) error
	GetByContentAndUser(ctx context.Context, contentID, userID uint) (*model.Like, error)
	CountByContentID(ctx context.Context, contentID uint) (int64, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, id uint) error
	RevokeAllByUserID(ctx context.Context, userID uint) error
}
