package dto

import (
	"time"

	"realworldapp/internal/models"
)

type CreateCommentRequest struct {
	Body string `json:"body" validate:"required"`
}

type CommentAuthorResponse struct {
	AuthorID uint   `json:"authorId"`
	Username string `json:"username"`
}

type CommentResponse struct {
	ID        uint                  `json:"id"`
	Body      string                `json:"body"`
	Author    CommentAuthorResponse `json:"author"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
}

type CommentListResponse struct {
	Comments []CommentResponse `json:"comments"`
}

func NewCommentResponse(comment *models.Comment) CommentResponse {
	return CommentResponse{
		ID:   comment.ID,
		Body: comment.Body,
		Author: CommentAuthorResponse{
			AuthorID: comment.Author.ID,
			Username: comment.Author.Username,
		},
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}

func NewCommentListResponse(comments []models.Comment) CommentListResponse {
	items := make([]CommentResponse, 0, len(comments))
	for index := range comments {
		items = append(items, NewCommentResponse(&comments[index]))
	}

	return CommentListResponse{Comments: items}
}
