package models

import "gorm.io/gorm"

type Subtitle struct {
	gorm.Model
	ID        int     `gorm:"primaryKey"`
	VideoID   int     `json:"video_id"`
	Text      string  `json:"text"`
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
}
