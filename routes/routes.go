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
		AllowedOrigins: []string{"https//*","http://*"},
		AllowedMethods: []string{"GET","POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type","X-XSRF=Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	})) 

	router.Use(middleware.Heartbeat("/ping"))
	
	router.Get("/getEvents", handlers.QueryEvents)
	router.Post("/register", handlers.Register)
	router.Post("/login", handlers.Login)

	return router
	
}