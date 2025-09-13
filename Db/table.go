package Db

import "gorm.io/gorm"

// Model URL
type Url struct {
	gorm.Model
	ShortCode   string `gorm:"uniqueIndex;size:50"`
	OriginalUrl string
	VisitCount  int
}
