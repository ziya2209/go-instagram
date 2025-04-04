package repo

import (
	"instagram/models"

	"gorm.io/gorm"
)

type CommentRepo interface {
	Create(comment *models.Comment) error
	GetById(id int) (*models.Comment, error)
	Update(comment *models.Comment) error
	Delete(id int) error
}

type GormCommentRepo struct {
	DB *gorm.DB
}

func (r *GormCommentRepo) Create(comment *models.Comment) error {
	return r.DB.Create(comment).Error
}

func (r *GormCommentRepo) GetById(id int) (*models.Comment, error) {
	var comment models.Comment
	if err := r.DB.Where("id = ?", id).First(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *GormCommentRepo) Update(comment *models.Comment) error {
	return r.DB.Save(comment).Error
}

func (r *GormCommentRepo) Delete(id int) error {
	return r.DB.Delete(&models.Comment{}, id).Error
}
