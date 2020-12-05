package model

// Asset Information from Mux
type Asset struct {
	ID                  string   `json:"id,omitempty"`
	URL                 string   `json:"url,omitempty"`
	CreatedAt           string   `json:"created_at,omitempty"`
	Status              string   `json:"status,omitempty"`
	Duration            float64  `json:"duration,omitempty"`
	MaxStoredResolution string   `json:"max_stored_resolution,omitempty"`
	MaxStoredFrameRate  float64  `json:"max_stored_frame_rate,omitempty"`
	AspectRatio         string   `json:"aspect_ratio,omitempty"`
	Passthrough         string   `json:"passthrough,omitempty"`
	Poster              string   `json:"poster,omitempty"`
	Thumbnail           string   `json:"thumbnail,omitempty"`
	Sources             []Source `json:"sources,omitempty"`
}

// Source manifests
type Source struct {
	Source string `json:"src"`
	Type   string `json:"type"`
}
