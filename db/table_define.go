package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId   string `gorm:"unique"`
	Password string
	Nickname string
}

type Post struct {
	gorm.Model
	UserId           string
	Title            string
	Content          string
	CreatedTimestamp string
}
