package handlers

import "github.com/labstack/echo/v5"

type ArticleHandler struct{}

func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

func (h *ArticleHandler) Create(c *echo.Context) error {
	return nil
}
