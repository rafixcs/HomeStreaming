package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

const PORT = "8080"

type Video struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	ThumbnailURL string    `json:"thumbnailUrl"`
	VideoURL     string    `json:"videoUrl"`
	Duration     string    `json:"duration"`
	Resolution   string    `json:"resolution"`
	Format       string    `json:"format"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `json:"createdAt"`
}

type VideoMetadata struct {
	Streams []struct {
		CodecType string `json:"codec_type"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		Duration  string `json:"duration"`
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
		Size     string `json:"size"`
		Format   string `json:"format_name"`
	} `json:"format"`
}

type VideoListResponse []Video

func main() {

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Add your React app URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//r.Use(JsonApplicationHeader)
	r.Use(corsHandler.Handler)
	//  r.Use(cacheControl)

	r.Get("/videos", func(w http.ResponseWriter, r *http.Request) {
		//path := http.Dir("./assets/videos/")
		videos, err := scanVideoDirectory("./assets/videos/")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseContent := videos
		json.NewEncoder(w).Encode(&responseContent)
	})

	//fs := http.FileServer(http.Dir("./assets/videos/"))
	//r.Handle("/video/", http.StripPrefix("/video/", fs))

	//fs = http.FileServer(http.Dir("./assets/thumbnails/"))
	//r.Handle("/thumbnail/", http.StripPrefix("/thumbnail/", fs))

	workDir, _ := os.Getwd()
	thumbnailsDir := http.Dir(filepath.Join(workDir, "assets/thumbnails"))
	//videosDir := http.Dir(filepath.Join(workDir, "assets/videos"))
	// Serve files with custom handling
	ServeFiles(r, "/thumbnail", thumbnailsDir)
	r.Get("/video/{}", ServeVideo)

	// Serve static files (JS, CSS, images, etc.)
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Printf("Server running on port %v\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, r))
}

func JsonApplicationHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func scanVideoDirectory(dirPath string) ([]Video, error) {
	var videos []Video

	// Supported video formats
	supportedFormats := map[string]bool{
		".mp4": true,
		".mkv": true,
		".avi": true,
		".mov": true,
		".wmv": true,
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if supportedFormats[ext] {
			video, err := processVideoFile(path, info)
			if err != nil {
				log.Printf("Error processing video %s: %v", path, err)
				return nil
			}
			videos = append(videos, video)
		}

		return nil
	})

	return videos, err
}

func processVideoFile(path string, info os.FileInfo) (Video, error) {
	filename := filepath.Base(path)
	id := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Generate thumbnail if it doesn't exist
	thumbnailPath := filepath.Join("./assets/thumbnails", id+".jpg")
	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		err = generateThumbnail(path, thumbnailPath)
		if err != nil {
			return Video{}, fmt.Errorf("thumbnail generation failed: %v", err)
		}
	}

	// Get video metadata
	metadata, err := getVideoMetadata(path)
	if err != nil {
		return Video{}, fmt.Errorf("metadata extraction failed: %v", err)
	}

	// Create video object
	video := Video{
		ID:           id,
		Title:        strings.TrimSuffix(filename, filepath.Ext(filename)),
		ThumbnailURL: fmt.Sprintf("http://localhost:8080/thumbnail/%s.jpg", id),
		VideoURL:     fmt.Sprintf("http://localhost:8080/video/%s", filename),
		Size:         info.Size(),
		CreatedAt:    info.ModTime(),
	}

	// Add metadata if available
	if metadata != nil && len(metadata.Streams) > 0 {
		for _, stream := range metadata.Streams {
			if stream.CodecType == "video" {
				video.Resolution = fmt.Sprintf("%dx%d", stream.Width, stream.Height)
				video.Duration = stream.Duration
				break
			}
		}
		video.Format = metadata.Format.Format
	}

	return video, nil
}

func generateThumbnail(videoPath, thumbnailPath string) error {
	// Generate thumbnail at 1 second mark
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:01.000",
		"-vframes", "1",
		"-vf", "scale=320:-1",
		thumbnailPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %v, output: %s", err, string(output))
	}

	return nil
}

func getVideoMetadata(videoPath string) (*VideoMetadata, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		videoPath)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %v", err)
	}

	var metadata VideoMetadata
	if err := json.Unmarshal(output, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %v", err)
	}

	return &metadata, nil
}

// ServeFiles serves static files with custom handling
func ServeFiles(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	//path += "*"

	r.Get(path+"{}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request")
		thumb := strings.Split(r.URL.Path, "/")
		path := filepath.Join("./assets/thumbnails/", thumb[len(thumb)-1])
		log.Println(path)

		// Custom error handling
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// Set cache headers
		w.Header().Set("Cache-Control", "max-age=31536000, public")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		fs.ServeHTTP(w, r)
	})
}

// Optional: Add middleware for specific routes
func thumbnailMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add any specific handling for thumbnail requests
		next.ServeHTTP(w, r)
	})
}

// Cache control middleware
func cacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set cache headers for thumbnails
		w.Header().Set("Cache-Control", "max-age=31536000, public")
		w.Header().Set("Expires", time.Now().Add(time.Hour*24*365).Format(http.TimeFormat))
		next.ServeHTTP(w, r)
	})
}

func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".ogg":
		return "video/ogg"
	case ".mkv":
		return "video/x-matroska"
	case ".avi":
		return "video/x-msvideo"
	case ".mov":
		return "video/quicktime"
	default:
		return "application/octet-stream"
	}
}

func ServeVideo(w http.ResponseWriter, r *http.Request) {
	// Extract the filename from the URL
	filePath := strings.TrimPrefix(r.URL.Path, "/video/")
	videoPath := filepath.Join("./assets/videos", filePath)

	// Open the video file
	video, err := os.Open(videoPath)
	if err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}
	defer video.Close()

	// Get video file information
	videoInfo, err := video.Stat()
	if err != nil {
		http.Error(w, "Failed to get video info", http.StatusInternalServerError)
		return
	}

	// Get the video size
	videoSize := videoInfo.Size()

	// Get the content type based on file extension
	contentType := getContentType(videoPath)
	w.Header().Set("Content-Type", contentType)

	// Handle range requests for video streaming
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		// Parse the range header
		ranges := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
		if len(ranges) != 2 {
			http.Error(w, "Invalid range header", http.StatusBadRequest)
			return
		}

		// Parse start range
		start, err := strconv.ParseInt(ranges[0], 10, 64)
		if err != nil {
			http.Error(w, "Invalid range header", http.StatusBadRequest)
			return
		}

		// Parse end range
		var end int64
		if ranges[1] == "" {
			end = videoSize - 1
		} else {
			end, err = strconv.ParseInt(ranges[1], 10, 64)
			if err != nil {
				http.Error(w, "Invalid range header", http.StatusBadRequest)
				return
			}
		}

		// Validate ranges
		if start >= videoSize || end >= videoSize {
			w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", videoSize))
			http.Error(w, "Requested range not satisfiable", http.StatusRequestedRangeNotSatisfiable)
			return
		}

		// Set headers for partial content
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, videoSize))
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
		w.WriteHeader(http.StatusPartialContent)

		// Seek to start position
		_, err = video.Seek(start, io.SeekStart)
		if err != nil {
			http.Error(w, "Failed to seek video", http.StatusInternalServerError)
			return
		}

		// Stream the video chunk
		_, err = io.CopyN(w, video, end-start+1)
		if err != nil {
			return // Client probably disconnected
		}
	} else {
		// No range requested, serve the entire file
		w.Header().Set("Content-Length", fmt.Sprintf("%d", videoSize))
		w.Header().Set("Accept-Ranges", "bytes")
		_, err = io.Copy(w, video)
		if err != nil {
			return // Client probably disconnected
		}
	}
}
