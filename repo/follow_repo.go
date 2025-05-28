package repo

import (
	"errors"
	"instagram/models"
	"strings"

	"gorm.io/gorm"
)

type FollowRepo interface {
	Create(follow *models.Follow) error
	GetById(id int) (*models.Follow, error)
	Delete(id int) error
	FollowUser(followerID int, followeeID int) error
	GetFollowers(userId int) ([]*models.Follow, error)
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

var ErrDuplicateFollow = errors.New("you are already following this user")

func (r *GormFollowRepo) FollowUser(followerID int, followedID int) error {
	follow := &models.Follow{
		FollowerId: followerID,
		FollowedId: followedID,
	}
	err := r.Create(follow)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") || strings.Contains(err.Error(), "Duplicate entry") {
			return ErrDuplicateFollow
		}
		return err
	}
	return nil
}

func (r *GormFollowRepo) GetFollowers(userId int) ([]*models.Follow, error) {
	var followers []*models.Follow
	if err := r.DB.Where("follower_id = ?", userId).Find(&followers).Error; err != nil {
		return nil, err
	}
	return followers, nil
}
