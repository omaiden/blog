package comment

import (
	"context"
	"errors"
	"time"

	"github.com/moonrhythm/randid"
	"github.com/moonrhythm/validator"

	"blog/internal/kctx"
	"blog/pkg/sql"
)

type Comment struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	AuthorID  string `json:"author_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type ListCommentRequest struct {
	PostID string `json:"post_id"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type ListCommentResponse struct {
	Items  []*Comment `json:"items"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
	Total  int        `json:"total"`
}

type CreateCommentRequest struct {
	PostID  string `json:"post_id"`
	Content string `json:"content"`
}

func (r *CreateCommentRequest) Valid() error {
	v := validator.New()
	v.Must(r.PostID != "", "post id required")
	v.Must(r.Content != "", "content required")
	return v.Error()
}

func GetComments(ctx context.Context, req *ListCommentRequest) (*ListCommentResponse, error) {
	currentUserID := kctx.GetUserID(ctx)
	if currentUserID == "" {
		return nil, errors.New("unauthorized")
	}

	comments, err := sql.ListComments(ctx, &sql.ListCommentFilter{
		PostID: req.PostID,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return nil, err
	}

	resp := make([]*Comment, 0, len(comments))
	for _, c := range comments {
		resp = append(resp, convertToCommentResp(c))
	}

	return &ListCommentResponse{
		Items:  resp,
		Limit:  req.Limit,
		Offset: req.Offset,
		Total:  0, // TODO
	}, nil
}

func CreateComment(ctx context.Context, req *CreateCommentRequest) (*Comment, error) {
	if err := req.Valid(); err != nil {
		return nil, err
	}

	now := time.Now()
	currentUserID := kctx.GetUserID(ctx)

	sqlComment := sql.Comment{
		ID:        randid.MustGenerate().String(),
		PostID:    req.PostID,
		Content:   req.Content,
		AuthorID:  currentUserID,
		CreatedAt: now,
	}
	err := sql.CreateComment(ctx, &sqlComment)
	if err != nil {
		return nil, err
	}

	return convertToCommentResp(&sqlComment), nil
}
