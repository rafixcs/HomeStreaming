package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/rafixcs/homestreaming/src/presentation/api/controller"
	"github.com/rafixcs/homestreaming/src/presentation/api/middleware"
)

func LoadRoutes(mux *chi.Mux) error {
	mux.Group(func(r chi.Router) {
		r.Use(middleware.JsonApplicationHeader)

		r.Get("/videos", controller.ListVideos)
	})

	mux.Group(func(r chi.Router) {
		controller.ServeVideos(r)
	})

	return nil
}
