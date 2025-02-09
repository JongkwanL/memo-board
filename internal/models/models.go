package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"type:varchar(255);uniqueIndex"`
	Password  string `gorm:"type:varchar(255)"`
	Email     string `gorm:"type:varchar(255);uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Post struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Content   string `gorm:"type:text"`
	AuthorID  uint
	Author    User `gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
