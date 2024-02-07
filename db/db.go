package db

import (
	"audit-log-service/models"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DBConn *DBConnection

func NewConnection(db *gorm.DB) error{
	err := createNewConnection(db)
	if err != nil {
		return err
	}
	return nil
}


type DBConnection struct {
	DB *gorm.DB
} 

func (connection *DBConnection) Init() {
	var err error

	if err = connection.DB.AutoMigrate(&models.Event{}); err != nil {
		fmt.Println("could not automigrate events")
		return
	}
}

func createNewConnection(db *gorm.DB) error {
	var err error

	db, err = gorm.Open(sqlite.Open("audit.db"),&gorm.Config{})

	if err != nil {
		return err
	}
	return nil
}

func createNewTestConnection(db *gorm.DB) error {
	var err error

	db, err = gorm.Open(sqlite.Open("test.db"),&gorm.Config{})

	if err != nil {
		return err
	}
	return nil
}