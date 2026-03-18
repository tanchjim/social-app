package service

import (
	"context"

	"github.com/yourorg/social-app/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, username, password, nickname string) (*model.User, error)
	Login(ctx context.Context, username, password string) (accessToken, refreshToken string, user *model.User, err error)
	Refresh(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error)
	Logout(ctx context.Context, userID uint) error
}
