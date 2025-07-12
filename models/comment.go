package models

import "time"

type Comment struct {
	Id         int       `gorm:"primaryKey;column:id"`
	ParentId   int       `gorm:"column:parent_id"`
	ParentType string    `gorm:"column:parent_type"`
	UserId     int       `gorm:"column:user_id"`
	Message    string    `gorm:"column:message"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
