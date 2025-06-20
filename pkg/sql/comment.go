package sql

import (
	"context"
	"github.com/acoshift/pgsql/pgctx"
	"time"
)

type Comment struct {
	ID        string
	PostID    string
	AuthorID  string
	Content   string
	CreatedAt time.Time
}

type ListCommentFilter struct {
	PostID string
	Limit  int
	Offset int
}

func CreateComment(ctx context.Context, comment *Comment) error {
	return pgctx.QueryRow(ctx, `
        INSERT INTO comments (id, post_id, author_id, content, created_at) 
        	VALUES ($1, $2, $3, $4, $5)
    `, comment.ID, comment.PostID, comment.AuthorID, comment.Content, comment.CreatedAt).Err()
}

func ListComments(ctx context.Context, filter *ListCommentFilter) ([]*Comment, error) {
	rows, err := pgctx.Query(ctx, `
        SELECT *
        FROM comments 
        WHERE post_id = $1
        ORDER BY created_at DESC 
        LIMIT $2
        OFFSET $3
    `, filter.PostID, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		var a Comment
		err := rows.Scan(&a.ID, &a.PostID, &a.Content, &a.AuthorID, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &a)
	}

	return comments, nil
}
