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
	authMiddleware echo.MiddlewareFunc,
) {
	api := e.Group("/api")

	registerUserRoutes(api, NewUserHandler(userService), authMiddleware)
	registerArticleRoutes(api, NewArticleHandler(articleService), authMiddleware)
	registerTagRoutes(api, NewTagHandler(tagService))
}

func registerUserRoutes(api *echo.Group, handler *UserHandler, authMiddleware echo.MiddlewareFunc) {
	users := api.Group("/users")
	users.POST("/login", handler.Login)
	users.POST("", handler.Create)
	api.GET("/user", handler.Current, authMiddleware)
	api.GET("/profiles/:username", handler.Profile, authMiddleware)
}

func registerArticleRoutes(api *echo.Group, handler *ArticleHandler, authMiddleware echo.MiddlewareFunc) {
	articles := api.Group("/articles")
	articles.GET("/:slug/comments", handler.ListComments)
	articles.POST("/:slug/comments", handler.CreateComment, authMiddleware)
	articles.GET("", handler.List)
	articles.GET("/:slug", handler.GetBySlug)
	articles.POST("", handler.Create, authMiddleware)
	articles.PUT("/:slug", handler.Update, authMiddleware)
	articles.DELETE("/:slug", handler.Delete, authMiddleware)
}

func registerTagRoutes(api *echo.Group, handler *TagHandler) {
	api.GET("/tags", handler.List)
}
