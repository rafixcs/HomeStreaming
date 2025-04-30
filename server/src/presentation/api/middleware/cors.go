package middleware

import "github.com/rs/cors"

func NewCorsMiddleware() *cors.Cors {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Add your React app URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	return corsHandler
}
