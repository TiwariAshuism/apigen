package api

// Package api defines API interfaces for the blog application.
// Run go generate to regenerate handler/service/repository code.
//
//go:generate go run ../../cmd/apigen -input routes.go -output ..

import "context"

// ── Post resource ────────────────────────────────────────────────────────────

// PostAPI defines HTTP endpoints for blog posts.
type PostAPI interface {
	// GET /posts
	ListPosts(ctx context.Context, req ListPostsRequest) ([]Post, error)
	// GET /posts/:id
	GetPost(ctx context.Context, id string) (Post, error)
	// POST /posts
	CreatePost(ctx context.Context, req CreatePostRequest) (Post, error)
	// PUT /posts/:id
	UpdatePost(ctx context.Context, id string, req UpdatePostRequest) (Post, error)
	// DELETE /posts/:id
	DeletePost(ctx context.Context, id string) (Post, error)
}

// ── Comment resource ──────────────────────────────────────────────────────────

// CommentAPI defines HTTP endpoints for comments on a post.
type CommentAPI interface {
	// GET /posts/:postId/comments
	ListComments(ctx context.Context, postId string, req ListCommentsRequest) ([]Comment, error)
	// POST /posts/:postId/comments
	CreateComment(ctx context.Context, postId string, req CreateCommentRequest) (Comment, error)
	// DELETE /posts/:postId/comments/:id
	DeleteComment(ctx context.Context, postId string, id string) (Comment, error)
}

// ── Type stubs (used by the interface; real definitions live in internal/model) ─

type Post struct {
	ID        string
	Title     string
	Body      string
	AuthorID  string
	CreatedAt string
}

type Comment struct {
	ID        string
	PostID    string
	Body      string
	AuthorID  string
	CreatedAt string
}

type ListPostsRequest struct {
	Page     int
	Limit    int
	AuthorID string
}

type CreatePostRequest struct {
	Title    string
	Body     string
	AuthorID string
}

type UpdatePostRequest struct {
	Title string
	Body  string
}

type ListCommentsRequest struct {
	Page  int
	Limit int
}

type CreateCommentRequest struct {
	Body     string
	AuthorID string
}
