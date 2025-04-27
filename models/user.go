package models

type User struct {
	Id           int    `gorm:"primaryKey;column:id"`
	Username     string `gorm:"column:username"`
	Age          int    `gorm:"column:age"`
	Email        string `gorm:"column:email"`
	PasswordHash string `gorm:"column:password_hash"`
	Bio          string `gorm:"column:bio"`
}
