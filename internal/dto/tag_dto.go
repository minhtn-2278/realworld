package dto

import "realworldapp/internal/models"

type TagListResponse struct {
	Tags []string `json:"tags"`
}

func NewTagListResponse(tags []models.Tag) TagListResponse {
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		names = append(names, tag.Name)
	}

	return TagListResponse{Tags: names}
}
