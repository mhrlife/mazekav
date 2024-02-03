package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model

	DisplayName string `json:"display_name"`
	PhoneNumber string `json:"phone_number"`
}

func (u User) Table() string {
	return "users"
}
