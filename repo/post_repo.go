package repo

import (
	"instagram/models"

	"gorm.io/gorm"
)

type PostRepo interface {
	Create(post *models.Post) error
	GetById(id int) (*models.Post, error)
	Update(post *models.Post) error
	Delete(id int) error
}

type GormPostRepo struct {
	DB *gorm.DB
}

func (r *GormPostRepo) Create(post *models.Post) error {
	return r.DB.Create(post).Error
}

func (r *GormPostRepo) GetById(id int) (*models.Post, error) {
	var post models.Post
	if err := r.DB.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *GormPostRepo) Update(post *models.Post) error {
	return r.DB.Save(post).Error
}

func (r *GormPostRepo) Delete(id int) error {
	return r.DB.Delete(&models.Post{}, id).Error
}
