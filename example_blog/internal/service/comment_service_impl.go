package service

import (
	"context"
	"errors"

	"github.com/example/blog/internal/model"
	"github.com/example/blog/internal/repository"
)

type commentServiceImpl struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentServiceImpl{repo: repo}
}

func (s *commentServiceImpl) ListComments(ctx context.Context, req model.ListCommentsRequest, postId string) ([]model.Comment, error) {
	if postId == "" {
		return nil, errors.New("postId is required")
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	return s.repo.ListComments(ctx, req, postId)
}

func (s *commentServiceImpl) CreateComment(ctx context.Context, req model.CreateCommentRequest, postId string) (model.Comment, error) {
	if postId == "" {
		return model.Comment{}, errors.New("postId is required")
	}
	if req.Body == "" {
		return model.Comment{}, errors.New("body is required")
	}
	return s.repo.CreateComment(ctx, req, postId)
}

func (s *commentServiceImpl) DeleteComment(ctx context.Context, postId string, id string) (model.Comment, error) {
	if postId == "" || id == "" {
		return model.Comment{}, errors.New("postId and id are required")
	}
	return s.repo.DeleteComment(ctx, postId, id)
}
