package service

import (
	"context"

	"github.com/yourorg/social-app/internal/model"
)

type CommentService interface {
	Create(ctx context.Context, contentID, userID uint, content string) (*model.Comment, error)
	List(ctx context.Context, contentID uint, page, pageSize int) ([]*model.Comment, int64, error)
}
