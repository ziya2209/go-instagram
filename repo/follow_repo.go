package repo

import (
	"instagram/models"

	"gorm.io/gorm"
)

type FollowRepo interface {
	Create(follow *models.Follow) error
	GetById(id int) (*models.Follow, error)
	Delete(id int) error
}

type GormFollowRepo struct {
	DB *gorm.DB
}

func (r *GormFollowRepo) Create(follow *models.Follow) error {
	return r.DB.Create(follow).Error
}

func (r *GormFollowRepo) GetById(id int) (*models.Follow, error) {
	var follow models.Follow
	if err := r.DB.First(&follow, id).Error; err != nil {
		return nil, err
	}
	return &follow, nil
}

func (r *GormFollowRepo) Delete(id int) error {
	return r.DB.Delete(&models.Follow{}, id).Error
}
