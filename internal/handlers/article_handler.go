package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"

	"realworldapp/internal/dto"
	"realworldapp/internal/services"
)

type ArticleHandler struct {
	articleService services.ArticleService
}

func NewArticleHandler(articleService services.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

func (h *ArticleHandler) Create(c *echo.Context) error {
	var request dto.CreateArticleRequest
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	article, err := h.articleService.Create(c.Request().Context(), services.CreateArticleInput{
		Title:       request.Title,
		Description: request.Description,
		Body:        request.Body,
		AuthorID:    request.AuthorID,
		TagList:     request.TagList,
	})
	if err != nil {
		return serviceError(err)
	}

	return c.JSON(http.StatusCreated, dto.ArticleResponseEnvelope{
		Article: dto.NewArticleResponse(article),
	})
}

func (h *ArticleHandler) List(c *echo.Context) error {
	articles, err := h.articleService.List(c.Request().Context())
	if err != nil {
		return serviceError(err)
	}

	return c.JSON(http.StatusOK, dto.NewArticleListResponse(articles))
}

func (h *ArticleHandler) GetBySlug(c *echo.Context) error {
	article, err := h.articleService.GetBySlug(c.Request().Context(), c.Param("slug"))
	if err != nil {
		return serviceError(err)
	}

	return c.JSON(http.StatusOK, dto.ArticleResponseEnvelope{
		Article: dto.NewArticleResponse(article),
	})
}

func (h *ArticleHandler) Update(c *echo.Context) error {
	var request dto.UpdateArticleRequest
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	article, err := h.articleService.Update(c.Request().Context(), c.Param("slug"), services.UpdateArticleInput{
		Title:       request.Title,
		Description: request.Description,
		Body:        request.Body,
		TagList:     request.TagList,
	})
	if err != nil {
		return serviceError(err)
	}

	return c.JSON(http.StatusOK, dto.ArticleResponseEnvelope{
		Article: dto.NewArticleResponse(article),
	})
}

func (h *ArticleHandler) Delete(c *echo.Context) error {
	if err := h.articleService.Delete(c.Request().Context(), c.Param("slug")); err != nil {
		return serviceError(err)
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{Message: "ok"})
}

func serviceError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "resource not found")
	}

	return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
}
