package models

type User struct {
	Base
	Email     string `json:"email" gorm:"unique;not null"`
	Username  string `json:"username" gorm:"not null"`
	Password  string `json:"-" gorm:"not null"` // Password is never returned in JSON
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
