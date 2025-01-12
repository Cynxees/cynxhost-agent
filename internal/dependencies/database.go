package dependencies

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DatabaseClient struct {
	Db *gorm.DB
}

func NewDatabaseClient(config *Config) (*DatabaseClient, error) {
	// Construct the DSN (Data Source Name)
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.Database.MySQL.Username,
		config.Database.MySQL.Password,
		config.Database.MySQL.Host,
		config.Database.MySQL.Port,
		config.Database.MySQL.Database,
	)

	// Open a connection with GORM using the MySQL driver
	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable GORM's built-in logging
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Check the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %w", err)
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DatabaseClient{Db: db}, nil
}

func (client *DatabaseClient) Close() error {
	sqlDB, err := client.Db.DB()
	if err != nil {
		return fmt.Errorf("failed to get generic database object: %w", err)
	}
	return sqlDB.Close()
}
