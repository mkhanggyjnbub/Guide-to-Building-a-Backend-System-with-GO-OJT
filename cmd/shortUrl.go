package main

import (
	"baitapweek1/Db"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// Sinh random string 6 ký tự
func RandomString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Sinh short code không trùng database
func GenerateUniqueCode(db *gorm.DB, length int) string {
	for {
		code := RandomString(length)
		var url Db.Url
		if err := db.Where("short_code = ?", code).First(&url).Error; err != nil {
			// Không tìm thấy → unique
			return code
		}
		// Nếu trùng → sinh lại
	}
}
