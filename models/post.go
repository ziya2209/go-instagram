package models

type Post struct {
	Id        int    `gorm:"primaryKey;column:id"`
	UserId    int    `gorm:"column:user_id"`
	Url       string `gorm:"column:url"`
	Caption   string `gorm:"column:caption"`
	CreatedAt string `gorm:"column:created_at"`
	UpdatedAt string `gorm:"column:updated_at"`
}
