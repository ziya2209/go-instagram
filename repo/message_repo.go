package repo

import (
	"instagram/models"

	"gorm.io/gorm"
)

type MessageRepo interface {
	Create(message *models.Message) error
	GetById(id int) (*models.Message, error)
	Update(message *models.Message) error
	Delete(id int) error
}

type GormMessageRepo struct {
	DB *gorm.DB
}

func (r *GormMessageRepo) Create(message *models.Message) error {
	return r.DB.Create(message).Error
}

func (r *GormMessageRepo) GetById(id int) (*models.Message, error) {
	var message models.Message
	if err := r.DB.First(&message, id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *GormMessageRepo) Update(message *models.Message) error {
	return r.DB.Save(message).Error
}

func (r *GormMessageRepo) Delete(id int) error {
	return r.DB.Delete(&models.Message{}, id).Error
}
