package db

import (
	"fmt"
	"os"

	"drug-info/backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		user := getEnv("MYSQL_USER", "root")
		password := getEnv("MYSQL_PASSWORD", "root")
		host := getEnv("MYSQL_HOST", "127.0.0.1")
		port := getEnv("MYSQL_PORT", "3306")
		database := getEnv("MYSQL_DATABASE", "medical_info")

		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user,
			password,
			host,
			port,
			database,
		)
	}

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := conn.AutoMigrate(&models.Drug{}, &models.SpecimenApplication{}, &models.SysUser{}); err != nil {
		return nil, err
	}

	return conn, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
