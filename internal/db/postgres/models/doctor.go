package models

type Doctor struct {
	ID       int    `json:"id" gorm:"unique"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}
