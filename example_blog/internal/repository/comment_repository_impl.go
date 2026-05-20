package repository

import (
	"context"
	"database/sql"

	"github.com/example/blog/internal/model"
)

type commentRepositoryImpl struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepositoryImpl{db: db}
}

func (r *commentRepositoryImpl) ListComments(ctx context.Context, req model.ListCommentsRequest, postId string) ([]model.Comment, error) {
	// TODO: implement SELECT WHERE post_id = $1
	panic("not implemented")
}

func (r *commentRepositoryImpl) CreateComment(ctx context.Context, req model.CreateCommentRequest, postId string) (model.Comment, error) {
	// TODO: implement INSERT RETURNING
	panic("not implemented")
}

func (r *commentRepositoryImpl) DeleteComment(ctx context.Context, postId string, id string) (model.Comment, error) {
	// TODO: implement DELETE WHERE post_id = $1 AND id = $2 RETURNING
	panic("not implemented")
}
