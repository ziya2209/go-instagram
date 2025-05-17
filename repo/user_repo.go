package repo

import (
	"errors"
	"instagram/models"

	"gorm.io/gorm" // object relation mapper
)

type UserRepo interface {
	Create(user *models.User) error
	GetById(id int) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
	GetAll() ([]*models.User, error)
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &gormUserRepo{DB: db}
}

type gormUserRepo struct {
	DB *gorm.DB
}

func (r *gormUserRepo) Create(user *models.User) error {
	// Check if email already exists
	var existingUser models.User
	if err := r.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("email already exists")
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	return r.DB.Create(user).Error
}

func (r *gormUserRepo) GetById(id int) (*models.User, error) {
	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepo) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
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
func (r *gormUserRepo) GetAll() ([]*models.User, error) {
	var users []*models.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}




