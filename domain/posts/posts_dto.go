package posts

import "time"

const (
	StatusPublic = "public"
	layout       = "2006-01-02"
)

type Post struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Title     string    `json:"title"`
	Corpus    string    `json:"corpus"`
	MediaType string    `json:"mediatype"`
	Media     []byte    `json:"media"`
	PostDate  time.Time `json:"postdate"`
}

type Posts []Post
