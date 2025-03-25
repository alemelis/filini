package models

import "gorm.io/gorm"

type Gif struct {
	gorm.Model
	ID         string `json:"id"`
	VideoId    string `json:"video_id"`
	SubtitleId string `json:"subtitle_id"`
	FilePath   string `json:"file_path"`
}
