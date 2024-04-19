package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string `gorm:"type:varchar(100);not null";json:"Nombre"`
	NikName        string `gorm:"type:varchar(60);not null;unique_index";json:"NikName"`
	Email          string `gorm:"not null;unique_index";json:"Correo"`
	Password       string `gorm:json:"Contraseña"`
	ProfilePicture string
	Triks          []Trik `gorm:"foreignKey:UserID"json:"Trucos"`
}
