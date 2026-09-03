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

// Create creates an article for the authenticated user.
// @Summary Create an article
// @Tags articles
// @Accept json
// @Produce json
// @Param request body dto.CreateArticleRequest true "Article payload"
// @Success 201 {object} dto.ArticleResponseEnvelope
// @Failure 400 {object} utils.ValidationError
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 409 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Security BearerAuth
// @Router /articles [post]
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

// List lists articles with optional filters and pagination.
// @Summary List articles
// @Tags articles
// @Produce json
// @Param tag query string false "Tag name"
// @Param author query string false "Author username"
// @Param limit query int false "Page size" default(10)
// @Param page query int false "Page number" default(1)
// @Success 200 {object} dto.ArticleListResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /articles [get]
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

// ListFeed lists articles from users followed by the authenticated user.
// @Summary List the authenticated user's feed
// @Tags articles
// @Produce json
// @Param page query int false "Page number" default(1)
// @Success 200 {object} dto.ArticleListResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Security BearerAuth
// @Router /articles/feed [get]
func (h *ArticleHandler) ListFeed(c *echo.Context) error {
	pagination, err := utils.ParsePagination(c)
	if err != nil {
		return err
	}
	pagination.Limit = utils.DefaultPaginationLimit
	userID, err := appmiddleware.CurrentUserID(c)
	if err != nil {
		return err
	}

	articles, err := h.articleService.ListFeed(c.Request().Context(), userID, pagination)
	if err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.NewArticleListResponse(
		articles.Articles,
		articles.FavoriteCounts,
		pagination.Metadata(articles.Total),
	))
}

// GetBySlug returns an article and its comments.
// @Summary Get an article by slug
// @Tags articles
// @Produce json
// @Param slug path string true "Article slug"
// @Success 200 {object} dto.ArticleDetailResponseEnvelope
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /articles/{slug} [get]
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

// Favorite creates or removes an article favorite for the authenticated user.
// @Summary Favorite or unfavorite an article
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Param request body dto.FavoriteArticleRequest true "Favorite payload"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} utils.ValidationError
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Security BearerAuth
// @Router /articles/{slug}/favorite [post]
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

// Update updates an article owned by the authenticated user.
// @Summary Update an article
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Param request body dto.UpdateArticleRequest true "Article payload"
// @Success 200 {object} dto.ArticleResponseEnvelope
// @Failure 400 {object} utils.ValidationError
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 403 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 409 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Security BearerAuth
// @Router /articles/{slug} [put]
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

// Delete deletes an article owned by the authenticated user.
// @Summary Delete an article
// @Tags articles
// @Produce json
// @Param slug path string true "Article slug"
// @Success 200 {object} dto.MessageResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 403 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Security BearerAuth
// @Router /articles/{slug} [delete]
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

// ListComments lists comments for an article.
// @Summary List article comments
// @Tags comments
// @Produce json
// @Param slug path string true "Article slug"
// @Success 200 {object} dto.CommentListResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /articles/{slug}/comments [get]
func (h *ArticleHandler) ListComments(c *echo.Context) error {
	comments, err := h.articleService.ListComments(c.Request().Context(), c.Param("slug"))
	if err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.NewCommentListResponse(comments))
}

// CreateComment creates a comment on an article.
// @Summary Create an article comment
// @Tags comments
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Param request body dto.CreateCommentRequest true "Comment payload"
// @Success 201 {object} dto.CommentResponse
// @Failure 400 {object} utils.ValidationError
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Security BearerAuth
// @Router /articles/{slug}/comments [post]
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

// DeleteComment deletes a comment owned by the authenticated user.
// @Summary Delete an article comment
// @Tags comments
// @Produce json
// @Param slug path string true "Article slug"
// @Param id path int true "Comment ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 403 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Security BearerAuth
// @Router /articles/{slug}/comments/{id} [delete]
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
