package repo

import (
	"instagram/models"

	"gorm.io/gorm" // object relation mapper
)

type UserRepo interface {
	Create(user *models.User) error
	GetById(id int) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &gormUserRepo{DB: db}
}

type gormUserRepo struct {
	DB *gorm.DB
}

func (r *gormUserRepo) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *gormUserRepo) GetById(id int) (*models.User, error) {
	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepo) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *gormUserRepo) Delete(id int) error {
	return r.DB.Delete(&models.User{}, id).Error
}
