package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/thkx/graphql-sqlc-go-dome/internal/config"
	"github.com/thkx/graphql-sqlc-go-dome/internal/graph/model"
	"github.com/thkx/graphql-sqlc-go-dome/internal/repository/mysql/db"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLRepository struct {
	queries db.Querier
	db      *sql.DB
}

func NewMySQLRepository(cfg *config.Config) (*MySQLRepository, error) {
	// 创建或打开数据库
	db_conn, err := sql.Open("mysql", cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	db_conn.SetMaxOpenConns(cfg.MaxOpenConns)       // 最大打开连接数
	db_conn.SetMaxIdleConns(cfg.MaxIdleConns)       // 最大空闲连接数
	db_conn.SetConnMaxLifetime(cfg.ConnMaxLifetime) // 连接最大存活时间
	db_conn.SetConnMaxIdleTime(cfg.ConnMaxIdleTime) // 连接最大空闲时间

	// 验证连接有效性（sql.Open 仅初始化连接池，不实际建立连接）
	if err := db_conn.Ping(); err != nil {
		db_conn.Close() // 验证失败关闭连接池，避免资源泄露
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	// 执行迁移（简单方式：读取嵌入或文件执行）
	if err := migrate(db_conn); err != nil {
		return nil, fmt.Errorf("failed to open MySQL db: %w", err)
	}

	return &MySQLRepository{
		queries: db.New(db_conn),
		db:      db_conn,
	}, nil
}

func (r *MySQLRepository) Close() error {
	return r.db.Close()
}

func (r *MySQLRepository) GetUserByID(id string) (*model.User, error) {

	return nil, nil
}

func (r *MySQLRepository) GetUserByEmail(email string) (*model.User, error) {

	return nil, nil
}

func (r *MySQLRepository) GetUserByName(name string) (*model.User, error) {

	return nil, nil
}

func (r *MySQLRepository) CreateUser(user *model.User) error {
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

func (r *MySQLRepository) UpdateUser(user *model.User) (*model.User, error) {

	return &model.User{ID: user.ID}, nil
}

func (r *MySQLRepository) DeleteUser(id string) error {
	err := r.queries.SoftDeleteUser(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}

func (r *MySQLRepository) ListPosts() ([]*model.Post, error) {

	return nil, nil
}

func (r *MySQLRepository) CreatePost(post *model.Post) error {

	return nil
}

func (r *MySQLRepository) GetPostByID(id string) (*model.Post, error) {

	return nil, nil
}

func (r *MySQLRepository) ListPostsByTitle(title string) ([]*model.Post, error) {

	return nil, nil
}

func (r *MySQLRepository) UpdatePost(post *model.Post) (*model.Post, error) {

	return &model.Post{ID: post.ID}, nil
}

func (r *MySQLRepository) DeletePost(id string) error {
	return nil
}

func (r *MySQLRepository) ListComments() ([]*model.Comment, error) {

	return nil, nil
}

func (r *MySQLRepository) GetCommentByID(id string) (*model.Comment, error) {

	return nil, nil
}

func (r *MySQLRepository) ListCommentsByUserID(user_id string) ([]*model.Comment, error) {

	return nil, nil
}

func (r *MySQLRepository) ListCommentsByPostID(post_id string) ([]*model.Comment, error) {

	return nil, nil
}

func (r *MySQLRepository) ListCommentsByContent(content string) ([]*model.Comment, error) {

	return nil, nil
}

func (r *MySQLRepository) CreateComment(comment *model.Comment) error {

	return nil
}

func (r *MySQLRepository) UpdateComment(comment *model.Comment) (*model.Comment, error) {

	return &model.Comment{ID: comment.ID}, nil
}

func (r *MySQLRepository) DeleteComment(id string) error {

	return nil
}

func (r *MySQLRepository) CreateIndicator(ind *model.Indicator) error {
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

func (r *MySQLRepository) ListIndicators() ([]*model.Indicator, error) {
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
