package main

import (
	"net/http"
	"strings"

	"baitapweek1/Db"

	"github.com/gin-gonic/gin"
)

func main() {
	db := Db.InitDB()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Trang nhập URL
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})

	// Tạo short URL
	r.POST("/save", func(c *gin.Context) {
		longUrl := c.PostForm("longUrl")
		shortCode := c.PostForm("shortCode")

		// Nếu người dùng nhập URL không có http:// hoặc https:// → thêm https://
		if !strings.HasPrefix(longUrl, "http://") && !strings.HasPrefix(longUrl, "https://") {
			longUrl = "https://" + longUrl
		}

		// Kiểm tra URL đã tồn tại
		var existing Db.Url
		if err := db.Where("original_url = ?", longUrl).First(&existing).Error; err == nil {
			shortCode = existing.ShortCode
		} else {
			// Sinh short code nếu chưa có
			if shortCode == "" {
				shortCode = GenerateUniqueCode(db, 6)
			}
			url := Db.Url{
				OriginalUrl: longUrl,
				ShortCode:   shortCode,
				VisitCount:  0,
			}
			db.Create(&url)
		}

		// Hiển thị short URL cho người dùng click
		c.HTML(http.StatusOK, "home.html", gin.H{
			"message":  "URL đã lưu thành công!",
			"shortUrl": "http://localhost:8080/" + shortCode,
			"original": longUrl,
		})
	})

	// Redirect từ short URL
	r.GET("/:shortCode", func(c *gin.Context) {
		shortCode := c.Param("shortCode")
		var url Db.Url
		if err := db.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
			c.String(http.StatusNotFound, "URL không tồn tại")
			return
		}

		// Tăng số lượt truy cập
		db.Model(&url).Update("visit_count", url.VisitCount+1)

		// Redirect tới URL gốc
		c.Redirect(http.StatusFound, url.OriginalUrl)
	})

	r.Run(":8080")
}
