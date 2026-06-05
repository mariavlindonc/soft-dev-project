package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	conn *gorm.DB
}

func NewDatabase() (*Database, error) {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found, using system env vars")
	}

	// Read credentials
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "3306")
	user := getEnvOrDefault("DB_USER", "root")
	pass := getEnvOrDefault("DB_PASSWORD", "")
	name := getEnvOrDefault("DB_NAME", "ceibo_db")

	// Create database if it does not exist
	if err := ensureDatabase(host, port, user, pass, name); err != nil {
		return nil, fmt.Errorf("failed to ensure database: %w", err)
	}

	// Connect
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("database connected successfully")
	return &Database{conn: conn}, nil
}

func (d *Database) GetConnection() *gorm.DB {
	return d.conn
}

func (d *Database) Migrate(models ...interface{}) error {
	if err := d.conn.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}

	d.conn.Exec("ALTER TABLE events DROP CHECK chk_events_capacity")
	if err := d.conn.Exec(
		"ALTER TABLE events ADD CONSTRAINT chk_events_capacity CHECK (tickets_sold >= 0 AND tickets_sold <= capacity)",
	).Error; err != nil {
		return fmt.Errorf("failed to add CHECK constraint: %w", err)
	}

	log.Println("migrations completed successfully")
	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func ensureDatabase(host, port, user, pass, name string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port)

	tmpDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %w", err)
	}

	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", name)
	if err := tmpDB.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	sqlDB, err := tmpDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sqlDB: %w", err)
	}
	defer sqlDB.Close()

	log.Printf("database `%s` ensured", name)
	return nil
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
