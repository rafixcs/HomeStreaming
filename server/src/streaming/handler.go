package streaming

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

	"github.com/go-chi/chi/v5"
)

func ScanVideoDirectory(dirPath string) ([]Video, error) {
	var videos []Video

	// Supported video formats
	supportedFormats := map[string]bool{
		".mp4": true,
		".mkv": true,
		".avi": true,
		".mov": true,
		".wmv": true,
	}

	log.Println(dirPath)

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
		ThumbnailURL: fmt.Sprintf("http://192.168.0.14:8080/thumbnail/%s.jpg", id),
		VideoURL:     fmt.Sprintf("http://192.168.0.14:8080/video/%s", filename),
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

	r.Get(path+"{}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request")
		thumb := strings.Split(r.URL.Path, "/")
		path := filepath.Join("./assets/thumbnails/", thumb[len(thumb)-1])
		log.Println(path)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Cache-Control", "max-age=31536000, public")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		fs.ServeHTTP(w, r)
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
	filePath := strings.TrimPrefix(r.URL.Path, "/video/")
	videoPath := filepath.Join("./assets/videos", filePath)

	video, err := os.Open(videoPath)
	if err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}
	defer video.Close()

	videoInfo, err := video.Stat()
	if err != nil {
		http.Error(w, "Failed to get video info", http.StatusInternalServerError)
		return
	}

	videoSize := videoInfo.Size()

	contentType := getContentType(videoPath)
	w.Header().Set("Content-Type", contentType)

	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		ranges := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
		if len(ranges) != 2 {
			http.Error(w, "Invalid range header", http.StatusBadRequest)
			return
		}

		start, err := strconv.ParseInt(ranges[0], 10, 64)
		if err != nil {
			http.Error(w, "Invalid range header", http.StatusBadRequest)
			return
		}

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

		if start >= videoSize || end >= videoSize {
			w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", videoSize))
			http.Error(w, "Requested range not satisfiable", http.StatusRequestedRangeNotSatisfiable)
			return
		}

		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, videoSize))
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
		w.WriteHeader(http.StatusPartialContent)

		_, err = video.Seek(start, io.SeekStart)
		if err != nil {
			http.Error(w, "Failed to seek video", http.StatusInternalServerError)
			return
		}

		_, err = io.CopyN(w, video, end-start+1)
		if err != nil {
			return
		}
	} else {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", videoSize))
		w.Header().Set("Accept-Ranges", "bytes")
		_, err = io.Copy(w, video)
		if err != nil {
			return
		}
	}
}
