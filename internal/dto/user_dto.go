package dto

import "realworldapp/internal/models"

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=100"`
}

type RegisterUserResponse struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required,max=255"`
	Password string `json:"password" validate:"required,max=100"`
}

type LoginUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserArticleResponse struct {
	ArticleID uint   `json:"article_id"`
	Title     string `json:"title"`
}

type UserResponse struct {
	ID       uint                  `json:"id"`
	Username string                `json:"username"`
	Email    string                `json:"email"`
	Bio      *string               `json:"bio"`
	Image    *string               `json:"image"`
	Articles []UserArticleResponse `json:"articles"`
}

func NewRegisterUserResponse(user *models.User) RegisterUserResponse {
	return RegisterUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
		Image:    user.Image,
	}
}

func NewLoginUserResponse(accessToken string, refreshToken string) LoginUserResponse {
	return LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
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
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
		Image:    user.Image,
		Articles: articles,
	}
}
