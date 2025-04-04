package models

type Comment struct {
	Id       int
	ParentId int
	// parenttype can have two values "post"or "comment"
	ParentType string
	UserId     int
	Message    string
	CreatedAt  string
	UpdatedAt  string
}
