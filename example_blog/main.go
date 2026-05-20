package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/example/blog/internal/handler"
	"github.com/example/blog/internal/repository"
	"github.com/example/blog/internal/service"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 dbname=blog sslmode=disable")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// Wire dependency chain: repository → service → handler
	postRepo := repository.NewPostRepository(db)
	postSvc := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postSvc)

	commentRepo := repository.NewCommentRepository(db)
	commentSvc := service.NewCommentService(commentRepo)
	commentHandler := handler.NewCommentHandler(commentSvc)

	r := gin.Default()
	v1 := r.Group("/v1")
	postHandler.RegisterRoutes(v1)
	commentHandler.RegisterRoutes(v1)

	log.Println("blog API listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("run: %v", err)
	}
}
