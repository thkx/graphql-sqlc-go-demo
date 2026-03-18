package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/thkx/graphql-sqlc-go-dome/internal/config"
	"github.com/thkx/graphql-sqlc-go-dome/internal/graph/model"
	"github.com/thkx/graphql-sqlc-go-dome/internal/repository/sqlite/db"
	_ "modernc.org/sqlite"
)

type SQLiteRepository struct {
	queries db.Querier
	db      *sql.DB
}

func NewSQLiteRepository(cfg *config.Config) (*SQLiteRepository, error) {
	// 创建或打开数据库
	db_conn, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	db_conn.SetMaxOpenConns(cfg.MaxOpenConns)       // 最大打开连接数
	db_conn.SetMaxIdleConns(cfg.MaxIdleConns)       // 最大空闲连接数
	db_conn.SetConnMaxLifetime(cfg.ConnMaxLifetime) // 连接最大存活时间
	db_conn.SetConnMaxIdleTime(cfg.ConnMaxIdleTime) // 连接最大空闲时间

	// 设置常用 PRAGMA
	if _, err := db_conn.Exec(`
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = NORMAL;
		PRAGMA foreign_keys = ON;
		PRAGMA busy_timeout = 5000;
	`); err != nil {
		db_conn.Close()
		return nil, fmt.Errorf("failed to set pragma: %w", err)
	}

	// 执行迁移（简单方式：读取嵌入或文件执行）
	if err := migrate(db_conn); err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	return &SQLiteRepository{
		queries: db.New(db_conn),
		db:      db_conn,
	}, nil
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

func (r *SQLiteRepository) ListUsers() ([]*model.User, error) {
	users, err := r.queries.ListUsers(context.Background())
	if err != nil {
		return nil, err
	}

	result := make([]*model.User, len(users))
	for i, u := range users {
		result[i] = &model.User{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Password:  u.Password,
			Avatar:    u.Avatar,
			Gender:    u.Gender.String,
			Bio:       u.Bio,
			CreatedAt: u.CreatedAt, // 如果需要在 model 中暴露时间，可添加字段
			UpdatedAt: u.UpdatedAt,
		}
	}

	return result, nil
}

func (r *SQLiteRepository) GetUserByID(id string) (*model.User, error) {
	row, err := r.queries.GetUserByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:       row.ID,
		Name:     row.Name,
		Email:    row.Email,
		Gender:   row.Gender.String,
		Avatar:   row.Avatar,
		Bio:      row.Bio,
		Password: row.Password,
	}, nil
}

func (r *SQLiteRepository) GetUserByEmail(email string) (*model.User, error) {
	row, err := r.queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:       row.ID,
		Name:     row.Name,
		Email:    row.Email,
		Gender:   row.Gender.String,
		Avatar:   row.Avatar,
		Bio:      row.Bio,
		Password: row.Password,
	}, nil
}

func (r *SQLiteRepository) GetUserByName(name string) (*model.User, error) {
	row, err := r.queries.GetUserByName(context.Background(), name)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:       row.ID,
		Name:     row.Name,
		Email:    row.Email,
		Gender:   row.Gender.String,
		Avatar:   row.Avatar,
		Bio:      row.Bio,
		Password: row.Password,
	}, nil
}

func (r *SQLiteRepository) CreateUser(user *model.User) error {
	params := db.CreateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	err := r.queries.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteRepository) UpdateUser(user *model.User) (*model.User, error) {
	params := db.UpdateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Avatar:   user.Avatar,
		Gender: sql.NullString{
			String: user.Gender,
			Valid:  true,
		},
		Bio: user.Bio,
	}
	err := r.queries.UpdateUser(context.Background(), params)
	if err != nil {
		return nil, err
	}
	return &model.User{ID: user.ID}, nil
}

func (r *SQLiteRepository) DeleteUser(id string) error {
	err := r.queries.SoftDeleteUser(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteRepository) ListPosts() ([]*model.Post, error) {
	posts, err := r.queries.ListPosts(context.Background())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Post, len(posts))
	for i, u := range posts {
		result[i] = &model.Post{
			ID:      u.ID,
			Title:   u.Title,
			Content: u.Content,
			Pv:      int32(u.Pv.Int64),
			Author: &model.User{
				ID:     u.UserID,
				Name:   u.UserName,
				Email:  u.UserEmail,
				Gender: u.UserGender.String,
				Avatar: u.UserAvatar,
				Bio:    u.UserBio,
			},
			CreatedAt: u.CreatedAt, // 如果需要在 model 中暴露时间，可添加字段
			UpdatedAt: u.UpdatedAt,
		}
	}

	return result, nil
}

func (r *SQLiteRepository) CreatePost(post *model.Post) error {
	params := db.CreatePostParams{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Pv: sql.NullInt64{
			Int64: 0,
			Valid: true,
		},
		UserID: post.Author.ID,
	}
	err := r.queries.CreatePost(context.Background(), params)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteRepository) GetPostByID(id string) (*model.Post, error) {
	row, err := r.queries.GetPostByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	result := &model.Post{
		ID:      row.ID,
		Title:   row.Title,
		Pv:      int32(row.Pv.Int64),
		Content: row.Content,
		Author: &model.User{
			ID:     row.UserID,
			Name:   row.UserName,
			Email:  row.UserEmail,
			Gender: row.UserGender.String,
			Avatar: row.UserAvatar,
			Bio:    row.UserBio,
		},
		CreatedAt: row.CreatedAt, // 如果需要在 model 中暴露时间，可添加字段
		UpdatedAt: row.UpdatedAt,
	}

	return result, err
}

func (r *SQLiteRepository) ListPostsByTitle(title string) ([]*model.Post, error) {
	rows, err := r.queries.ListPostsByTitle(context.Background(), title)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Post, len(rows))
	for i, u := range rows {
		result[i] = &model.Post{
			ID:      u.ID,
			Title:   u.Title,
			Pv:      int32(u.Pv.Int64),
			Content: u.Content,
			Author: &model.User{
				ID:     u.UserID,
				Name:   u.UserName,
				Email:  u.UserEmail,
				Gender: u.UserGender.String,
				Avatar: u.UserAvatar,
				Bio:    u.UserBio,
			},
			CreatedAt: u.CreatedAt, // 如果需要在 model 中暴露时间，可添加字段
			UpdatedAt: u.UpdatedAt,
		}
	}
	return result, err
}

func (r *SQLiteRepository) UpdatePost(post *model.Post) (*model.Post, error) {
	params := db.UpdatePostParams{
		ID:    post.ID,
		Title: post.Title,
		Pv: sql.NullInt64{
			Int64: int64(post.Pv),
			Valid: true,
		},
		Content: post.Content,
		UserID:  post.Author.ID,
	}

	err := r.queries.UpdatePost(context.Background(), params)
	if err != nil {
		return nil, err
	}
	return &model.Post{ID: post.ID}, nil
}

func (r *SQLiteRepository) DeletePost(id string) error {
	return r.queries.SoftDeletePost(context.Background(), id)
}

func (r *SQLiteRepository) ListComments() ([]*model.Comment, error) {
	rows, err := r.queries.ListComments(context.Background())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Comment, len(rows))
	for i, u := range rows {
		result[i] = &model.Comment{
			ID:      u.ID,
			UserID:  u.UserID,
			PostID:  u.PostID,
			Content: u.Content,
			Author: &model.User{
				ID:     u.UserID,
				Name:   u.UserName,
				Email:  u.UserEmail,
				Gender: u.UserGender.String,
				Avatar: u.UserAvatar,
				Bio:    u.UserBio,
			},
			CreatedAt: u.CreatedAt, // 如果需要在 model 中暴露时间，可添加字段
			UpdatedAt: u.UpdatedAt,
		}
	}
	return result, err
}

func (r *SQLiteRepository) GetCommentByID(id string) (*model.Comment, error) {
	row, err := r.queries.GetCommentByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	result := &model.Comment{
		ID:      row.ID,
		UserID:  row.UserID,
		PostID:  row.PostID,
		Content: row.Content,
		Author: &model.User{
			ID:     row.UserID,
			Name:   row.UserName,
			Email:  row.UserEmail,
			Gender: row.UserGender.String,
			Avatar: row.UserAvatar,
			Bio:    row.UserBio,
		},
		CreatedAt: row.CreatedAt, // 如果需要在 model 中暴露时间，可添加字段
		UpdatedAt: row.UpdatedAt,
	}
	return result, nil
}

func (r *SQLiteRepository) ListCommentsByUserID(user_id string) ([]*model.Comment, error) {
	rows, err := r.queries.ListCommentsByUserID(context.Background(), user_id)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Comment, len(rows))
	for i, u := range rows {
		result[i] = &model.Comment{
			ID:      u.ID,
			UserID:  u.UserID,
			PostID:  u.PostID,
			Content: u.Content,
			Author: &model.User{
				ID:     u.UserID,
				Name:   u.UserName,
				Email:  u.UserEmail,
				Gender: u.UserGender.String,
				Avatar: u.UserAvatar,
				Bio:    u.UserBio,
			},
			CreatedAt: u.CreatedAt, // 如果需要在 model 中暴露时间，可添加字段
			UpdatedAt: u.UpdatedAt,
		}
	}
	return result, err
}

func (r *SQLiteRepository) ListCommentsByPostID(post_id string) ([]*model.Comment, error) {
	rows, err := r.queries.ListCommentsByPostID(context.Background(), post_id)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Comment, len(rows))
	for i, u := range rows {
		result[i] = &model.Comment{
			ID:      u.ID,
			UserID:  u.UserID,
			PostID:  u.PostID,
			Content: u.Content,
			Author: &model.User{
				ID:     u.UserID,
				Name:   u.UserName,
				Email:  u.UserEmail,
				Gender: u.UserGender.String,
				Avatar: u.UserAvatar,
				Bio:    u.UserBio,
			},
			CreatedAt: u.CreatedAt, // 如果需要在 model 中暴露时间，可添加字段
			UpdatedAt: u.UpdatedAt,
		}
	}
	return result, err
}

func (r *SQLiteRepository) ListCommentsByContent(content string) ([]*model.Comment, error) {
	rows, err := r.queries.ListCommentsByContent(context.Background(), content)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Comment, len(rows))
	for i, u := range rows {
		result[i] = &model.Comment{
			ID:      u.ID,
			UserID:  u.UserID,
			PostID:  u.PostID,
			Content: u.Content,
			Author: &model.User{
				ID:     u.UserID,
				Name:   u.UserName,
				Email:  u.UserEmail,
				Gender: u.UserGender.String,
				Avatar: u.UserAvatar,
				Bio:    u.UserBio,
			},
			CreatedAt: u.CreatedAt, // 如果需要在 model 中暴露时间，可添加字段
			UpdatedAt: u.UpdatedAt,
		}
	}
	return result, err
}

func (r *SQLiteRepository) CreateComment(comment *model.Comment) error {
	params := db.CreateCommentParams{
		ID:      comment.ID,
		UserID:  comment.UserID,
		Content: comment.Content,
		PostID:  comment.PostID,
	}
	err := r.queries.CreateComment(context.Background(), params)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteRepository) UpdateComment(comment *model.Comment) (*model.Comment, error) {
	params := db.UpdateCommentParams{
		ID:      comment.ID,
		UserID:  comment.UserID,
		Content: comment.Content,
		PostID:  comment.GetPostID(),
	}
	err := r.queries.UpdateComment(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return &model.Comment{ID: comment.ID}, nil
}

func (r *SQLiteRepository) DeleteComment(id string) error {
	err := r.queries.SoftDeleteComment(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteRepository) CreateIndicator(ind *model.Indicator) error {
	params := db.CreateIndicatorParams{
		Indicator:     ind.Indicator,
		IndicatorType: ind.IndicatorType,
		MetaSource:    ind.MetaSource,
	}
	err := r.queries.CreateIndicator(context.Background(), params)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteRepository) ListIndicators() ([]*model.Indicator, error) {
	rows, err := r.queries.ListIndicators(context.Background())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Indicator, len(rows))
	for i, u := range rows {
		result[i] = &model.Indicator{
			ID:            u.ID,
			Indicator:     u.Indicator,
			IndicatorType: u.IndicatorType,
			MetaSource:    u.MetaSource,
		}
	}
	return result, err
}
