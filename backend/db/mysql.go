package db

import (
	"fmt"
	"os"

	"drug-info/backend/config"
	"drug-info/backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(mysqlConfig config.MySQLConfig) (*gorm.DB, error) {
	dsn := buildDSN(mysqlConfig)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := conn.AutoMigrate(&models.Drug{}, &models.SpecimenApplication{}, &models.SysUser{}, &models.CasbinRule{}, &models.FileUploadAndDownload{}); err != nil {
		return nil, err
	}

	return conn, nil
}

func buildDSN(mysqlConfig config.MySQLConfig) string {
	if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
		return dsn
	}
	if mysqlConfig.DSN != "" {
		return mysqlConfig.DSN
	}

	username := getEnv("MYSQL_USER", defaultString(mysqlConfig.Username, "root"))
	password := getEnv("MYSQL_PASSWORD", defaultString(mysqlConfig.Password, "root"))
	host := getEnv("MYSQL_HOST", defaultString(mysqlConfig.Host, "127.0.0.1"))
	port := getEnv("MYSQL_PORT", defaultString(mysqlConfig.Port, "3306"))
	database := getEnv("MYSQL_DATABASE", defaultString(mysqlConfig.Database, "medical_info"))
	charset := getEnv("MYSQL_CHARSET", defaultString(mysqlConfig.Charset, "utf8mb4"))
	loc := getEnv("MYSQL_LOC", defaultString(mysqlConfig.Loc, "Local"))
	parseTime := "False"
	if mysqlConfig.ParseTime || getEnv("MYSQL_PARSE_TIME", "") == "true" {
		parseTime = "True"
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		parseTime,
		loc,
	)
}

func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
