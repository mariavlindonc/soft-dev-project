package db

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"backend/logger"

	sqlDriver "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const maxRetries = 15
const retryDelay = 2 * time.Second

type Database struct {
	conn *gorm.DB
}

func NewDatabase() (*Database, error) {
	// Load .env
	if err := godotenv.Load(); err != nil {
		logger.Warn(".env file not found, using system env vars")
	}

	// Read credentials
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "3306")
	user := getEnvOrDefault("DB_USER", "root")
	pass := getEnvOrDefault("DB_PASSWORD", "")
	name := getEnvOrDefault("DB_NAME", "ceibo_db")

	// Configure TLS for MySQL if CA cert is provided
	caPath := os.Getenv("DB_SSL_CA")
	if caPath != "" {
		if err := configureMySQLTLS(caPath); err != nil {
			return nil, fmt.Errorf("failed to configure mysql TLS: %w", err)
		}
	}

	// Create database if it does not exist
	if err := ensureDatabase(host, port, user, pass, name, caPath != ""); err != nil {
		return nil, fmt.Errorf("failed to ensure database: %w", err)
	}

	// Connect with retry (handles Docker depends_on race condition)
	tlsParam := "false"
	if caPath != "" {
		tlsParam = "custom"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%s",
		user, pass, host, port, name, tlsParam)

	var conn *gorm.DB
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		conn, lastErr = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Info),
		})
		if lastErr == nil {
			break
		}
		logger.Warn("database connection attempt %d/%d failed: %v", attempt, maxRetries, lastErr)
		time.Sleep(retryDelay)
	}
	if lastErr != nil {
		return nil, fmt.Errorf("failed to connect to database after %d retries: %w", maxRetries, lastErr)
	}

	logger.Info("database connected successfully")
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

	logger.Info("migrations completed successfully")
	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func configureMySQLTLS(caPath string) error {
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(caPath)
	if err != nil {
		return fmt.Errorf("failed to read CA cert: %w", err)
	}
	if !rootCertPool.AppendCertsFromPEM(pem) {
		return fmt.Errorf("failed to parse CA cert")
	}

	tlsConfig := &tls.Config{
		RootCAs: rootCertPool,
	}
	return sqlDriver.RegisterTLSConfig("custom", tlsConfig)
}

func ensureDatabase(host, port, user, pass, name string, useTLS bool) error {
	tlsParam := "false"
	if useTLS {
		tlsParam = "custom"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local&tls=%s",
		user, pass, host, port, tlsParam)

	tmpDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
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

	logger.Info("database `%s` ensured", name)
	return nil
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
