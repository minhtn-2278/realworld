package handlers

import (
	"github.com/labstack/echo/v5"

	"realworldapp/internal/services"
)

func RegisterRoutes(
	e *echo.Echo,
	articleService services.ArticleService,
	tagService services.TagService,
	userService services.UserService,
) {
	api := e.Group("/api")

	registerUserRoutes(api, NewUserHandler(userService))
	registerArticleRoutes(api, NewArticleHandler(articleService))
	registerTagRoutes(api, NewTagHandler(tagService))
}

func registerUserRoutes(api *echo.Group, handler *UserHandler) {
	users := api.Group("/users")
	users.POST("", handler.Create)
	api.GET("/profiles/:username", handler.Profile)
}

func registerArticleRoutes(api *echo.Group, handler *ArticleHandler) {
	articles := api.Group("/articles")
	articles.GET("/:slug/comments", handler.ListComments)
	articles.POST("/:slug/comments", handler.CreateComment)
	articles.GET("", handler.List)
	articles.GET("/:slug", handler.GetBySlug)
	articles.POST("", handler.Create)
	articles.PUT("/:slug", handler.Update)
	articles.DELETE("/:slug", handler.Delete)
}

func registerTagRoutes(api *echo.Group, handler *TagHandler) {
	api.GET("/tags", handler.List)
}
