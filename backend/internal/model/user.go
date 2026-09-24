package model

import "time"

// User 用户实体。
type User struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Phone        string    `gorm:"size:20;uniqueIndex;not null" json:"phone"`
	PasswordHash string    `gorm:"size:100;not null" json:"-"`
	Name         string    `gorm:"size:50;not null;default:''" json:"name"`
	Avatar       string    `gorm:"size:255;not null;default:''" json:"avatar"`
	Role         string    `gorm:"size:30;not null;default:worker" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 指定表名。
func (User) TableName() string { return "users" }
