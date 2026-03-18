package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Pagination 独立分页结构体（可复用）
type Pagination struct {
	Limit  int32 `json:"limit"`  // 每页条数
	Offset int32 `json:"offset"` // 偏移量
}

type Comment struct {
	ID          uuid.UUID     `json:"id"`
	Author      uuid.UUID     `json:"author"`
	Content     string        `json:"content"`
	PostID      uuid.UUID     `json:"post_id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedType sql.NullInt16 `json:"deleted_type"`
	DeletedAt   sql.NullTime  `json:"deleted_at"`
	User        User          `json:"user"`
}

type Indicator struct {
	ID            uuid.UUID    `json:"id"`
	Indicator     string       `json:"indicator"`
	IndicatorType string       `json:"indicator_type"`
	MetaSource    string       `json:"meta_source"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	DeletedAt     sql.NullTime `json:"deleted_at"`
}

type Post struct {
	ID          uuid.UUID     `json:"id"`
	Author      uuid.UUID     `json:"author"`
	Title       string        `json:"title"`
	Content     string        `json:"content"`
	Pv          int32         `json:"pv"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedType sql.NullInt16 `json:"deleted_type"`
	DeletedAt   sql.NullTime  `json:"deleted_at"`
	User        User          `json:"user"`
}

type User struct {
	ID          uuid.UUID      `json:"id"`
	Email       string         `json:"email"`
	Name        string         `json:"name"`
	Password    string         `json:"password"`
	Avatar      string         `json:"avatar"`
	Gender      sql.NullString `json:"gender"`
	Bio         string         `json:"bio"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedType sql.NullInt16  `json:"deleted_type"`
	DeletedAt   sql.NullTime   `json:"deleted_at"`
}

type UserRepository interface {
	ListUsers(ctx context.Context, arg *User, pagination *Pagination) ([]*User, error)
	ListUsersForAdmin(ctx context.Context, arg *User, pagination *Pagination) ([]*User, error)
	CreateUser(ctx context.Context, arg *User) error
	UpdateUser(ctx context.Context, arg *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type PostRepository interface {
	ListPosts(ctx context.Context, arg *Post, pagination *Pagination) ([]*Post, error)
	ListPostsForAdmin(ctx context.Context, arg *Post, pagination *Pagination) ([]*Post, error)
	CreatePost(ctx context.Context, arg *Post) error
	UpdatePost(ctx context.Context, arg *Post) error
	DeletePost(ctx context.Context, id uuid.UUID) error
}

type CommentRepository interface {
	ListComments(ctx context.Context, arg *Comment) ([]*Comment, error)
	CreateComment(ctx context.Context, arg *Comment) error
	UpdateComment(ctx context.Context, arg *Comment) error
	DeleteComment(ctx context.Context, id uuid.UUID) error
}

type IndicatorRepository interface {
	ListIndicators(ctx context.Context, arg *Indicator, pagination *Pagination) ([]*Indicator, error)
	CreateIndicator(ctx context.Context, arg *Indicator) error
	UpdateIndicator(ctx context.Context, arg *Indicator) error
	DeleteIndicator(ctx context.Context, id uuid.UUID) error
}

type Repository interface {
	UserRepository
	PostRepository
	CommentRepository
	IndicatorRepository

	Close() error
}
