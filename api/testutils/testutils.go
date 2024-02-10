package testutils

import (
	"audit-log-service/db"
	"log"
	"os"
)


func InitTestDb() (error) {
	if err := db.NewTestConnection(db.DB); err != nil {
		return err
	}
	db.DBConn = &db.DBConnection{DB: db.DB}

	db.DBConn.Init()
	return nil
}

func RemoveTestDb() {
	e := os.Remove("test.db") 
    if e != nil { 
        log.Fatal(e) 
    } 
}