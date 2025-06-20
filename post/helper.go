package post

import (
	"blog/pkg/api"
	"blog/pkg/sql"
)

func convertToPostResp(sqlPost *sql.Post) *Post {
	if sqlPost == nil {
		return nil
	}

	return &Post{
		ID:        sqlPost.ID,
		Title:     sqlPost.Title,
		Content:   sqlPost.Content,
		AuthorID:  sqlPost.AuthorID,
		CreatedAt: api.ConvertTimeToStr(sqlPost.CreatedAt),
		UpdatedAt: api.ConvertTimeToStr(sqlPost.UpdatedAt),
	}
}
