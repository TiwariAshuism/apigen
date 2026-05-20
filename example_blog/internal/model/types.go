package model

// Post represents a blog post.
type Post struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	AuthorID  string `json:"author_id"`
	CreatedAt string `json:"created_at"`
}

// Comment represents a comment on a blog post.
type Comment struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	Body      string `json:"body"`
	AuthorID  string `json:"author_id"`
	CreatedAt string `json:"created_at"`
}

// ListPostsRequest holds query params for listing posts.
type ListPostsRequest struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	AuthorID string `form:"author_id"`
}

// CreatePostRequest holds the payload for creating a post.
type CreatePostRequest struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	AuthorID string `json:"author_id"`
}

// UpdatePostRequest holds the payload for updating a post.
type UpdatePostRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// ListCommentsRequest holds query params for listing comments on a post.
type ListCommentsRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

// CreateCommentRequest holds the payload for creating a comment.
type CreateCommentRequest struct {
	Body     string `json:"body"`
	AuthorID string `json:"author_id"`
}
