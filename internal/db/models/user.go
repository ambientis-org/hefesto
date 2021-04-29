package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" validate:"required"`
	Name 	 string `json:"name" validate:"required"`
	Email	 string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}