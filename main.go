package main

import (
	"audit-log-service/db"
	"audit-log-service/routes"
	"fmt"
	"log"
	"net/http"
)

const (
	webPort = "80"
	rpcPort = "5001"
	gRpcPort = "50001"
)

type Config struct {
	DB *db.DBConnection
}


func main() {

	// http.HandleFunc("/events", handlers.HandleEvents)
	fmt.Println("Connecting to db")
	if err := db.NewConnection(db.DBConn.DB); err != nil {
		panic("could not connect to the database")
	}

	db.DBConn.Init()
	// fmt.Println("Starting server")
	app := Config{
		DB:db.DBConn,
	}

	go app.serve()
	
}

func (app *Config) serve() {
	server := &http.Server{
		Addr: fmt.Sprintf("%s", webPort),
		Handler: routes.SetupRoutes(),
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("failed to listen for requests", err)
		log.Panic()
	}
}


func basicHandler(w http.ResponseWriter,r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
