package post

import (
	"context"
	"errors"
)


type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{repo: repo}
}

// CreatePost menghandle logika pembuatan post baru
func (s *Service) CreatePost(ctx context.Context, req *CreatePostRequest) (*Post, error) {
	//Validasi input
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.Content == "" {
		return nil, errors.New("content is required")
	}
	Post := &Post{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Tags:     req.Tags,
	}
	


	if err := s.repo.CreatePost(ctx, Post); err != nil {
		return nil, err
	}
	return Post, nil
}

// GetAllPosts menghandle logika pengambilan semua post
func (s *Service) GetAllPosts(ctx context.Context) ([]*Post, error) {
	return s.repo.GetAll(ctx)
}

// GetPostByID menghandle logika pengambilan post berdasarkan ID
func (s *Service) GetPostByID(ctx context.Context, id int64) (*Post, error) {
	return s.repo.GetByID(ctx, id)
}

//UpdatePost menghandle logika pembaruan post
func (s *Service) UpdatePost(ctx context.Context, ID int64, req *updatePostInput) (*Post, error) {
	// cek apakah post ada
	existingPost, err := s.repo.GetByID(ctx, ID)
	if err != nil {
		return nil, errors.New("Post not found")
	}
	// validasi input
	if req.Title == nil || *req.Title == "" {
		return nil, errors.New("title is required")
	}
	// update field yang dikirim
	if req.Title != nil {
		existingPost.Title = *req.Title
	}
	if req.Content != nil {
		existingPost.Content = *req.Content
	}
	if req.Category != nil {
		existingPost.Category = *req.Category
	}
	if req.Tags != nil {
		existingPost.Tags = req.Tags
	}

	// simpan perubahan
	if err := s.repo.UpdatePost(ctx, ID, req); err != nil {
		return nil, err
	}

	return existingPost, nil
}
// DeletePost menghandle logika penghapusan post
func (s *Service) DeletePost(ctx context.Context, id int64) error {
	// cek apakah post ada
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.DeletePost(ctx, id)
}

// SearchPosts menghandle logika pencarian post berdasarkan kata kunci
func (s *Service) SearchPosts(ctx context.Context, keyword string) ([]*Post, error){
	if keyword == "" {
		return s.repo.GetAll(ctx)
	}
	return s.repo.SearchPosts(ctx, keyword)
}