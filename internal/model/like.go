package model

import (
	"time"
)

type Like struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ContentID uint      `gorm:"uniqueIndex:idx_content_user;not null;index" json:"content_id"`
	UserID    uint      `gorm:"uniqueIndex:idx_content_user;not null;index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (Like) TableName() string {
	return "likes"
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	TokenHash string    `gorm:"uniqueIndex;size:64;not null" json:"-"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
