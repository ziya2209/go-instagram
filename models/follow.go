package models

type Follow struct {
	Id         int    `gorm:"primaryKey;column:id"`
	FollowerId int    `gorm:"column:follower_id;uniqueIndex:idx_follower_followed"`
	FollowedId int    `gorm:"column:followed_id;uniqueIndex:idx_follower_followed"`
	CreatedAt  string `gorm:"column:created_at"`
}
