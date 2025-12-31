package post

import (
	"context"
	"database/sql"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}
// CreatePost: Membuat post baru
func (r *repository) CreatePost(ctx context.Context, p *Post) error {
	query := "INSERT INTO posts (title, content, category, tags, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query, 
		p.Title, 
		p.Content, 
		p.Category, 
		p.Tags, 
		p.CreatedAt, 
		p.UpdatedAt,
	).Scan(&p.ID)
	
	if err != nil {
		return err
	}
	return nil
}
// GetAll: Mengambil semua post
func (r *repository) GetAll(ctx context.Context) ([]*Post, error) {
	query := "SELECT id, title, content, created_at, updated_at FROM posts"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
// GETBYID: Mengambil post berdasarkan ID
func (r *repository) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := "SELECT id, title, content, created_at, updated_at FROM posts WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	var p Post
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &p, nil
}

// UPDATE: 	Mengubah data post
func (r *repository) UpdatePost(ctx context.Context, id int64, req *updatePostInput) error {
	query := "UPDATE posts SET title = COALESCE($1, title), content = COALESCE($2, content), updated_at = $3 WHERE id = $4"

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, req.Title, req.Content, now, id)
	if err != nil {
		return err
	}
	return nil
}
// DELETE: 	Menghapus post
func (r *repository) DeletePost(ctx context.Context, id int64) error {
	query := "DELETE FROM posts WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

// SEARCH: Mencari post berdasarkan kata kunci pada judul atau konten
func (r *repository) SearchPosts(ctx context.Context, keyword string) ([]*Post, error) {
	query := "SELECT id, title, content, created_at, updated_at FROM posts WHERE title ILIKE $1 OR content ILIKE $2"
	likePattern := "%" + keyword + "%"
	rows, err := r.db.QueryContext(ctx, query, likePattern, likePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}