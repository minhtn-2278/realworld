package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"

	"realworldapp/internal/dto"
	appmiddleware "realworldapp/internal/middleware"
	"realworldapp/internal/services"
	"realworldapp/internal/utils"
)

type ArticleHandler struct {
	articleService services.ArticleService
}

func NewArticleHandler(articleService services.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

func (h *ArticleHandler) Create(c *echo.Context) error {
	var request dto.CreateArticleRequest
	if err := utils.BindAndValidate(c, &request); err != nil {
		return err
	}
	authorID, err := appmiddleware.CurrentUserID(c)
	if err != nil {
		return err
	}

	article, err := h.articleService.Create(c.Request().Context(), services.CreateArticleInput{
		Title:       request.Title,
		Description: request.Description,
		Body:        request.Body,
		AuthorID:    authorID,
		TagList:     request.TagList,
	})
	if err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusCreated, dto.ArticleResponseEnvelope{
		Article: dto.NewArticleResponse(article),
	})
}

func (h *ArticleHandler) List(c *echo.Context) error {
	pagination, err := utils.ParsePagination(c)
	if err != nil {
		return err
	}

	tag := strings.TrimSpace(c.QueryParam("tag"))
	author := strings.TrimSpace(c.QueryParam("author"))

	articles, err := h.articleService.List(c.Request().Context(), services.ListArticlesInput{
		Tag:        tag,
		Author:     author,
		Pagination: pagination,
	})
	if err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.NewArticleListResponse(
		articles.Articles,
		articles.FavoriteCounts,
		pagination.Metadata(articles.Total),
	))
}

func (h *ArticleHandler) GetBySlug(c *echo.Context) error {
	detail, err := h.articleService.GetDetail(c.Request().Context(), c.Param("slug"))
	if err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.ArticleDetailResponseEnvelope{
		Article: dto.NewArticleDetailResponse(
			detail.Article,
			detail.Comments,
			detail.FavoritesCount,
		),
	})
}

func (h *ArticleHandler) Favorite(c *echo.Context) error {
	var request dto.FavoriteArticleRequest
	if err := utils.BindAndValidate(c, &request); err != nil {
		return err
	}
	userID, err := appmiddleware.CurrentUserID(c)
	if err != nil {
		return err
	}

	if err := h.articleService.Favorite(c.Request().Context(), c.Param("slug"), userID, *request.Favorite); err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{Message: "ok"})
}

func (h *ArticleHandler) Update(c *echo.Context) error {
	var request dto.UpdateArticleRequest
	if err := utils.BindAndValidate(c, &request); err != nil {
		return err
	}
	authorID, err := appmiddleware.CurrentUserID(c)
	if err != nil {
		return err
	}

	article, err := h.articleService.Update(
		c.Request().Context(),
		c.Param("slug"),
		authorID,
		services.UpdateArticleInput{
			Title:       request.Title,
			Description: request.Description,
			Body:        request.Body,
			TagList:     request.TagList,
		},
	)
	if err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.ArticleResponseEnvelope{
		Article: dto.NewArticleResponse(article),
	})
}

func (h *ArticleHandler) Delete(c *echo.Context) error {
	authorID, err := appmiddleware.CurrentUserID(c)
	if err != nil {
		return err
	}

	if err := h.articleService.Delete(c.Request().Context(), c.Param("slug"), authorID); err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{Message: "ok"})
}

func (h *ArticleHandler) ListComments(c *echo.Context) error {
	comments, err := h.articleService.ListComments(c.Request().Context(), c.Param("slug"))
	if err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.NewCommentListResponse(comments))
}

func (h *ArticleHandler) CreateComment(c *echo.Context) error {
	var request dto.CreateCommentRequest
	if err := utils.BindAndValidate(c, &request); err != nil {
		return err
	}
	authorID, err := appmiddleware.CurrentUserID(c)
	if err != nil {
		return err
	}

	comment, err := h.articleService.CreateComment(
		c.Request().Context(),
		c.Param("slug"),
		request.Body,
		authorID,
	)
	if err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusCreated, dto.NewCommentResponse(comment))
}

func (h *ArticleHandler) DeleteComment(c *echo.Context) error {
	commentID, err := parseCommentID(c.Param("id"))
	if err != nil {
		return err
	}
	authorID, err := appmiddleware.CurrentUserID(c)
	if err != nil {
		return err
	}

	if err := h.articleService.DeleteComment(c.Request().Context(), c.Param("slug"), commentID, authorID); err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{Message: "ok"})
}

func parseCommentID(value string) (uint, error) {
	commentID, err := strconv.ParseUint(value, 10, 0)
	if err != nil || commentID == 0 {
		return 0, utils.APIError(http.StatusBadRequest, "comment id must be a positive integer")
	}

	return uint(commentID), nil
}
