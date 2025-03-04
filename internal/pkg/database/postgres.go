package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
func ConnectPgsql(host string, port string, username string, password string, dbname string) (db *gorm.DB) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable", host, username, password, port, dbname)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connected to PostgreSQL!")
	
	// 连接池配置（可选）
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)   // 空闲连接池大小
	sqlDB.SetMaxOpenConns(100)  // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间
	return db
}