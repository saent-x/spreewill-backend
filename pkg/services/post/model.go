package post

import (
	"gorm.io/gorm"
)

type (
	PostDTO struct {
		gorm.Model
		Images   []string `json:"images"`
		Caption  string   `json:"caption"`
		Likes    uint     `json:"like"`
		Dislikes uint     `json:"dislikes"`
		Views    uint     `json:"views"`
	}
)
