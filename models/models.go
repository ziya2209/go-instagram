package models

import "instagram/database"

func init() {
	database.RegisterModels(
		&User{},
		&Comment{},
		&Like{},
		&Follow{},
		&Message{},
		&Post{},
	)
}