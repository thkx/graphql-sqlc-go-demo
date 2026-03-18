package usecase

import (
	"context"
	"net/http"

	"github.com/thkx/graphql-sqlc-go-dome/internal/graph/model"
)

type Usecase interface {
	UserUsecase
	PostUsecase
	IndicatorUsecase
}

// 可单独定义子接口
type UserUsecase interface {
	ListUsers(ctx context.Context) ([]*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, input model.CreateUserInput) error
	DeleteUser(ctx context.Context, id string) error
	Login(ctx context.Context, input model.LoginInput, w http.ResponseWriter) (*model.Login, error)
}

type PostUsecase interface {
	ListPosts(ctx context.Context) ([]*model.Post, error)
	CreatePost(ctx context.Context, input model.CreatePostInput) error
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
	ListPostsByTitle(ctx context.Context, title string) ([]*model.Post, error)
	UpdatePost(ctx context.Context, input *model.UpdatePostInput) (*model.Post, error)
	DeletePost(ctx context.Context, id string) error
}

type CommentUsecase interface {
	ListComments(ctx context.Context) ([]*model.Comment, error)
	GetCommentByID(ctx context.Context, id string) (*model.Comment, error)
	ListCommentsByUserID(ctx context.Context, userID string) ([]*model.Comment, error)
	ListCommentsByPostID(ctx context.Context, postID string) ([]*model.Comment, error)
	ListCommentsByContent(ctx context.Context, content string) ([]*model.Comment, error)
	CreateComment(ctx context.Context, input model.CreateCommentInput) error
	UpdateComment(ctx context.Context, input model.UpdateCommentInput) (*model.Comment, error)
	DeleteComment(ctx context.Context, id string) error
}

type IndicatorUsecase interface {
	CreateIndicator(ctx context.Context, input model.IndicatorInput) error
}

type compositeUsecase struct {
	UserUsecase
	PostUsecase
	CommentUsecase
	IndicatorUsecase
}

func NewUsecase(user UserUsecase, post PostUsecase, comment CommentUsecase, indicator IndicatorUsecase) Usecase {
	return &compositeUsecase{UserUsecase: user, PostUsecase: post, CommentUsecase: comment, IndicatorUsecase: indicator}
}
