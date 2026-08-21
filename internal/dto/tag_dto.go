package dto

import "realworldapp/internal/models"

// TagListResponse is the public JSON representation of the tag list.
type TagListResponse struct {
	Tags []string `json:"tags"`
}

// NewTagListResponse maps tag models to a tag list response DTO.
func NewTagListResponse(tags []models.Tag) TagListResponse {
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		names = append(names, tag.Name)
	}

	return TagListResponse{Tags: names}
}
