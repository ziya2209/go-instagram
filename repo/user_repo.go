package repo

import (
	"instagram/models"
)

type UserRepo interface {
	Create(user *models.User) error
	GetById(id int) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
	GetAll() ([]*models.User, error)
}
