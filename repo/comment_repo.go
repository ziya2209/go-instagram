package repo

import (
	"fmt"
	"instagram/models"

	"gorm.io/gorm"
)

type CommentRepo interface {
	Create(comment *models.Comment) error
	GetById(id int) (*models.Comment, error)
	Update(comment *models.Comment) error
	Delete(id int) error
	GetByPostId(postId int) ([]CommentTree, error)
}

type gormCommentRepo struct {
	db *gorm.DB
}

func NewGormCommentRepo(db *gorm.DB) *gormCommentRepo {
	return &gormCommentRepo{db: db}
}

func (r *gormCommentRepo) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *gormCommentRepo) GetById(id int) (*models.Comment, error) {
	var comment models.Comment
	if err := r.db.Where("id = ?", id).First(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *gormCommentRepo) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *gormCommentRepo) Delete(id int) error {
	return r.db.Delete(&models.Comment{}, id).Error
}

func (r *gormCommentRepo) GetByPostId(postId int) ([]CommentTree, error) {
	// Get all comments for this post (root level comments)
	rootComments, err := r.GetByTypeAndIds("post", []int{postId})
	if err != nil {
		return nil, err
	}

	// If no root comments, return empty slice
	if len(rootComments) == 0 {
		return []CommentTree{}, nil
	}

	// Create root nodes
	var roots []CommentTree
	for _, com := range rootComments {
		roots = append(roots, CommentTree{
			Comment: com,
		})
	}

	// Build a map of parent_id -> children for all nested comments
	childrenMap := make(map[int][]CommentTree)

	// Get all nested comments (replies) in a breadth-first manner
	currentLevelComments := rootComments
	processedIds := make(map[int]bool) // Prevent infinite loops

	for len(currentLevelComments) > 0 {
		var nextLevelComments []models.Comment

		// Get comment IDs for the current level
		var currentIds []int
		for _, comment := range currentLevelComments {
			// Skip if we've already processed this comment (prevent cycles)
			if processedIds[comment.Id] {
				continue
			}
			currentIds = append(currentIds, comment.Id)
			processedIds[comment.Id] = true
		}

		// If no new IDs to process, break to prevent infinite loop
		if len(currentIds) == 0 {
			break
		}

		// Get all replies to the current level comments
		replies, err := r.GetByTypeAndIds("comment", currentIds)
		if err != nil {
			return nil, err
		}

		// Group replies by their parent_id
		for _, reply := range replies {
			childrenMap[reply.ParentId] = append(childrenMap[reply.ParentId], CommentTree{
				Comment: reply,
			})
			nextLevelComments = append(nextLevelComments, reply)
		}

		currentLevelComments = nextLevelComments
	}

	// Build the complete tree structure
	for i := range roots {
		buildCommentTree(&roots[i], childrenMap)
	}

	return roots, nil
}

// buildCommentTree recursively builds the comment tree structure
func buildCommentTree(node *CommentTree, childrenMap map[int][]CommentTree) {
	children, exists := childrenMap[node.Comment.Id]
	if !exists {
		return
	}

	node.Children = children
	for i := range node.Children {
		buildCommentTree(&node.Children[i], childrenMap)
	}
}

type CommentTree struct {
	Comment  models.Comment
	Children []CommentTree
}

func (ct CommentTree) PrintIndent() {
	ct.printWithIndent(0)
}

func (ct CommentTree) printWithIndent(level int) {
	// Print current comment with indentation
	indent := ""
	for i := 0; i < level; i++ {
		indent += "  " // Two spaces per level
	}
	fmt.Printf("%s- %s (ID: %d, User: %d)\n", indent, ct.Comment.Message, ct.Comment.Id, ct.Comment.UserId)

	// Recursively print children with increased indentation
	for _, child := range ct.Children {
		child.printWithIndent(level + 1)
	}
}

func (r *gormCommentRepo) GetByTypeAndIds(parentType string, ids []int) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.Where("parent_type = ? AND parent_id IN ?", parentType, ids).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
