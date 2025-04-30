package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rafixcs/homestreaming/src/presentation/api/middleware"
	"github.com/rafixcs/homestreaming/src/presentation/api/routes"
)

const PORT = "8080"

func RunServer() {
	r := chi.NewRouter()
	r.Use(chi_middleware.Logger)

	corsHandler := middleware.NewCorsMiddleware()
	r.Use(corsHandler.Handler)

	routes.LoadRoutes(r)

	log.Printf("Server running on port %v\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, r))
}
