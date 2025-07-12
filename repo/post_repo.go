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
	LikePost(postId int, userId int) error
	GetPostsByUserIds(userIds []int, page int, limit int) ([]models.Post, int, error)
	GetPostsByUserId(userId int, page int, limit int) ([]models.Post, int, error)
	GetLikeCount(postId int) (int, error)
	GetCommentCount(postId int) (int, error)
	IsPostOwner(postId int, userId int) (bool, error)
}

type GormPostRepo struct {
	DB *gorm.DB
}

func NewGormPostRepo(db *gorm.DB) *GormPostRepo {
	return &GormPostRepo{DB: db}
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
func (r *GormPostRepo) LikePost(postId int, userId int) error {
	// Assuming you have a Like model and a Likes table
	like := models.Like{
		PostId: postId,
		UserId: userId,
	}
	return r.DB.Create(&like).Error
}

func (r *GormPostRepo) Update(post *models.Post) error {
	return r.DB.Save(post).Error
}

func (r *GormPostRepo) Delete(id int) error {
	return r.DB.Delete(&models.Post{}, id).Error
}

// GetPostsByUserId retrieves posts from a specific user with pagination, ordered by creation date (newest first)
func (r *GormPostRepo) GetPostsByUserId(userId int, page int, limit int) ([]models.Post, int, error) {
	var posts []models.Post
	var totalCount int64

	// Calculate offset
	offset := (page - 1) * limit

	// Get total count for pagination
	if err := r.DB.Model(&models.Post{}).Where("user_id = ?", userId).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get posts with pagination, ordered by creation date descending (newest first)
	if err := r.DB.Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, int(totalCount), nil
}

// GetPostsByUserIds retrieves posts from specific users with pagination, ordered by creation date (newest first)
func (r *GormPostRepo) GetPostsByUserIds(userIds []int, page int, limit int) ([]models.Post, int, error) {
	var posts []models.Post
	var totalCount int64

	// Calculate offset
	offset := (page - 1) * limit

	// Get total count for pagination
	if err := r.DB.Model(&models.Post{}).Where("user_id IN ?", userIds).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get posts with pagination, ordered by creation date descending (newest first)
	if err := r.DB.Where("user_id IN ?", userIds).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, int(totalCount), nil
}

// GetLikeCount returns the number of likes for a specific post
func (r *GormPostRepo) GetLikeCount(postId int) (int, error) {
	var count int64
	if err := r.DB.Model(&models.Like{}).Where("post_id = ?", postId).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetCommentCount returns the number of comments for a specific post
func (r *GormPostRepo) GetCommentCount(postId int) (int, error) {
	var count int64
	if err := r.DB.Model(&models.Comment{}).Where("parent_type = ? AND parent_id = ?", "post", postId).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// IsPostOwner checks if a user owns a specific post
func (r *GormPostRepo) IsPostOwner(postId int, userId int) (bool, error) {
	var count int64
	if err := r.DB.Model(&models.Post{}).Where("id = ? AND user_id = ?", postId, userId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
