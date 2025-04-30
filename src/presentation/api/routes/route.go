package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/rafixcs/homestreaming/src/presentation/api/middleware"
	"github.com/rafixcs/homestreaming/src/streaming"
	"github.com/rafixcs/homestreaming/src/utils"
)

func LoadRoutes(mux *chi.Mux) error {
	projectPaths, err := utils.GetProjectPaths()
	if err != nil {
		err = fmt.Errorf("LoadRoutes: failed to get project paths: %w", err)
		return err
	}

	mux.Group(func(r chi.Router) {
		r.Use(middleware.JsonApplicationHeader)

		r.Get("/videos", func(w http.ResponseWriter, r *http.Request) {
			videosPath := projectPaths.ProjectRoot + "/assets/videos/"
			videos, err := streaming.ScanVideoDirectory(videosPath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			responseContent := videos
			json.NewEncoder(w).Encode(&responseContent)
		})
	})

	mux.Group(func(r chi.Router) {
		workDir, _ := os.Getwd()
		thumbnailsDir := http.Dir(filepath.Join(workDir, "assets/thumbnails"))
		streaming.ServeFiles(r, "/thumbnail", thumbnailsDir)
		r.Get("/video/{}", streaming.ServeVideo)
	})

	return nil
}
