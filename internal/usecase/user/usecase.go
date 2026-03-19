package user

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/thkx/graphql-sqlc-go-dome/internal/auth"
	"github.com/thkx/graphql-sqlc-go-dome/internal/auth/jwt"
	"github.com/thkx/graphql-sqlc-go-dome/internal/graph/model"
	"github.com/thkx/graphql-sqlc-go-dome/internal/repository"
	"github.com/thkx/graphql-sqlc-go-dome/internal/usecase"
)

type user struct {
	repo repository.UserRepository
}

func NewUsecase(repo repository.UserRepository) usecase.UserUsecase {
	return &user{repo: repo}
}

func (u *user) ListUsers(ctx context.Context) ([]*model.User, error) {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return nil, err
	}
	params := &repository.User{}
	pagination := &repository.Pagination{
		Limit:  0,
		Offset: 10,
	}
	rows, err := u.repo.ListUsers(ctx, params, pagination)
	if err != nil {
		return nil, err
	}
	var result []*model.User
	for _, v := range rows {
		row := &model.User{
			ID:        v.ID.String(),
			Name:      v.Name,
			Email:     v.Email,
			Avatar:    v.Avatar,
			Gender:    v.Gender.String,
			Bio:       v.Bio,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		result = append(result, row)
	}
	return result, nil
}

func (u *user) GetUserByID(ctx context.Context, input string) (*model.User, error) {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return nil, err
	}
	params := &repository.User{
		ID: uuid.MustParse(input),
	}
	pagination := &repository.Pagination{
		Limit:  0,
		Offset: 10,
	}
	rows, err := u.repo.ListUsers(ctx, params, pagination)
	if err != nil {
		return nil, err
	}
	var result []*model.User
	for _, v := range rows {
		row := &model.User{
			ID:        v.ID.String(),
			Name:      v.Name,
			Email:     v.Email,
			Avatar:    v.Avatar,
			Gender:    v.Gender.String,
			Bio:       v.Bio,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		result = append(result, row)
	}
	return result[0], nil
}

func (u *user) CreateUser(ctx context.Context, input model.CreateUserInput) error {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return err
	}
	params := &repository.User{
		ID:       uuid.New(),
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := u.repo.CreateUser(ctx, params); err != nil {
		return err
	}
	return nil
}

func (u *user) GetUserByEmail(ctx context.Context, input model.UserEmailInput) (*model.User, error) {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return nil, err
	}
	params := &repository.User{
		Email: input.Email,
	}
	pagination := &repository.Pagination{
		Limit:  0,
		Offset: 10,
	}
	rows, err := u.repo.ListUsers(ctx, params, pagination)
	if err != nil {
		return nil, err
	}
	var result []*model.User
	for _, v := range rows {
		row := &model.User{
			ID:        v.ID.String(),
			Name:      v.Name,
			Email:     v.Email,
			Avatar:    v.Avatar,
			Gender:    v.Gender.String,
			Bio:       v.Bio,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		result = append(result, row)
	}
	return result[0], nil
}

func (u *user) GetUserByName(ctx context.Context, input model.UserNameInput) (*model.User, error) {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return nil, err
	}
	params := &repository.User{
		Name: input.Name,
	}
	pagination := &repository.Pagination{
		Limit:  0,
		Offset: 10,
	}
	rows, err := u.repo.ListUsers(ctx, params, pagination)
	if err != nil {
		return nil, err
	}
	var result []*model.User
	for _, v := range rows {
		row := &model.User{
			ID:        v.ID.String(),
			Name:      v.Name,
			Email:     v.Email,
			Avatar:    v.Avatar,
			Gender:    v.Gender.String,
			Bio:       v.Bio,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		result = append(result, row)
	}
	return result[0], nil
}

func (u *user) UpdateUser(ctx context.Context, input model.UpdateUserInput) (*model.User, error) {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return nil, err
	}

	params := &repository.User{
		ID:       uuid.MustParse(input.ID),
		Name:     *input.Name,
		Email:    *input.Email,
		Password: *input.Password,
		Avatar:   *input.Avatar,
		Bio:      *input.Bio,
		Gender: sql.NullString{
			String: *input.Gender,
			Valid:  *input.Gender != "",
		},
	}
	if err := u.repo.UpdateUser(ctx, params); err != nil {
		return nil, err
	}
	return nil, nil
}

func (u *user) DeleteUser(ctx context.Context, input string) error {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return err
	}

	return u.repo.DeleteUser(ctx, uuid.MustParse(input))
}

func (u *user) Login(ctx context.Context, input model.LoginInput, w http.ResponseWriter) (*model.Login, error) {
	user, err := u.GetUserByEmail(ctx, model.UserEmailInput{Email: input.Email})
	if err != nil || user.Password != input.Password {
		return nil, errors.New("invalid email or password")
	}
	token, err := jwt.GenerateToken(user.ID, user.Name, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth-cookie",
		Value:    url.QueryEscape(token),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   86400,
		// Secure:   os.Getenv("ENV") == "production", // 生产开启
	})
	return &model.Login{Token: token, User: user}, nil
}
