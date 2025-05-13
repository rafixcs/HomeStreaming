package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/rafixcs/homestreaming/src/streaming"
	"github.com/rafixcs/homestreaming/src/utils"
)

func ServeVideos(r chi.Router) {
	workDir, _ := os.Getwd()
	thumbnailsDir := http.Dir(filepath.Join(workDir, "assets/thumbnails"))
	streaming.ServeFiles(r, "/thumbnail", thumbnailsDir)
	r.Get("/video/{}", streaming.ServeVideo)
}

func ListVideos(w http.ResponseWriter, r *http.Request) {
	projectPaths, err := utils.GetProjectPaths()
	if err != nil {
		err = fmt.Errorf("LoadRoutes: failed to get project paths: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	videosPath := projectPaths.ProjectRoot + "/assets/videos/"
	videos, err := streaming.ScanVideoDirectory(videosPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseContent := videos
	json.NewEncoder(w).Encode(&responseContent)
}
