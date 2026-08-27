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

type FollowUserRequest struct {
	Follow *bool `json:"follow" validate:"required"`
}

type UserArticleResponse struct {
	ArticleID uint   `json:"article_id"`
	Title     string `json:"title"`
}

type UserResponse struct {
	ID             uint                  `json:"id"`
	Username       string                `json:"username"`
	Email          string                `json:"email"`
	Bio            *string               `json:"bio"`
	Image          *string               `json:"image"`
	Following      bool                  `json:"following"`
	FollowersCount int64                 `json:"followersCount"`
	FollowingCount int64                 `json:"followingCount"`
	Articles       []UserArticleResponse `json:"articles"`
}

type CurrentUserResponse struct {
	ID             uint                  `json:"id"`
	Username       string                `json:"username"`
	Email          string                `json:"email"`
	Bio            *string               `json:"bio"`
	Image          *string               `json:"image"`
	FollowersCount int64                 `json:"followersCount"`
	FollowingCount int64                 `json:"followingCount"`
	Articles       []UserArticleResponse `json:"articles"`
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
	articles := newUserArticleResponses(user)

	return UserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Bio:            user.Bio,
		Image:          user.Image,
		Following:      user.Following,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
		Articles:       articles,
	}
}

func NewCurrentUserResponse(user *models.User) CurrentUserResponse {
	return CurrentUserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Bio:            user.Bio,
		Image:          user.Image,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
		Articles:       newUserArticleResponses(user),
	}
}

func newUserArticleResponses(user *models.User) []UserArticleResponse {
	articles := make([]UserArticleResponse, 0, len(user.Articles))
	for _, article := range user.Articles {
		articles = append(articles, UserArticleResponse{
			ArticleID: article.ID,
			Title:     article.Title,
		})
	}

	return articles
}
