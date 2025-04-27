package models

type Follow struct {
	Id         int    `gorm:"primaryKey;column:id"`
	FollowerId int    `gorm:"column:follower_id"`
	FollowedId int    `gorm:"column:followed_id"`
	CreatedAt  string `gorm:"column:created_at"`
}
