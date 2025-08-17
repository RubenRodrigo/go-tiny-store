package models

type Category struct {
	Base
	Name     string    `json:"name" gorm:"not null;size:100"`
	Products []Product `json:"products"`
}
