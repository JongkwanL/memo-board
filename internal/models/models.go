package models

import "time"

type Role uint8

const (
	ADMIN Role = iota
	USER
)

type User struct {
	ID         uint   `gorm:"primaryKey"`
	Username   string `gorm:"type:varchar(255);uniqueIndex"`
	Password   string `gorm:"type:varchar(255)"`
	Email      string `gorm:"type:varchar(255);uniqueIndex"`
	Role       Role   `gorm:"type:tinyint(1);default:1"`
	IsApproved bool   `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
