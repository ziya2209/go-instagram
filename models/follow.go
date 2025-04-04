package models

type Follow struct {
	Id         int
	FollowerId int
	FollowedId int
	CreatedAt  string
}
