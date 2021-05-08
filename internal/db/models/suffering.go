package models

import "gorm.io/gorm"

type Suffering struct {
	gorm.Model
	Name	 	string
	Severity 	uint8
	RequiresMed bool
}