package main

import (
	"audit-log-service/config"
	"audit-log-service/db"
	"audit-log-service/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

const (
	webPort = "8000"
	rpcPort = "5001"
)

type Config struct {
	DB *db.DBConnection
	logger log.Logger
}


func main() {
	var err error

	fmt.Println("Connecting to db")
	if err := db.NewConnection(db.DB); err != nil {
		panic("could not connect to the database")
	}
	db.DBConn = &db.DBConnection{DB: db.DB}

	db.DBConn.Init()
	cfgFile := "config.yml"
	
	config.AuthConfiguration, err = config.LoadConfig(cfgFile)
	
	app := Config{DB: db.DBConn, logger: *log.Default()}
	app.logger.Printf("Starting Server at port %s", webPort)

	go func() {
		s := rpc.NewServer()
		s.RegisterCodec(json.NewCodec(), "application/json")
		// s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
		rpcServer := new(RPCServer)
		s.RegisterService(rpcServer, "")
		r := mux.NewRouter()
		r.Handle("/rpc",s)
		fmt.Println("serving on 5001")
        err := http.ListenAndServe(":5001", r)
		fmt.Println("Failed to serve")
        if err != nil {
            panic("ListenAndServe: " + err.Error())
        }
	  }()

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: routes.SetupRoutes(),
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("failed to listen for requests", err)
		log.Panic()
	}	
}

func basicHandler(w http.ResponseWriter,r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
