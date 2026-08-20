package handlers

import "github.com/labstack/echo/v5"

func RegisterRoutes(e *echo.Echo) {
	api := e.Group("/api")

	registerUserRoutes(api, NewUserHandler())
	registerArticleRoutes(api, NewArticleHandler())
	registerTagRoutes(api, NewTagHandler())
}

func registerUserRoutes(api *echo.Group, handler *UserHandler) {
	users := api.Group("/users")
	users.POST("", handler.Create)
}

func registerArticleRoutes(api *echo.Group, handler *ArticleHandler) {
	articles := api.Group("/articles")
	articles.POST("", handler.Create)
}

func registerTagRoutes(api *echo.Group, handler *TagHandler) {
	api.GET("/tags", handler.List)
}
