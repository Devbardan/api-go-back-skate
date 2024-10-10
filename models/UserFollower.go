package models

import (
	"time"
)

type UserFollower struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserID     uint //// USUARIO A QUIEN SIGUEN.
	FollowerID uint /// USUARIO QUE SIGUE.
	User       User `gorm:"foreignKey:UserID"`
	Follower   User `gorm:"foreignKey:FollowerID"`
}

func (UserFollower) TableName() string {
	return "follower_user"
}
