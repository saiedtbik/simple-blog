package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)
var DB *gorm.DB

func ConnectDatabase() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:  logger.Info, // Log level Info will output everything
		},
	)
	database, err := gorm.Open(sqlite.Open("test2.db"),  &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&User{},&Role{},&UserRole{},&RolePermissionResource{},&Resource{})
	if err != nil {
		return
	}

	DB = database
}
