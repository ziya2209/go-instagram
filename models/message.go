package models

type Message struct {
	Id         int    `gorm:"primaryKey;column:id"`
	SenderId   int    `gorm:"column:sender_id"`
	ReceiverId int    `gorm:"column:receiver_id"`
	Content    string `gorm:"column:content"`
	CreatedAt  string `gorm:"column:created_at"`
}
