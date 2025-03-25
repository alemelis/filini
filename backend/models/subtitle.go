package models

import "gorm.io/gorm"

type Subtitle struct {
	gorm.Model
	ID        uint32  `gorm:"primaryKey"`
	VideoID   uint32  `json:"video_id"`
	Text      string  `json:"text"`
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
}
