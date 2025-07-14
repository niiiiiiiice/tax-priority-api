package persistence

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig конфигурация базы данных
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "root",
		Password: "root",
		DBName:   "tax_priority",
		SSLMode:  "disable",
	}
}

func Connect(config *DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Получаем underlying sql.DB для настройки пула соединений
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Настраиваем пул соединений
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to database")
	return db, nil
}

// Migrate выполняет миграции базы данных
func Migrate(db *gorm.DB) error {
	log.Println("Starting database migration...")

	// Автоматическая миграция всех моделей
	err := db.AutoMigrate()

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// Close закрывает соединение с базой данных
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Close()
}
