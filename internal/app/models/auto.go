package models

import "gorm.io/gorm"

//Auto models...omitempty
type Auto struct {
	gorm.Model
	ID       uint   `gorm:"primarykey" json:"id"`
	Mark     string `json:"mark"`
	MaxSpeed int64  `json:"max_speed"`
	Distance int64  `json:"distance"`
	Handler  string `json:"handler"`
	Stock    string `json:"stock"`
}

type AutoApi struct {
	MaxSpeed int64  `json:"max_speed"`
	Distance int64  `json:"distance"`
	Handler  string `json:"handler"`
	Stock    string `json:"stock"`
}
