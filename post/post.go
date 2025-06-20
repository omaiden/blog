package post

import (
	"context"
	"errors"
	"time"

	"github.com/moonrhythm/randid"
	"github.com/moonrhythm/validator"

	"blog/internal/kctx"
	"blog/pkg/api"
	"blog/pkg/sql"
)

type Post struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	AuthorID  string `json:"author_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ListPostRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (r *ListPostRequest) Valid() error {
	v := validator.New()
	v.Must(r.Limit > 0 && r.Limit <= 100, "limit must be between 1 and 100")
	return v.Error()
}

type ListPostResponse struct {
	Items  []*Post
	Limit  int
	Offset int
	Total  int
}

type CreatePostRequest struct {
	Title   string
	Content string
}

func (r *CreatePostRequest) Valid() error {
	v := validator.New()
	v.Must(r.Title != "", "title required")
	v.Must(r.Content != "", "content required")
	return v.Error()
}

type UpdatePostRequest struct {
	ID      string
	Title   string
	Content string
}

func (r *UpdatePostRequest) Valid() error {
	v := validator.New()
	v.Must(r.ID != "", "id required")
	v.Must(r.Title != "", "title required")
	v.Must(r.Content != "", "content required")
	return v.Error()
}

type GetPostRequest struct {
	ID string `json:"id"`
}

type DeletePostRequest struct {
	ID string `json:"id"`
}

func ListPosts(ctx context.Context, req *ListPostRequest) (*ListPostResponse, error) {
	if err := req.Valid(); err != nil {
		return nil, err
	}

	currentUserID := kctx.GetUserID(ctx)
	if currentUserID == "" {
		return nil, errors.New("unauthorized")
	}

	posts, err := sql.ListPosts(ctx, &sql.ListPostFilter{
		AuthorID: currentUserID,
		Offset:   req.Offset,
		Limit:    req.Limit,
	})
	if err != nil {
		return nil, err
	}

	resp := make([]*Post, 0, len(posts))
	for _, p := range posts {
		resp = append(resp, convertToPostResp(p))
	}

	return &ListPostResponse{
		Items:  resp,
		Limit:  req.Limit,
		Offset: req.Offset,
		Total:  0, // TODO
	}, nil
}

func CreatePost(ctx context.Context, req *CreatePostRequest) (*Post, error) {
	if err := req.Valid(); err != nil {
		return nil, err
	}

	now := time.Now()
	currentUserID := kctx.GetUserID(ctx)

	sqlPost := sql.Post{
		ID:        randid.MustGenerate().String(),
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  currentUserID,
		Published: true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := sql.CreatePost(ctx, &sqlPost)
	if err != nil {
		return nil, err
	}

	return convertToPostResp(&sqlPost), nil
}

func GetPost(ctx context.Context, req *GetPostRequest) (*Post, error) {
	currentUserID := kctx.GetUserID(ctx)

	sqlPost, err := sql.GetPostByID(ctx, req.ID)
	if err != nil {
		return nil, api.WrapError(err)
	}

	if currentUserID != sqlPost.AuthorID {
		return nil, api.ErrRecordNotFound
	}

	return convertToPostResp(sqlPost), err
}

func UpdatePost(ctx context.Context, req *UpdatePostRequest) (*Post, error) {
	currentUserID := kctx.GetUserID(ctx)

	sqlPost, err := sql.GetPostByID(ctx, req.ID)
	if err != nil {
		return nil, api.WrapError(err)
	}

	if currentUserID != sqlPost.AuthorID {
		return nil, api.ErrRecordNotFound
	}

	sqlPost.Title = req.Title
	sqlPost.Content = req.Content
	sqlPost.UpdatedAt = time.Now()

	err = sql.UpdatePost(ctx, sqlPost)
	if err != nil {
		return nil, err
	}

	return convertToPostResp(sqlPost), nil

}

func DeletePost(ctx context.Context, req *DeletePostRequest) error {
	currentUserID := kctx.GetUserID(ctx)

	sqlPost, err := sql.GetPostByID(ctx, req.ID)
	if err != nil {
		return api.WrapError(err)
	}

	if currentUserID != sqlPost.AuthorID {
		return api.ErrRecordNotFound
	}

	sqlPost.Published = false
	sqlPost.UpdatedAt = time.Now()

	err = sql.UpdatePost(ctx, sqlPost)

	return err
}
