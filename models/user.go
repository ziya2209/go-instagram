package models

type User struct {
	Id           int
	Username     string
	Age          int
	Email        string
	PasswordHash string
	Bio          string
}
