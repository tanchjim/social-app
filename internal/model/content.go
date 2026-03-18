package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ContentType string

const (
	ContentTypeImage ContentType = "image"
	ContentTypeVideo ContentType = "video"
)

type ContentStatus string

const (
	StatusDraft     ContentStatus = "draft"
	StatusReviewing ContentStatus = "reviewing"
	StatusPublished ContentStatus = "published"
	StatusRejected  ContentStatus = "rejected"
	StatusDeleted   ContentStatus = "deleted"
)

type Content struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	Type         ContentType    `gorm:"size:20;not null" json:"type"`
	Title        string         `gorm:"size:200" json:"title"`
	Description  string         `json:"description"`
	MediaURL     string         `gorm:"size:500;not null" json:"media_url"`
	CoverURL     string         `gorm:"size:500" json:"cover_url"`
	Status       ContentStatus  `gorm:"size:20;default:draft;index" json:"status"`
	RejectReason string         `json:"reject_reason,omitempty"`
	LikeCount    int            `gorm:"default:0" json:"like_count"`
	CommentCount int            `gorm:"default:0" json:"comment_count"`
	Tags         pq.StringArray `gorm:"type:text[]" json:"tags"`
	DeletedAt    *time.Time     `json:"deleted_at,omitempty"`
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`

	// Relations
	User     User      `gorm:"foreignKey:UserID" json:"author,omitempty"`
	Comments []Comment `gorm:"foreignKey:ContentID" json:"-"`
	Likes    []Like    `gorm:"foreignKey:ContentID" json:"-"`
}

func (Content) TableName() string {
	return "contents"
}
