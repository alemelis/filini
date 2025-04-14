package models

import "gorm.io/gorm"

type Clip struct {
	gorm.Model
	Quote    string `json:"quote"`
	WebmPath string `json:"webm_path"`
}
