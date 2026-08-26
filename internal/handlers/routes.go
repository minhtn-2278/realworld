package handlers

import (
	"context"

	"github.com/labstack/echo/v5"

	appmiddleware "realworldapp/internal/middleware"
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
	articleOwnerMiddleware := newArticleOwnerMiddleware(articleService, authMiddleware)
	registerArticleRoutes(
		api,
		NewArticleHandler(articleService),
		authMiddleware,
		articleOwnerMiddleware,
	)
	registerTagRoutes(api, NewTagHandler(tagService))
}

func newArticleOwnerMiddleware(
	articleService services.ArticleService,
	authMiddleware echo.MiddlewareFunc,
) echo.MiddlewareFunc {
	return appmiddleware.RequireOwner(
		authMiddleware,
		func(ctx context.Context, c *echo.Context) (uint, error) {
			article, err := articleService.GetBySlug(ctx, c.Param("slug"))
			if err != nil {
				return 0, err
			}

			return article.AuthorID, nil
		},
	)
}

func registerUserRoutes(api *echo.Group, handler *UserHandler, authMiddleware echo.MiddlewareFunc) {
	users := api.Group("/users")
	users.POST("/login", handler.Login)
	users.POST("", handler.Create)
	api.GET("/user", handler.Current, authMiddleware)
	api.GET("/profiles/:username", handler.Profile, authMiddleware)
}

func registerArticleRoutes(
	api *echo.Group,
	handler *ArticleHandler,
	authMiddleware echo.MiddlewareFunc,
	ownerMiddleware echo.MiddlewareFunc,
) {
	articles := api.Group("/articles")
	articles.GET("/:slug/comments", handler.ListComments)
	articles.POST("/:slug/comments", handler.CreateComment, authMiddleware)
	articles.GET("", handler.List)
	articles.GET("/:slug", handler.GetBySlug)
	articles.POST("", handler.Create, authMiddleware)
	articles.PUT("/:slug", handler.Update, ownerMiddleware)
	articles.DELETE("/:slug", handler.Delete, ownerMiddleware)
}

func registerTagRoutes(api *echo.Group, handler *TagHandler) {
	api.GET("/tags", handler.List)
}
