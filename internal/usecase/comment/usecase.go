package comment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/thkx/graphql-sqlc-go-dome/internal/auth"
	"github.com/thkx/graphql-sqlc-go-dome/internal/graph/model"
	"github.com/thkx/graphql-sqlc-go-dome/internal/repository"
	"github.com/thkx/graphql-sqlc-go-dome/internal/usecase"
)

type comment struct {
	repo repository.CommentRepository
}

func NewUsecase(repo repository.CommentRepository) usecase.CommentUsecase {
	return &comment{repo: repo}
}

func (c *comment) ListComments(ctx context.Context) ([]*model.Comment, error) {
	params := &repository.Comment{}
	rows, err := c.repo.ListComments(ctx, params)
	if err != nil {
		return nil, err
	}
	var result []*model.Comment
	for _, row := range rows {
		com := &model.Comment{
			ID:        row.ID.String(),
			Content:   row.Content,
			PostID:    row.PostID.String(),
			CreatedAt: row.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: row.UpdatedAt.Format("2006-01-02 15:04:05"),
			Author: &model.User{
				ID:     row.User.ID.String(),
				Name:   row.User.Name,
				Email:  row.User.Email,
				Avatar: row.User.Avatar,
				Gender: row.User.Gender.String,
				Bio:    row.User.Bio,
			},
		}
		result = append(result, com)
	}
	return result, nil
}

func (c *comment) GetCommentByID(ctx context.Context, id string) (*model.Comment, error) {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := c.repo.ListComments(ctx, &repository.Comment{ID: uuid.MustParse(id)})
	if err != nil {
		return nil, err
	}

	if rows == nil {
		return nil, fmt.Errorf("not data")
	}

	row := &model.Comment{
		ID:        rows[0].ID.String(),
		Content:   rows[0].Content,
		PostID:    rows[0].PostID.String(),
		CreatedAt: rows[0].CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: rows[0].UpdatedAt.Format("2006-01-02 15:04:05"),
		Author: &model.User{
			ID:     rows[0].User.ID.String(),
			Name:   rows[0].User.Name,
			Email:  rows[0].User.Email,
			Avatar: rows[0].User.Avatar,
			Gender: rows[0].User.Gender.String,
			Bio:    rows[0].User.Bio,
		},
	}

	return row, nil
}

func (c *comment) ListCommentsByUserID(ctx context.Context, userID string) ([]*model.Comment, error) {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := c.repo.ListComments(ctx, &repository.Comment{Author: uuid.MustParse(userID)})
	if err != nil {
		return nil, err
	}

	var result []*model.Comment

	for _, row := range rows {
		com := &model.Comment{
			ID:        row.ID.String(),
			Content:   row.Content,
			PostID:    row.PostID.String(),
			CreatedAt: row.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: row.UpdatedAt.Format("2006-01-02 15:04:05"),
			Author: &model.User{
				ID:     row.User.ID.String(),
				Name:   row.User.Name,
				Email:  row.User.Email,
				Avatar: row.User.Avatar,
				Gender: row.User.Gender.String,
				Bio:    row.User.Bio,
			},
		}

		result = append(result, com)
	}

	return result, nil
}

func (c *comment) ListCommentsByPostID(ctx context.Context, postID string) ([]*model.Comment, error) {
	rows, err := c.repo.ListComments(ctx, &repository.Comment{PostID: uuid.MustParse(postID)})
	if err != nil {
		return nil, err
	}

	var result []*model.Comment

	for _, row := range rows {
		com := &model.Comment{
			ID:        row.ID.String(),
			Content:   row.Content,
			PostID:    row.PostID.String(),
			CreatedAt: row.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: row.UpdatedAt.Format("2006-01-02 15:04:05"),
			Author: &model.User{
				ID:     row.User.ID.String(),
				Name:   row.User.Name,
				Email:  row.User.Email,
				Avatar: row.User.Avatar,
				Gender: row.User.Gender.String,
				Bio:    row.User.Bio,
			},
		}

		result = append(result, com)
	}

	return result, nil
}

func (c *comment) ListCommentsByContent(ctx context.Context, content string) ([]*model.Comment, error) {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := c.repo.ListComments(ctx, &repository.Comment{Content: content})
	if err != nil {
		return nil, err
	}

	var result []*model.Comment

	for _, row := range rows {
		com := &model.Comment{
			ID:        row.ID.String(),
			Content:   row.Content,
			PostID:    row.PostID.String(),
			CreatedAt: row.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: row.UpdatedAt.Format("2006-01-02 15:04:05"),
			Author: &model.User{
				ID:     row.User.ID.String(),
				Name:   row.User.Name,
				Email:  row.User.Email,
				Avatar: row.User.Avatar,
				Gender: row.User.Gender.String,
				Bio:    row.User.Bio,
			},
		}

		result = append(result, com)
	}

	return result, nil
}

func (c *comment) CreateComment(ctx context.Context, input model.CreateCommentInput) error {
	userID, err := auth.RequireAuth(ctx)
	if err != nil {
		return err
	}

	params := &repository.Comment{
		ID:      uuid.New(),
		Content: input.Content,
		PostID:  uuid.MustParse(input.PostID),
		Author:  uuid.MustParse(userID),
	}

	return c.repo.CreateComment(ctx, params)
}

func (c *comment) UpdateComment(ctx context.Context, input model.UpdateCommentInput) (*model.Comment, error) {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return nil, err
	}
	params := &repository.Comment{
		ID:      uuid.MustParse(input.ID),
		Content: input.Content,
	}

	if err := c.repo.UpdateComment(ctx, params); err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *comment) DeleteComment(ctx context.Context, id string) error {
	return c.repo.DeleteComment(ctx, uuid.MustParse(id))
}
