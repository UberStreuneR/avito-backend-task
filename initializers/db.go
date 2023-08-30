package initializers

import (
	"avito-task/entity"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	var dsn string
	if config.Dev {
		dsn = fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5433 sslmode=disable", config.DBUserName, config.DBUserPassword, config.DBName)
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)
	}
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to the database \n", err.Error())
		os.Exit(1)
	}
	err = DB.AutoMigrate(&entity.User{}, &entity.Segment{})
	if err != nil {
		log.Fatal("Migration failed \n", err.Error())
		os.Exit(1)
	}
	log.Println("Connected to the database")
}

func GetTestDB(config *Config) *gorm.DB {
	ConnectDB(config)
	DB.Exec("DROP DATABASE IF EXISTS test")
	result := DB.Exec("CREATE DATABASE test")
	if result.Error != nil {
		log.Fatal("Failed to create database 'test'. ", result.Error.Error())
		os.Exit(1)
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.TestDBHost, config.DBUserName, config.DBUserPassword, config.TestDBName, config.TestDBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database \n", err.Error())
		os.Exit(1)
	}
	err = db.AutoMigrate(&entity.User{}, &entity.Segment{})
	if err != nil {
		log.Fatal("Migration failed \n", err.Error())
		os.Exit(1)
	}
	return db
}
