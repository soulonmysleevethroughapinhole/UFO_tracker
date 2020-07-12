package comments

import "time"

type Comment struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	ThreadID int64     `json:"threadid"`
	ParentID int64     `json:"parentid"`
	Content  string    `json:"content"`
	PostDate time.Time `json:"postdate"`
}

type Comments []Comment
