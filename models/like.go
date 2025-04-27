package models

type Like struct {
	Id        int    `gorm:"primaryKey;column:id"`
	UserId    int    `gorm:"column:user_id"`
	PostId    int    `gorm:"column:post_id"`
	CreatedAt string `gorm:"column:created_at"`
}
