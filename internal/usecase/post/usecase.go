package post

import (
	"context"

	"github.com/google/uuid"
	"github.com/thkx/graphql-sqlc-go-dome/internal/auth"
	"github.com/thkx/graphql-sqlc-go-dome/internal/graph/model"
	"github.com/thkx/graphql-sqlc-go-dome/internal/repository"
	"github.com/thkx/graphql-sqlc-go-dome/internal/usecase"
)

type post struct {
	repo repository.PostRepository
}

func NewUsecase(repo repository.PostRepository) usecase.PostUsecase {
	return &post{repo: repo}
}

func (p *post) ListPosts(ctx context.Context) ([]*model.Post, error) {
	params := &repository.Post{}
	pagination := &repository.Pagination{
		Limit:  0,
		Offset: 10,
	}
	rows, err := p.repo.ListPosts(ctx, params, pagination)
	if err != nil {
		return nil, err
	}

	var result []*model.Post
	for _, v := range rows {
		ind := &model.Post{
			ID:        v.ID.String(),
			Title:     v.Title,
			Content:   v.Content,
			Pv:        v.Pv,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Author: &model.User{
				ID:     v.User.ID.String(),
				Name:   v.User.Name,
				Email:  v.User.Email,
				Avatar: v.User.Avatar,
				Gender: v.User.Gender.String,
				Bio:    v.User.Bio,
			},
		}

		result = append(result, ind)
	}
	return result, nil
}

func (p *post) CreatePost(ctx context.Context, input model.CreatePostInput) error {
	user_id, err := auth.RequireAuth(ctx)
	if err != nil {
		return err
	}

	post := &repository.Post{
		ID:      uuid.New(),
		Title:   input.Title,
		Content: input.Content,
		Pv:      0,
		Author:  uuid.MustParse(user_id),
	}

	return p.repo.CreatePost(ctx, post)
}

func (p *post) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	params := &repository.Post{
		ID: uuid.MustParse(id),
	}
	pagination := &repository.Pagination{
		Limit:  0,
		Offset: 10,
	}
	rows, err := p.repo.ListPosts(ctx, params, pagination)
	if err != nil {
		return nil, err
	}

	var result []*model.Post
	for _, v := range rows {
		ind := &model.Post{
			ID:        v.ID.String(),
			Title:     v.Title,
			Content:   v.Content,
			Pv:        v.Pv,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Author: &model.User{
				ID:     v.User.ID.String(),
				Name:   v.User.Name,
				Email:  v.User.Email,
				Avatar: v.User.Avatar,
				Gender: v.User.Gender.String,
				Bio:    v.User.Bio,
			},
		}

		result = append(result, ind)
	}
	return result[0], nil
}

func (p *post) ListPostsByTitle(ctx context.Context, title string) ([]*model.Post, error) {
	params := &repository.Post{
		Title: title,
	}
	pagination := &repository.Pagination{
		Limit:  0,
		Offset: 10,
	}
	rows, err := p.repo.ListPosts(ctx, params, pagination)
	if err != nil {
		return nil, err
	}

	var result []*model.Post
	for _, v := range rows {
		ind := &model.Post{
			ID:        v.ID.String(),
			Title:     v.Title,
			Content:   v.Content,
			Pv:        v.Pv,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Author: &model.User{
				ID:     v.User.ID.String(),
				Name:   v.User.Name,
				Email:  v.User.Email,
				Avatar: v.User.Avatar,
				Gender: v.User.Gender.String,
				Bio:    v.User.Bio,
			},
		}

		result = append(result, ind)
	}
	return result, nil
}

func (p *post) UpdatePost(ctx context.Context, input *model.UpdatePostInput) (*model.Post, error) {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return nil, err
	}

	params := &repository.Post{
		ID:      uuid.MustParse(input.ID),
		Title:   *input.Title,
		Content: *input.Content,
		Pv:      *input.Pv,
	}
	if err := p.repo.UpdatePost(ctx, params); err != nil {
		return nil, err
	}
	return nil, nil
}

func (p *post) DeletePost(ctx context.Context, id string) error {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return err
	}

	if err := p.repo.DeletePost(ctx, uuid.MustParse(id)); err != nil {
		return err
	}

	return nil
}
