package models

import "gorm.io/gorm"

type Webm struct {
	gorm.Model
	ID         uint32 `json:"id"`
	VideoId    uint32 `json:"video_id"`
	SubtitleId uint32 `json:"subtitle_id"`
	FilePath   string `json:"file_path"`
}
