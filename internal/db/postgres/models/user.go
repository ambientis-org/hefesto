package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username  string 	 `json:"username" validate:"required" gorm:"unique"'`
	Email	  string 	 `json:"email" validate:"required,email"  gorm:"unique"`
	Name 	  string 	 `json:"name" validate:"required"`
	Lastname  string 	 `json:"lastname" validate:"required"`
	Password  string 	 `json:"password" validate:"required"`
	Birthdate *time.Time `json:"birthdate" validate:"required"`
	Genre	  string 	 `json:"genre" validate:"required"`
}