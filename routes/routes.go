package routes

import (
	"audit-log-service/handlers"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)




func SetupRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"httpsL//*","http://*"},
		AllowedMethods: []string{"GET","POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type","X-XSRF=Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	})) 

	router.Use(middleware.Heartbeat("/ping"))

	router.Post("/addEvent",handlers.HandleEvent)
	router.Get("getEvents", handlers.QueryEvents)

	return router
	
}