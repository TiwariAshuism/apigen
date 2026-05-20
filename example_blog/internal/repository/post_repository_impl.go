package repository

import (
	"context"
	"database/sql"

	"github.com/example/blog/internal/model"
)

type postRepositoryImpl struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepositoryImpl{db: db}
}

func (r *postRepositoryImpl) ListPosts(ctx context.Context, req model.ListPostsRequest) ([]model.Post, error) {
	// TODO: implement SELECT with optional author_id filter
	panic("not implemented")
}

func (r *postRepositoryImpl) GetPost(ctx context.Context, id string) (model.Post, error) {
	// TODO: implement SELECT WHERE id = $1
	panic("not implemented")
}

func (r *postRepositoryImpl) CreatePost(ctx context.Context, req model.CreatePostRequest) (model.Post, error) {
	// TODO: implement INSERT RETURNING
	panic("not implemented")
}

func (r *postRepositoryImpl) UpdatePost(ctx context.Context, req model.UpdatePostRequest, id string) (model.Post, error) {
	// TODO: implement UPDATE WHERE id = $1 RETURNING
	panic("not implemented")
}

func (r *postRepositoryImpl) DeletePost(ctx context.Context, id string) (model.Post, error) {
	// TODO: implement DELETE WHERE id = $1 RETURNING
	panic("not implemented")
}
