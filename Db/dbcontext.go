package Db

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	server := "MKHANGDZ1ST"
	port := 1433
	user := "sa"
	password := "mkhang123abc"
	database := "BtWeek1"

	// 1 Kết nối tới master trước
	masterDSN := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=master&encrypt=true&trustServerCertificate=true",
		user, password, server, port)

	masterDB, err := gorm.Open(sqlserver.Open(masterDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Không kết nối được SQL Server: %v", err)
	}

	// 2️ Tạo database nếu chưa có
	createDBSQL := fmt.Sprintf("IF DB_ID(N'%s') IS NULL CREATE DATABASE %s;", database, database)
	if err := masterDB.Exec(createDBSQL).Error; err != nil {
		log.Fatalf("❌ Không thể tạo database: %v", err)
	}

	// 3️ Kết nối lại tới database mới tạo
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=true&trustServerCertificate=true",
		user, password, server, port, database)

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Không kết nối được database %s: %v", database, err)
	}

	//  Auto migrate bảng Url
	if err := db.AutoMigrate(&Url{}); err != nil {
		log.Fatalf("❌ Lỗi migrate DB: %v", err)
	}

	fmt.Println("✅ Database + bảng Url đã tạo thành công!")
	return db
}
