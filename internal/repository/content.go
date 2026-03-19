package repository

import (
	"context"
	"time"

	"social-app/internal/model"
)

type ContentRepository interface {
	Create(ctx context.Context, content *model.Content) error
	GetByID(ctx context.Context, id uint) (*model.Content, error)
	List(ctx context.Context, status string, offset, limit int) ([]*model.Content, int64, error)
	ListByUser(ctx context.Context, userID uint, status string, offset, limit int) ([]*model.Content, int64, error)
	Update(ctx context.Context, content *model.Content) error
	SoftDelete(ctx context.Context, id uint, deletedAt time.Time) error
}
