package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Fullname string `json:"Fullname"`
	Age      string `json:"Age"`
	Location string `json:"Location"`
	Email    string `gorm:"unique" json:"Email"`
	Zipcode  string `json:"Zipcode"`
}
