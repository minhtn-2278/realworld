package dto

import "realworldapp/internal/models"

type UserArticleResponse struct {
	ArticleID uint   `json:"article_id"`
	Title     string `json:"title"`
}

type UserResponse struct {
	Username string                `json:"username"`
	Bio      *string               `json:"bio"`
	Image    *string               `json:"image"`
	Articles []UserArticleResponse `json:"articles"`
}

func NewUserResponse(user *models.User) UserResponse {
	articles := make([]UserArticleResponse, 0, len(user.Articles))
	for _, article := range user.Articles {
		articles = append(articles, UserArticleResponse{
			ArticleID: article.ID,
			Title:     article.Title,
		})
	}

	return UserResponse{
		Username: user.Username,
		Bio:      user.Bio,
		Image:    user.Image,
		Articles: articles,
	}
}
