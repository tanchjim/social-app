package model

import (
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ContentID uint      `gorm:"index;not null" json:"content_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"author,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}
