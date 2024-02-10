package db

import (
	events "audit-log-service/db/eventsRepository"
	users "audit-log-service/db/usersRepository"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DBConn *DBConnection
var DB *gorm.DB

func NewConnection(db *gorm.DB) error{
	err := createNewConnection(db)
	if err != nil {
		return err
	}
	return nil
}

func createNewConnection(db *gorm.DB) error {
	var err error

	db, err = gorm.Open(sqlite.Open("audit.db"),&gorm.Config{})

	if err != nil {
		return err
	}
	DB = db
	return nil
}

type DBConnection struct {
	DB *gorm.DB
} 

func (connection *DBConnection) Init() {
	var err error

	if err = connection.DB.AutoMigrate(&events.Event{}); err != nil {
		fmt.Println("could not automigrate events")
		return
	}

	if err = connection.DB.AutoMigrate(&users.User{}); err != nil {
		fmt.Println("could not automigrate events")
		return
	}
}



func (connection *DBConnection) NewTestConnection(db *gorm.DB) error {
	var err error

	db, err = gorm.Open(sqlite.Open("test.db"),&gorm.Config{})

	if err != nil {
		return err
	}
	return nil
}