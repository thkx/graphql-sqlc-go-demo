package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/thkx/graphql-sqlc-go-dome/internal/config"
	"github.com/thkx/graphql-sqlc-go-dome/internal/repository"
	"github.com/thkx/graphql-sqlc-go-dome/internal/repository/postgres/db"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// var _ repository.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	queries db.Querier
	db      *sql.DB
}

func NewPostgresRepository(cfg *config.Config) (*PostgresRepository, error) {
	// 创建或打开数据库
	db_conn, err := sql.Open("pgx", cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	db_conn.SetMaxOpenConns(cfg.MaxOpenConns)       // 最大打开连接数
	db_conn.SetMaxIdleConns(cfg.MaxIdleConns)       // 最大空闲连接数
	db_conn.SetConnMaxLifetime(cfg.ConnMaxLifetime) // 连接最大存活时间
	db_conn.SetConnMaxIdleTime(cfg.ConnMaxIdleTime) // 连接最大空闲时间

	// 验证连接有效性（可选但推荐）
	if err := db_conn.Ping(); err != nil {
		db_conn.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	// 执行数据库迁移
	if err := migrate(db_conn); err != nil {
		db_conn.Close()                                                  // 迁移失败关闭连接池，避免资源泄露
		return nil, fmt.Errorf("failed to migrate postgres db: %w", err) // 修正错误提示
	}

	return &PostgresRepository{
		queries: db.New(db_conn),
		db:      db_conn,
	}, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

func (r *PostgresRepository) ListUsers(ctx context.Context, arg *repository.User, pagination *repository.Pagination) ([]*repository.User, error) {
	params := db.ListUsersParams{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
		ID:     arg.ID.String(),
		Name:   arg.Name,
		Email:  arg.Email,
	}

	rows, err := r.queries.ListUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []*repository.User
	for _, row := range rows {
		user := &repository.User{
			ID:        row.ID,
			Name:      row.Name,
			Email:     row.Email,
			Password:  row.Password,
			Avatar:    row.Avatar,
			Gender:    row.Gender,
			Bio:       row.Bio,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
		result = append(result, user)
	}

	return result, nil
}

func (r *PostgresRepository) ListUsersForAdmin(ctx context.Context, arg *repository.User, pagination *repository.Pagination) ([]*repository.User, error) {
	params := db.ListUsersForAdminParams{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
		Name:   arg.Name,
		Email:  arg.Email,
		Gender: arg.Gender.String,
	}

	rows, err := r.queries.ListUsersForAdmin(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []*repository.User
	for _, row := range rows {
		user := &repository.User{
			ID:          row.ID,
			Name:        row.Name,
			Email:       row.Email,
			Password:    row.Password,
			Avatar:      row.Avatar,
			Gender:      row.Gender,
			Bio:         row.Bio,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
			DeletedType: row.DeletedType,
			DeletedAt:   row.DeletedAt,
		}
		result = append(result, user)
	}

	return result, nil
}

func (r *PostgresRepository) CreateUser(ctx context.Context, arg *repository.User) error {
	params := db.CreateUserParams{
		ID:       arg.ID,
		Email:    arg.Email,
		Name:     arg.Name,
		Password: arg.Password,
		Avatar:   arg.Avatar,
		Gender:   arg.Gender,
		Bio:      arg.Bio,
	}
	if err := r.queries.CreateUser(ctx, params); err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, arg *repository.User) error {
	params := db.UpdateUserParams{
		ID:     arg.ID,
		Gender: arg.Gender,
		Bio:    arg.Bio,
		Email: sql.NullString{
			String: arg.Email,
			Valid:  arg.Email != "",
		},
		Name: sql.NullString{
			String: arg.Name,
			Valid:  arg.Name != "",
		},
		Password: sql.NullString{
			String: arg.Password,
			Valid:  arg.Password != "",
		},
		Avatar: sql.NullString{
			String: arg.Avatar,
			Valid:  arg.Avatar != "",
		},
	}

	if err := r.queries.UpdateUser(ctx, params); err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.SoftDeleteUser(ctx, id); err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) ListPosts(ctx context.Context, arg *repository.Post, pagination *repository.Pagination) ([]*repository.Post, error) {
	params := db.ListPostsParams{
		Limit:   pagination.Limit,
		Offset:  pagination.Offset,
		Title:   arg.Title,
		Content: arg.Content,
		Pv:      arg.Pv,
		Author:  arg.Author.String(),
		ID:      arg.ID.String(),
	}

	rows, err := r.queries.ListPosts(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []*repository.Post
	for _, row := range rows {
		post := &repository.Post{
			ID:        row.ID,
			Author:    row.UserID,
			Title:     row.Title,
			Content:   row.Content,
			Pv:        row.Pv,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			User: repository.User{
				ID:     row.UserID,
				Email:  row.UserEmail,
				Name:   row.UserName,
				Avatar: row.UserAvatar,
				Gender: row.UserGender,
				Bio:    row.UserBio,
			},
		}

		result = append(result, post)
	}
	return result, nil
}

func (r *PostgresRepository) ListPostsForAdmin(ctx context.Context, arg *repository.Post, pagination *repository.Pagination) ([]*repository.Post, error) {
	params := db.ListPostsForAdminParams{
		Limit:   pagination.Limit,
		Offset:  pagination.Offset,
		Title:   arg.Title,
		Content: arg.Content,
		Pv:      arg.Pv,
		Author:  arg.Author.String(),
	}

	rows, err := r.queries.ListPostsForAdmin(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []*repository.Post
	for _, row := range rows {
		post := &repository.Post{
			ID:          row.ID,
			Author:      row.UserID,
			Title:       row.Title,
			Content:     row.Content,
			Pv:          row.Pv,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
			DeletedType: row.DeletedType,
			DeletedAt:   row.DeletedAt,
			User: repository.User{
				ID:     row.UserID,
				Email:  row.UserEmail,
				Name:   row.UserName,
				Avatar: row.UserAvatar,
				Gender: row.UserGender,
				Bio:    row.UserBio,
			},
		}

		result = append(result, post)
	}
	return result, nil
}

func (r *PostgresRepository) CreatePost(ctx context.Context, arg *repository.Post) error {
	params := db.CreatePostParams{
		ID:      arg.ID,
		Author:  arg.Author,
		Title:   arg.Title,
		Content: arg.Content,
		Pv:      arg.Pv,
	}
	if err := r.queries.CreatePost(ctx, params); err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) UpdatePost(ctx context.Context, arg *repository.Post) error {
	params := db.UpdatePostParams{
		ID: arg.ID,
		Author: uuid.NullUUID{
			UUID:  arg.Author,
			Valid: arg.Author != uuid.Nil,
		},
		Title: sql.NullString{
			String: arg.Title,
			Valid:  arg.Title != "",
		},
		Content: sql.NullString{
			String: arg.Content,
			Valid:  arg.Content != "",
		},
	}
	if err := r.queries.UpdatePost(ctx, params); err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) DeletePost(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.SoftDeletePost(ctx, id); err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) ListComments(ctx context.Context, arg *repository.Comment) ([]*repository.Comment, error) {
	params := db.ListCommentsParams{
		ID:      arg.ID.String(),
		Author:  arg.Author.String(),
		PostID:  arg.PostID.String(),
		Content: arg.Content,
	}
	rows, err := r.queries.ListComments(ctx, params)
	if err != nil {
		return nil, err
	}
	var result []*repository.Comment

	for _, row := range rows {
		com := &repository.Comment{
			ID:        row.ID,
			Content:   row.Content,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			User: repository.User{
				ID:     row.UserID,
				Email:  row.UserEmail,
				Name:   row.UserName,
				Avatar: row.UserAvatar,
				Gender: row.UserGender,
				Bio:    row.UserBio,
			},
		}

		result = append(result, com)
	}
	return result, nil
}

func (r *PostgresRepository) CreateComment(ctx context.Context, arg *repository.Comment) error {
	params := db.CreateCommentParams{
		ID:      arg.ID,
		Author:  arg.Author,
		PostID:  arg.PostID,
		Content: arg.Content,
	}
	if err := r.queries.CreateComment(ctx, params); err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) UpdateComment(ctx context.Context, arg *repository.Comment) error {
	params := db.UpdateCommentParams{
		ID: arg.ID,
		Content: sql.NullString{
			String: arg.Content,
			Valid:  arg.Content != "",
		},
	}
	if err := r.queries.UpdateComment(ctx, params); err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) DeleteComment(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.SoftDeleteComment(ctx, id); err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) ListIndicators(ctx context.Context, arg *repository.Indicator, pagination *repository.Pagination) ([]*repository.Indicator, error) {
	params := db.ListIndicatorsParams{
		Limit:         pagination.Limit,
		Offset:        pagination.Offset,
		ID:            arg.ID.String(),
		Indicator:     arg.Indicator,
		IndicatorType: arg.IndicatorType,
		MetaSource:    arg.MetaSource,
	}

	rows, err := r.queries.ListIndicators(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []*repository.Indicator
	for _, row := range rows {
		in := &repository.Indicator{
			ID:            row.ID,
			Indicator:     row.Indicator,
			IndicatorType: row.IndicatorType,
			MetaSource:    row.MetaSource,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		}
		result = append(result, in)
	}
	return result, nil
}

func (r *PostgresRepository) CreateIndicator(ctx context.Context, arg *repository.Indicator) error {
	params := db.CreateIndicatorParams{
		ID:            arg.ID,
		Indicator:     arg.Indicator,
		IndicatorType: arg.IndicatorType,
		MetaSource:    arg.MetaSource,
	}
	if err := r.queries.CreateIndicator(ctx, params); err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) UpdateIndicator(ctx context.Context, arg *repository.Indicator) error {
	params := db.UpdateIndicatorParams{
		ID: arg.ID,
		Indicator: sql.NullString{
			String: arg.Indicator,
			Valid:  arg.Indicator != "",
		},
		IndicatorType: sql.NullString{
			String: arg.IndicatorType,
			Valid:  arg.IndicatorType != "",
		},
		MetaSource: sql.NullString{
			String: arg.MetaSource,
			Valid:  arg.MetaSource != "",
		},
	}
	if err := r.queries.UpdateIndicator(ctx, params); err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) DeleteIndicator(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.SoftDeleteIndicator(ctx, id); err != nil {
		return err
	}

	return nil
}
