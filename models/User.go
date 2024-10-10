package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name              string `gorm:"type:varchar(100);not null" json:"Nombre"`
	NikName           string `gorm:"type:varchar(60);not null;unique_index" json:"NikName"`
	Email             string `gorm:"not null;unique_index" json:"Correo"`
	Password          string `json:"Contraseña"`
	ProfilePicture    string
	BannerOrDegradado string
	Triks             []Trik  `gorm:"foreignKey:UserID" json:"Trucos"`
	Nemesis           []Trik  `gorm:"foreignKey:UserID" json:"Nemesis"`
	Followers         []*User `gorm:"many2many:follower_user;"`
}
