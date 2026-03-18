package service

import (
	"context"

	"github.com/yourorg/social-app/internal/model"
)

type UserService interface {
	GetByID(ctx context.Context, id uint) (*model.User, error)
	Update(ctx context.Context, userID uint, nickname, bio, avatar string) error
	GetMyContents(ctx context.Context, userID uint, status string, page, pageSize int) ([]*model.Content, int64, error)
}
