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
	fmt.Printf("GetByPostId called with postId: %d\n", postId)

	// Get all comments for this post (root level comments) with usernames
	rootComments, err := r.getByTypeAndIdsWithUsernames("post", []int{postId})
	if err != nil {
		fmt.Printf("Error getting root comments: %v\n", err)
		return nil, err
	}

	fmt.Printf("Found %d root comments\n", len(rootComments))

	// If no root comments, return empty slice
	if len(rootComments) == 0 {
		return []CommentTree{}, nil
	}

	// Create root nodes
	var roots []CommentTree
	for _, com := range rootComments {
		fmt.Printf("Creating root node for comment ID %d, username: %s\n", com.Id, com.Username)
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
		var nextLevelComments []CommentWithUser

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

		// Get all replies to the current level comments with usernames
		replies, err := r.getByTypeAndIdsWithUsernames("comment", currentIds)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Found %d replies for current level\n", len(replies))

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

	fmt.Printf("Returning %d comment trees\n", len(roots))
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
	Comment  CommentWithUser `json:"comment"`
	Children []CommentTree   `json:"children"`
}

type CommentWithUser struct {
	Id         int    `json:"id"`
	ParentId   int    `json:"parent_id"`
	ParentType string `json:"parent_type"`
	UserId     int    `json:"user_id"`
	Username   string `json:"username"`
	Message    string `json:"message"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func (r *gormCommentRepo) getByTypeAndIdsWithUsernames(parentType string, ids []int) ([]CommentWithUser, error) {
	var results []CommentWithUser

	query := `
		SELECT c.id, c.parent_id, c.parent_type, c.user_id, u.username, c.message, c.created_at, c.updated_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.parent_type = ? AND c.parent_id IN ?
		ORDER BY c.created_at ASC
	`

	if err := r.db.Raw(query, parentType, ids).Scan(&results).Error; err != nil {
		return nil, err
	}

	// Debug log to see what we're getting
	fmt.Printf("GetByTypeAndIdsWithUsernames returning %d results\n", len(results))
	for i, result := range results {
		fmt.Printf("Result %d: ID=%d, Username=%s, Message=%s\n", i, result.Id, result.Username, result.Message)
	}

	return results, nil
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
