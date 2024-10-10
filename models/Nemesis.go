package models

import "gorm.io/gorm"

type Nemesis struct {
	gorm.Model
	UserID uint
	TrikID uint
	User   User `gorm:"foreignKey:UserID"`
	Trik   Trik `gorm:"foreignKey:TrikID"`
}
