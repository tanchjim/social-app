package service

import (
	"context"

	"github.com/yourorg/social-app/internal/model"
)

type ContentService interface {
	Create(ctx context.Context, userID uint, content *model.Content) error
	GetByID(ctx context.Context, id uint) (*model.Content, error)
	List(ctx context.Context, page, pageSize int) ([]*model.Content, int64, error)
	GetReviewResult(ctx context.Context, contentID, userID uint) (*model.Content, error)
	SoftDelete(ctx context.Context, contentID, userID uint) error
}
