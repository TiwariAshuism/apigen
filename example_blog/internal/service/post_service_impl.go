package service

import (
	"context"
	"errors"

	"github.com/example/blog/internal/model"
	"github.com/example/blog/internal/repository"
)

type postServiceImpl struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postServiceImpl{repo: repo}
}

func (s *postServiceImpl) ListPosts(ctx context.Context, req model.ListPostsRequest) ([]model.Post, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	return s.repo.ListPosts(ctx, req)
}

func (s *postServiceImpl) GetPost(ctx context.Context, id string) (model.Post, error) {
	if id == "" {
		return model.Post{}, errors.New("id is required")
	}
	return s.repo.GetPost(ctx, id)
}

func (s *postServiceImpl) CreatePost(ctx context.Context, req model.CreatePostRequest) (model.Post, error) {
	if req.Title == "" {
		return model.Post{}, errors.New("title is required")
	}
	if req.Body == "" {
		return model.Post{}, errors.New("body is required")
	}
	return s.repo.CreatePost(ctx, req)
}

func (s *postServiceImpl) UpdatePost(ctx context.Context, req model.UpdatePostRequest, id string) (model.Post, error) {
	if id == "" {
		return model.Post{}, errors.New("id is required")
	}
	return s.repo.UpdatePost(ctx, req, id)
}

func (s *postServiceImpl) DeletePost(ctx context.Context, id string) (model.Post, error) {
	if id == "" {
		return model.Post{}, errors.New("id is required")
	}
	return s.repo.DeletePost(ctx, id)
}
