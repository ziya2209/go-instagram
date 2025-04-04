package repo

import (
	"instagram/models"

	"gorm.io/gorm"
)

type LikeRepo interface {
	Create(like *models.Like) error
	GetById(id int) (*models.Like, error)
	Delete(id int) error
}

type GormLikeRepo struct {
	DB *gorm.DB
}

func (r *GormLikeRepo) Create(like *models.Like) error {
	return r.DB.Create(like).Error
}

func (r *GormLikeRepo) GetById(id int) (*models.Like, error) {
	var like models.Like
	if err := r.DB.First(&like, id).Error; err != nil {
		return nil, err
	}
	return &like, nil
}

func (r *GormLikeRepo) Delete(id int) error {
	return r.DB.Delete(&models.Like{}, id).Error
}
