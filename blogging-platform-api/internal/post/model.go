package post

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Tags      pq.StringArray  `json:"tags" db:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostRequest struct{
	Title   string `json:"title"`
	Content string `json:"content"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

type updatePostInput struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Category *string   `json:"category"`
	Tags     pq.StringArray `json:"tags"`
}