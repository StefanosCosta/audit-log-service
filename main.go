package main

import (
	"audit-log-service/db"
	"audit-log-service/routes"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

const (
	webPort = "8000"
	rpcPort = "5001"
	gRpcPort = "50001"
)

type Config struct {
	DB *db.DBConnection
	logger log.Logger
}


func main() {

	fmt.Println("Connecting to db")
	if err := db.NewConnection(db.DB); err != nil {
		panic("could not connect to the database")
	}
	db.DBConn = &db.DBConnection{DB: db.DB}

	db.DBConn.Init()
	
	app := Config{DB: db.DBConn, logger: *log.Default()}
	app.logger.Printf("Starting Server at port %s", webPort)

	go rpcListen()

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: routes.SetupRoutes(),
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("failed to listen for requests", err)
		log.Panic()
	}	
}

// func serve() {
// 	server := &http.Server{
// 		Addr: fmt.Sprintf(":%s", webPort),
// 		Handler: routes.SetupRoutes(),
// 	}
// 	err := server.ListenAndServe()
// 	if err != nil {
// 		fmt.Println("failed to listen for requests", err)
// 		log.Panic()
// 	}
// }

func rpcListen() error {
	log.Println("Starting RPC server on port ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}

}

func basicHandler(w http.ResponseWriter,r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
