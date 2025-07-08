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
	var roots []CommentTree
	comments, err := r.GetByTypeAndIds("post", []int{postId})
	if err != nil {
		return nil, err
	}
	for _, com := range comments {
		roots = append(roots, CommentTree{
			Comment: com,
		})
	}

	m := map[int][]CommentTree{}

	for len(comments) > 0 {
		ids := []int{}
		for _, v := range comments {
			ids = append(ids, v.Id)
			m[v.Id] = append(m[v.Id], CommentTree{
				Comment: v,
			})

		}
		comments, err = r.GetByTypeAndIds("comments", ids)
		if err != nil {
			return nil, err
		}
	}

	for _, root := range roots {
		buildCommentTree(&root, m)
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
