package comment

import (
	"blog/pkg/api"
	"blog/pkg/sql"
)

func convertToCommentResp(sqlComment *sql.Comment) *Comment {
	if sqlComment == nil {
		return nil
	}

	return &Comment{
		ID:        sqlComment.ID,
		PostID:    sqlComment.PostID,
		AuthorID:  sqlComment.AuthorID,
		Content:   sqlComment.Content,
		CreatedAt: api.ConvertTimeToStr(sqlComment.CreatedAt),
	}
}
