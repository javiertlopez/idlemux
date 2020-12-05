package model

// Video struct
type Video struct {
	ID          *string  `json:"id,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	SourceURL   string   `json:"source_url,omitempty"`
	Asset       *Asset   `json:"asset,omitempty"`
	Duration    *float64 `json:"duration,omitempty"`
	Poster      string   `json:"poster,omitempty"`
	Thumbnail   string   `json:"thumbnail,omitempty"`
	Sources     []Source `json:"sources,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}
