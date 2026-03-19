package service

import "context"

type LikeService interface {
	Toggle(ctx context.Context, contentID, userID uint, action string) (isLiked bool, likeCount int64, err error)
}
