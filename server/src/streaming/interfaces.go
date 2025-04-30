package streaming

import "time"

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
