package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserId   string `gorm:"unique;not null;size:255"`
	Password string `gorm:"not null;size:255"`
	Nickname string `gorm:"uniqueIndex:idx_nickname_tag;size:20;not null"`
	Tag      uint64 `gorm:"uniqueIndex:idx_nickname_tag;size:20;not null"`
}

type Post struct {
	gorm.Model
	UserIndex        uint   `gorm:"not null;index;column:user_index"`
	Title            string `gorm:"not null;size:255"`
	Content          string `gorm:"not null;size:3000"`
	CreatedTimestamp string `gorm:"not null"`
	User             User   `gorm:"foreignKey:UserIndex;references:ID"`
}

type Comment struct {
	gorm.Model
	PostIndex          uint   `gorm:"not null;index"`
	UserIndex          uint   `gorm:"not null;index"`
	ParentCommentIndex *uint  `gorm:"index"` // [중요] 포인터(*uint)로 변경하여 NULL 허용 (최상위 댓글을 위해)
	Content            string `gorm:"not null;size:2000"`
	User               User   `gorm:"foreignKey:UserIndex;references:ID"`
	Post               Post   `gorm:"foreignKey:PostIndex;references:ID"`
}
