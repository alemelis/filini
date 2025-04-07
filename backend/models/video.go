package models

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	ID    uint32 `json:"id"`
	Title string `json:"title"`
}
