package sql

import (
	"context"
	"github.com/acoshift/pgsql/pgctx"
	"time"
)

type Post struct {
	ID        string
	Title     string
	Content   string
	AuthorID  string
	Published bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListPostFilter struct {
	AuthorID string
	Offset   int
	Limit    int
}

func CreatePost(ctx context.Context, post *Post) error {
	return pgctx.QueryRow(ctx, `
        INSERT INTO posts (id, title, content, author_id, published, created_at, updated_at) 
        	VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, post.ID, post.Title, post.Content, post.AuthorID, post.Published, post.CreatedAt, post.UpdatedAt).Err()
}

func UpdatePost(ctx context.Context, post *Post) error {
	return pgctx.QueryRow(ctx, `
        UPDATE posts 
        SET title = $1 , content = $2, published = $3, updated_at = $4 
        WHERE id = $5
    `, post.Title, post.Content, post.Published, post.UpdatedAt, post.ID).Err()
}

func GetPostByID(ctx context.Context, postID string) (*Post, error) {
	var post Post
	err := pgctx.QueryRow(ctx, `
	SELECT *
		FROM posts
		WHERE id = $1 AND published = true
		`, postID).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.Published,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	return &post, err
}

func ListPosts(ctx context.Context, filter *ListPostFilter) ([]*Post, error) {
	rows, err := pgctx.Query(ctx, `
        SELECT *
        FROM posts 
        WHERE author_id = $1 AND published = true
        ORDER BY created_at DESC 
        LIMIT $2
        OFFSET $3
    `, filter.AuthorID, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var a Post
		err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.AuthorID, &a.Published, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &a)
	}

	return posts, nil
}
