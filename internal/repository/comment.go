package repository

import (
	"context"

	"github.com/yourorg/social-app/internal/model"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	ListByContentID(ctx context.Context, contentID uint, offset, limit int) ([]*model.Comment, int64, error)
}
