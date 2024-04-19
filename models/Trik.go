package models

import "gorm.io/gorm"

type Trik struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null;unique_index"`
	Description string `gorm:"not null"`
	TrikImage   string
	Level       int  `gorm:"default:1"`
	Done        bool `gorm:"default:false"`
	UserID      uint
}
